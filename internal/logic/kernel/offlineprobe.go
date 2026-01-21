package kernel

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/pkg/kernel"
)

// OfflineProbeManager schedules compensation probes for offline nodes.
type OfflineProbeManager struct {
	svcCtx *svc.ServiceContext
	mu     sync.Mutex
	probes map[controlKey]*offlineProbe
}

// NewOfflineProbeManager constructs an offline probe manager.
func NewOfflineProbeManager(svcCtx *svc.ServiceContext) *OfflineProbeManager {
	manager := &OfflineProbeManager{
		svcCtx: svcCtx,
		probes: make(map[controlKey]*offlineProbe),
	}
	return manager
}

// Update refreshes offline targets and adjusts probe goroutines.
func (m *OfflineProbeManager) Update(ctx context.Context) {
	if m == nil || m.svcCtx == nil {
		return
	}
	targets, err := m.listOfflineTargets(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("kernel offline probe: list targets failed: %v", err)
		return
	}

	m.applyTargets(ctx, targets)
}

type controlKey struct {
	endpoint string
	token    string
	timeout  time.Duration
}

type offlineTarget struct {
	key                controlKey
	nodeIDs            []uint64
	meta               authDebug
	maxIntervalSeconds int
}

func (m *OfflineProbeManager) listOfflineTargets(ctx context.Context) ([]offlineTarget, error) {
	nodes, err := m.svcCtx.Repositories.Node.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	targets := make(map[controlKey]*offlineTarget)
	for _, node := range nodes {
		status := strings.ToLower(strings.TrimSpace(node.Status))
		if status != "offline" {
			continue
		}
		if !node.StatusSyncEnabled {
			continue
		}
		if node.KernelStatusPollIntervalSeconds <= 0 {
			continue
		}
		endpoint := strings.TrimSpace(node.ControlEndpoint)
		if endpoint == "" {
			continue
		}
		token := resolveControlToken(node)
		timeout := resolveKernelHTTPTimeout(node)
		key := controlKey{endpoint: endpoint, token: token, timeout: timeout}
		target := targets[key]
		if target == nil {
			target = &offlineTarget{
				key:  key,
				meta: buildAuthDebug(node),
			}
			targets[key] = target
		}
		target.nodeIDs = append(target.nodeIDs, node.ID)
		target.maxIntervalSeconds = mergeOfflineProbeInterval(target.maxIntervalSeconds, node.KernelOfflineProbeMaxIntervalSeconds)
	}

	results := make([]offlineTarget, 0, len(targets))
	for _, target := range targets {
		results = append(results, *target)
	}
	return results, nil
}

func mergeOfflineProbeInterval(current int, candidate int) int {
	if candidate < 0 {
		candidate = 0
	}
	if candidate == 0 {
		if current == 0 {
			return 0
		}
		return current
	}
	if current == 0 {
		return candidate
	}
	if candidate < current {
		return candidate
	}
	return current
}

func (m *OfflineProbeManager) applyTargets(ctx context.Context, targets []offlineTarget) {
	if len(targets) == 0 {
		m.stopAll()
		return
	}

	targetMap := make(map[controlKey]offlineTarget, len(targets))
	for _, target := range targets {
		targetMap[target.key] = target
	}

	var newProbes []*offlineProbe
	var stopProbes []*offlineProbe

	m.mu.Lock()
	for key, target := range targetMap {
		if probe, ok := m.probes[key]; ok {
			probe.updateNodeIDs(target.nodeIDs)
			probe.meta = target.meta
			probe.updateMaxInterval(target.maxIntervalSeconds)
			continue
		}
		probe := newOfflineProbe(target.key, target.meta, target.nodeIDs, target.maxIntervalSeconds)
		m.probes[key] = probe
		newProbes = append(newProbes, probe)
	}
	for key, probe := range m.probes {
		if _, ok := targetMap[key]; ok {
			continue
		}
		delete(m.probes, key)
		stopProbes = append(stopProbes, probe)
	}
	m.mu.Unlock()

	for _, probe := range stopProbes {
		probe.stop()
	}
	for _, probe := range newProbes {
		go probe.run(ctx, m.svcCtx)
	}
}

func (m *OfflineProbeManager) stopAll() {
	m.mu.Lock()
	probes := make([]*offlineProbe, 0, len(m.probes))
	for key, probe := range m.probes {
		probes = append(probes, probe)
		delete(m.probes, key)
	}
	m.mu.Unlock()

	for _, probe := range probes {
		probe.stop()
	}
}

type offlineProbe struct {
	key                controlKey
	meta               authDebug
	stopCh             chan struct{}
	doneCh             chan struct{}
	nodeIDs            atomic.Value
	maxIntervalSeconds atomic.Int64
}

func newOfflineProbe(key controlKey, meta authDebug, nodeIDs []uint64, maxIntervalSeconds int) *offlineProbe {
	probe := &offlineProbe{
		key:    key,
		meta:   meta,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	probe.updateNodeIDs(nodeIDs)
	probe.updateMaxInterval(maxIntervalSeconds)
	return probe
}

func (p *offlineProbe) updateNodeIDs(nodeIDs []uint64) {
	cloned := append([]uint64(nil), nodeIDs...)
	p.nodeIDs.Store(cloned)
}

func (p *offlineProbe) updateMaxInterval(seconds int) {
	if seconds < 0 {
		seconds = 0
	}
	p.maxIntervalSeconds.Store(int64(seconds))
}

func (p *offlineProbe) stop() {
	select {
	case <-p.stopCh:
		return
	default:
		close(p.stopCh)
	}
}

func (p *offlineProbe) run(ctx context.Context, svcCtx *svc.ServiceContext) {
	defer close(p.doneCh)

	attempt := 0
	lastDay := dayKey(time.Now())

	for {
		now := time.Now()
		currentDay := dayKey(now)
		if currentDay != lastDay {
			attempt = 0
			lastDay = currentDay
		}

		delaySeconds := 1 + attempt*2
		maxInterval := int(p.maxIntervalSeconds.Load())
		if maxInterval > 0 && delaySeconds > maxInterval {
			delaySeconds = maxInterval
		}
		delay := time.Duration(delaySeconds) * time.Second

		untilNextDay := nextDayStart(now).Sub(now)
		if untilNextDay > 0 && delay > untilNextDay {
			delay = untilNextDay
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-p.stopCh:
			timer.Stop()
			return
		case <-timer.C:
		}

		if dayKey(time.Now()) != lastDay {
			attempt = 0
			lastDay = dayKey(time.Now())
			continue
		}

		nodeIDs, _ := p.nodeIDs.Load().([]uint64)
		if len(nodeIDs) == 0 {
			continue
		}

		if err := probeControlEndpoint(ctx, svcCtx, p.key, nodeIDs, p.meta); err != nil {
			attempt++
			continue
		}
		return
	}
}

func probeControlEndpoint(ctx context.Context, svcCtx *svc.ServiceContext, key controlKey, nodeIDs []uint64, meta authDebug) error {
	client, err := kernel.NewControlClient(kernel.HTTPOptions{
		BaseURL: key.endpoint,
		Token:   key.token,
		Timeout: key.timeout,
	})
	if err != nil {
		return err
	}

	_, err = client.GetStatus(ctx)
	if err != nil {
		if isUnauthorized(err) {
			logx.WithContext(ctx).Errorf(
				"kernel offline probe auth failed endpoint=%s auth=%s ak=%s sk_fp=%s token_fp=%s nodes=%v",
				key.endpoint,
				meta.AuthType,
				meta.AccessKeyMasked,
				meta.SecretFingerprint,
				meta.TokenFingerprint,
				nodeIDs,
			)
		}
		return err
	}

	statusByID, err := loadNodeStatusByID(ctx, svcCtx)
	if err != nil {
		return fmt.Errorf("offline probe load node status: %w", err)
	}
	recovered := resolveRecoveredNodes(nodeIDs, statusByID)
	markNodeStatus(ctx, svcCtx, nodeIDs, "online")
	if len(recovered) > 0 {
		triggerKernelRecovery(ctx, svcCtx, recovered)
	}
	return nil
}

func loadNodeStatusByID(ctx context.Context, svcCtx *svc.ServiceContext) (map[uint64]string, error) {
	nodes, err := svcCtx.Repositories.Node.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	statusByID := make(map[uint64]string, len(nodes))
	for _, node := range nodes {
		statusByID[node.ID] = strings.ToLower(strings.TrimSpace(node.Status))
	}
	return statusByID, nil
}

func dayKey(t time.Time) int {
	year, month, day := t.Date()
	return year*10000 + int(month)*100 + day
}

func nextDayStart(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, t.Location())
}
