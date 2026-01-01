package cli

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
	"github.com/zero-net-panel/zero-net-panel/internal/handler"
	publicsubscriptions "github.com/zero-net-panel/zero-net-panel/internal/handler/public/subscriptions"
	kernellogic "github.com/zero-net-panel/zero-net-panel/internal/logic/kernel"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
)

// RunAPIServer 保留旧名称以兼容现有入口，内部委托给 RunServices。
func RunAPIServer(ctx context.Context, cfg config.Config) error {
	return RunServices(ctx, cfg)
}

// RunServices 启动 HTTP 与 gRPC 服务，并在任一退出或外部取消时统一回收资源。
func RunServices(ctx context.Context, cfg config.Config) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	svcCtx, err := svc.NewServiceContext(cfg)
	if err != nil {
		return err
	}
	defer svcCtx.Cleanup()

	proc.AddShutdownListener(func() {
		cancel()
	})

	errCh := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runHTTPServer(runCtx, cfg, svcCtx); err != nil {
			errCh <- err
		}
	}()

	if cfg.Metrics.Enabled() && cfg.Metrics.Standalone() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runMetricsServer(runCtx, cfg.Metrics); err != nil {
				errCh <- err
			}
		}()
	}

	if cfg.GRPC.Enabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runGRPCServer(runCtx, cfg, svcCtx); err != nil {
				errCh <- err
			}
		}()
	}

	if svcCtx.KernelControl != nil && cfg.Kernel.StatusPollInterval > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runKernelStatusPoller(runCtx, svcCtx, cfg.Kernel.StatusPollInterval)
		}()
	}

	var runErr error
	select {
	case <-runCtx.Done():
	case runErr = <-errCh:
		cancel()
	}

	wg.Wait()

	return runErr
}

func runKernelStatusPoller(ctx context.Context, svcCtx *svc.ServiceContext, interval time.Duration) {
	logger := logx.WithContext(ctx)
	backoff := newKernelStatusBackoff(interval, svcCtx.Config.Kernel.StatusPollBackoff)
	nextDelay := interval

	for {
		if err := kernellogic.SyncStatus(ctx, svcCtx); err != nil {
			nextDelay = backoff.NextDelay()
			logger.Errorf("kernel status poll failed: %v (next=%s)", err, nextDelay)
		} else {
			backoff.Reset()
			nextDelay = interval
		}

		timer := time.NewTimer(nextDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

type kernelStatusBackoff struct {
	enabled    bool
	base       time.Duration
	max        time.Duration
	multiplier float64
	jitter     float64
	failures   int
	rng        *rand.Rand
}

func newKernelStatusBackoff(base time.Duration, cfg config.KernelBackoff) *kernelStatusBackoff {
	b := &kernelStatusBackoff{
		enabled:    cfg.Enabled,
		base:       base,
		max:        cfg.MaxInterval,
		multiplier: cfg.Multiplier,
		jitter:     cfg.Jitter,
	}
	if b.enabled {
		b.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if b.max <= 0 {
		b.max = base
	}
	if b.multiplier <= 1 {
		b.multiplier = 2
	}
	if b.jitter < 0 {
		b.jitter = 0
	}
	if b.jitter > 1 {
		b.jitter = 1
	}
	if b.max < b.base {
		b.max = b.base
	}
	return b
}

func (b *kernelStatusBackoff) NextDelay() time.Duration {
	if !b.enabled {
		return b.base
	}
	b.failures++
	delay := float64(b.base) * math.Pow(b.multiplier, float64(b.failures))
	if delay > float64(b.max) {
		delay = float64(b.max)
	}
	if b.jitter > 0 && b.rng != nil {
		jitter := (b.rng.Float64()*2 - 1) * b.jitter
		delay = delay * (1 + jitter)
	}
	if delay < float64(b.base) {
		delay = float64(b.base)
	}
	return time.Duration(delay)
}

func (b *kernelStatusBackoff) Reset() {
	b.failures = 0
}

func runHTTPServer(ctx context.Context, cfg config.Config, svcCtx *svc.ServiceContext) error {
	server := rest.MustNewServer(cfg.RestConf, corsOptions(cfg.CORS)...)
	defer server.Stop()

	if cfg.Metrics.Enabled() && !cfg.Metrics.Standalone() {
		metricsHandler := promhttp.Handler()
		server.AddRoute(rest.Route{
			Method:  http.MethodGet,
			Path:    cfg.Metrics.Path,
			Handler: metricsHandler.ServeHTTP,
		})
		fmt.Printf("Prometheus metrics available at http://%s:%d%s\n", cfg.Host, cfg.Port, cfg.Metrics.Path)
	}

	handler.RegisterHandlers(server, svcCtx)

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/subscriptions/:token",
		Handler: publicsubscriptions.PublicSubscriptionDownloadHandler(svcCtx),
	})

	fmt.Printf("Starting HTTP API at %s:%d...\n", cfg.Host, cfg.Port)

	done := make(chan struct{})
	go func() {
		server.Start()
		close(done)
	}()

	select {
	case <-ctx.Done():
		server.Stop()
		<-done
		return nil
	case <-done:
		return nil
	}
}

func corsOptions(cfg config.CORSConfig) []rest.RunOption {
	if !cfg.Enabled {
		return nil
	}

	origins := cfg.AllowOrigins
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	if len(cfg.AllowHeaders) == 0 {
		return []rest.RunOption{rest.WithCors(origins...)}
	}

	headers := append([]string(nil), cfg.AllowHeaders...)
	return []rest.RunOption{
		rest.WithCustomCors(func(header http.Header) {
			header.Add("Access-Control-Allow-Headers", strings.Join(headers, ", "))
		}, nil, origins...),
	}
}

func runMetricsServer(ctx context.Context, cfg config.MetricsConfig) error {
	mux := http.NewServeMux()
	mux.Handle(cfg.Path, promhttp.Handler())

	server := &http.Server{
		Addr:    cfg.ListenOn,
		Handler: mux,
	}

	fmt.Printf("Starting Prometheus metrics server at %s%s...\n", cfg.ListenOn, cfg.Path)

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		if err, ok := <-errCh; ok && err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}
