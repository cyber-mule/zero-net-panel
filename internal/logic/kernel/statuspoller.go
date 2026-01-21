package kernel

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/nodecfg"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/pkg/kernel"
)

const (
	statusPollTickInterval    = time.Second
	statusPollRefreshInterval = 10 * time.Second
)

type nodePollState struct {
	next       time.Time
	backoff    *statusBackoff
	cfgKey     string
	lastStatus string
}

// RunStatusPoller schedules per-node kernel status polling.
func RunStatusPoller(ctx context.Context, svcCtx *svc.ServiceContext) {
	if svcCtx == nil {
		return
	}
	poller := newStatusPoller(svcCtx)
	poller.run(ctx)
}

type statusPoller struct {
	svcCtx       *svc.ServiceContext
	lastRefresh  time.Time
	nodes        []repository.Node
	states       map[uint64]*nodePollState
	offlineProbe *OfflineProbeManager
}

func newStatusPoller(svcCtx *svc.ServiceContext) *statusPoller {
	return &statusPoller{
		svcCtx:       svcCtx,
		states:       make(map[uint64]*nodePollState),
		offlineProbe: NewOfflineProbeManager(svcCtx),
	}
}

func (p *statusPoller) run(ctx context.Context) {
	logger := logx.WithContext(ctx)
	ticker := time.NewTicker(statusPollTickInterval)
	defer ticker.Stop()

	for {
		if err := p.refreshNodes(ctx); err != nil {
			logger.Errorf("kernel status poll refresh failed: %v", err)
		}
		p.pollDue(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (p *statusPoller) refreshNodes(ctx context.Context) error {
	if !p.shouldRefresh() {
		return nil
	}

	nodes, err := p.svcCtx.Repositories.Node.ListAll(ctx)
	if err != nil {
		return err
	}
	p.nodes = nodes
	p.lastRefresh = time.Now().UTC()
	p.syncStates(nodes)
	p.offlineProbe.Update(ctx)
	return nil
}

func (p *statusPoller) shouldRefresh() bool {
	if p.lastRefresh.IsZero() {
		return true
	}
	return time.Since(p.lastRefresh) >= statusPollRefreshInterval
}

func (p *statusPoller) syncStates(nodes []repository.Node) {
	active := make(map[uint64]struct{}, len(nodes))
	for _, node := range nodes {
		if !isPollEligible(node) {
			delete(p.states, node.ID)
			continue
		}
		active[node.ID] = struct{}{}

		cfgKey := pollConfigKey(node)
		status := normalizeStatus(node.Status)
		interval := time.Duration(node.KernelStatusPollIntervalSeconds) * time.Second
		backoffCfg := nodecfg.KernelBackoffConfig{
			Enabled:            node.KernelStatusPollBackoffEnabled,
			MaxIntervalSeconds: node.KernelStatusPollBackoffMaxIntervalSeconds,
			Multiplier:         node.KernelStatusPollBackoffMultiplier,
			Jitter:             node.KernelStatusPollBackoffJitter,
		}

		state := p.states[node.ID]
		if state == nil {
			p.states[node.ID] = &nodePollState{
				backoff:    newStatusBackoff(interval, backoffCfg),
				cfgKey:     cfgKey,
				lastStatus: status,
			}
			continue
		}
		if state.cfgKey != cfgKey {
			state.backoff = newStatusBackoff(interval, backoffCfg)
			state.cfgKey = cfgKey
			state.next = time.Time{}
		}
		state.lastStatus = status
	}

	for nodeID := range p.states {
		if _, ok := active[nodeID]; !ok {
			delete(p.states, nodeID)
		}
	}
}

func (p *statusPoller) pollDue(ctx context.Context) {
	now := time.Now().UTC()
	for _, node := range p.nodes {
		if !isPollEligible(node) {
			continue
		}
		state := p.states[node.ID]
		if state == nil {
			continue
		}
		if !state.next.IsZero() && now.Before(state.next) {
			continue
		}

		interval := time.Duration(node.KernelStatusPollIntervalSeconds) * time.Second
		if interval <= 0 {
			continue
		}
		nextStatus, err := pollNodeStatus(ctx, p.svcCtx, node, state.lastStatus)
		if err != nil {
			state.next = now.Add(state.backoff.NextDelay())
			state.lastStatus = nextStatus
			continue
		}
		state.backoff.Reset()
		state.next = now.Add(interval)
		state.lastStatus = nextStatus
	}
}

func pollNodeStatus(ctx context.Context, svcCtx *svc.ServiceContext, node repository.Node, previousStatus string) (string, error) {
	endpoint := strings.TrimSpace(node.ControlEndpoint)
	if endpoint == "" {
		markNodeStatus(ctx, svcCtx, []uint64{node.ID}, "offline")
		return "offline", fmt.Errorf("node control endpoint not configured")
	}
	token := resolveControlToken(node)
	client, err := kernel.NewControlClient(kernel.HTTPOptions{
		BaseURL: endpoint,
		Token:   token,
		Timeout: resolveKernelHTTPTimeout(node),
	})
	if err != nil {
		markNodeStatus(ctx, svcCtx, []uint64{node.ID}, "offline")
		return "offline", err
	}

	_, err = client.GetStatus(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("kernel status poll failed for %s: %v", endpoint, err)
		if isUnauthorized(err) {
			meta := buildAuthDebug(node)
			logx.WithContext(ctx).Errorf(
				"kernel status auth debug endpoint=%s auth=%s ak=%s sk_fp=%s token_fp=%s nodes=%v",
				endpoint,
				meta.AuthType,
				meta.AccessKeyMasked,
				meta.SecretFingerprint,
				meta.TokenFingerprint,
				[]uint64{node.ID},
			)
		}
		markNodeStatus(ctx, svcCtx, []uint64{node.ID}, "offline")
		return "offline", err
	}

	markNodeStatus(ctx, svcCtx, []uint64{node.ID}, "online")
	prev := normalizeStatus(previousStatus)
	if prev != "online" && prev != "disabled" {
		triggerKernelRecovery(ctx, svcCtx, []uint64{node.ID})
	}
	return "online", nil
}

func isPollEligible(node repository.Node) bool {
	if !node.StatusSyncEnabled {
		return false
	}
	if strings.EqualFold(node.Status, "disabled") {
		return false
	}
	if strings.TrimSpace(node.ControlEndpoint) == "" {
		return false
	}
	return node.KernelStatusPollIntervalSeconds > 0
}

func normalizeStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func pollConfigKey(node repository.Node) string {
	endpoint := strings.TrimSpace(node.ControlEndpoint)
	tokenFingerprint := fingerprint(resolveControlToken(node))
	return fmt.Sprintf(
		"%s|%s|%d|%t|%d|%g|%g|%d",
		endpoint,
		tokenFingerprint,
		node.KernelStatusPollIntervalSeconds,
		node.KernelStatusPollBackoffEnabled,
		node.KernelStatusPollBackoffMaxIntervalSeconds,
		node.KernelStatusPollBackoffMultiplier,
		node.KernelStatusPollBackoffJitter,
		node.KernelHTTPTimeoutSeconds,
	)
}

type statusBackoff struct {
	enabled    bool
	base       time.Duration
	max        time.Duration
	multiplier float64
	jitter     float64
	failures   int
	rng        *rand.Rand
}

func newStatusBackoff(base time.Duration, cfg nodecfg.KernelBackoffConfig) *statusBackoff {
	b := &statusBackoff{
		enabled:    cfg.Enabled,
		base:       base,
		max:        time.Duration(cfg.MaxIntervalSeconds) * time.Second,
		multiplier: cfg.Multiplier,
		jitter:     cfg.Jitter,
	}
	if b.enabled {
		b.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if b.base <= 0 {
		b.base = time.Duration(nodecfg.DefaultKernelStatusPollIntervalSeconds) * time.Second
	}
	if b.max <= 0 {
		b.max = b.base
	}
	if b.multiplier <= 1 {
		b.multiplier = nodecfg.DefaultKernelStatusPollBackoffMultiplier
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

func (b *statusBackoff) NextDelay() time.Duration {
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

func (b *statusBackoff) Reset() {
	b.failures = 0
}
