package kernel

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/pkg/kernel"
)

// OfflineProbeManager schedules compensation probes for offline nodes.
type OfflineProbeManager struct {
	svcCtx                    *svc.ServiceContext
	defaultMaxIntervalSeconds int
	maxIntervalSeconds        atomic.Int64
	mu                        sync.Mutex
	probes                    map[controlKey]*offlineProbe
}

// NewOfflineProbeManager constructs an offline probe manager.
func NewOfflineProbeManager(svcCtx *svc.ServiceContext, defaultMaxIntervalSeconds int) *OfflineProbeManager {
	manager := &OfflineProbeManager{
		svcCtx:                    svcCtx,
		defaultMaxIntervalSeconds: defaultMaxIntervalSeconds,
		probes:                    make(map[controlKey]*offlineProbe),
	}
	manager.maxIntervalSeconds.Store(int64(defaultMaxIntervalSeconds))
	return manager
}

// Update refreshes offline targets and adjusts probe goroutines.
func (m *OfflineProbeManager) Update(ctx context.Context) {
	if m == nil || m.svcCtx == nil {
		return
	}
	maxInterval := m.resolveMaxIntervalSeconds(ctx)
	m.maxIntervalSeconds.Store(int64(maxInterval))

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
}

type offlineTarget struct {
	key     controlKey
	nodeIDs []uint64
	meta    authDebug
}

func (m *OfflineProbeManager) resolveMaxIntervalSeconds(ctx context.Context) int {
	defaults := repository.SiteSettingDefaults{
		Name:                                 m.svcCtx.Config.Site.Name,
		LogoURL:                              m.svcCtx.Config.Site.LogoURL,
		KernelOfflineProbeMaxIntervalSeconds: m.defaultMaxIntervalSeconds,
	}
	setting, err := m.svcCtx.Repositories.Site.GetSiteSetting(ctx, defaults)
	if err != nil {
		return m.defaultMaxIntervalSeconds
	}
	if setting.KernelOfflineProbeMaxIntervalSeconds < 0 {
		return 0
	}
	return setting.KernelOfflineProbeMaxIntervalSeconds
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
		endpoint := strings.TrimSpace(node.ControlEndpoint)
		if endpoint == "" {
			continue
		}
		token := resolveControlToken(node)
		key := controlKey{endpoint: endpoint, token: token}
		target := targets[key]
		if target == nil {
			target = &offlineTarget{
				key:  key,
				meta: buildAuthDebug(node),
			}
			targets[key] = target
		}
		target.nodeIDs = append(target.nodeIDs, node.ID)
	}

	results := make([]offlineTarget, 0, len(targets))
	for _, target := range targets {
		results = append(results, *target)
	}
	return results, nil
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
			continue
		}
		probe := newOfflineProbe(target.key, target.meta, target.nodeIDs)
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
		go probe.run(ctx, m.svcCtx, &m.maxIntervalSeconds)
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
	key     controlKey
	meta    authDebug
	stopCh  chan struct{}
	doneCh  chan struct{}
	nodeIDs atomic.Value
}

func newOfflineProbe(key controlKey, meta authDebug, nodeIDs []uint64) *offlineProbe {
	probe := &offlineProbe{
		key:    key,
		meta:   meta,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	probe.updateNodeIDs(nodeIDs)
	return probe
}

func (p *offlineProbe) updateNodeIDs(nodeIDs []uint64) {
	cloned := append([]uint64(nil), nodeIDs...)
	p.nodeIDs.Store(cloned)
}

func (p *offlineProbe) stop() {
	select {
	case <-p.stopCh:
		return
	default:
		close(p.stopCh)
	}
}

func (p *offlineProbe) run(ctx context.Context, svcCtx *svc.ServiceContext, maxIntervalSeconds *atomic.Int64) {
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
		maxInterval := int(maxIntervalSeconds.Load())
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
		Timeout: svcCtx.Config.Kernel.HTTP.Timeout,
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
