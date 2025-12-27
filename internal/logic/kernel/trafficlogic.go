package kernel

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// TrafficIngestLogic handles kernel traffic ingestion.
type TrafficIngestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewTrafficIngestLogic constructs TrafficIngestLogic.
func NewTrafficIngestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrafficIngestLogic {
	return &TrafficIngestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Ingest persists kernel traffic records and updates subscription usage.
func (l *TrafficIngestLogic) Ingest(req *types.KernelTrafficReportRequest) (*types.KernelTrafficIngestResponse, error) {
	if req == nil || len(req.Records) == 0 {
		return &types.KernelTrafficIngestResponse{Accepted: 0, Failed: 0}, nil
	}

	accepted := 0
	failed := 0
	for _, record := range req.Records {
		if err := l.ingestRecord(record); err != nil {
			failed++
			l.Errorf("kernel traffic ingest failed: %v", err)
			continue
		}
		accepted++
	}

	return &types.KernelTrafficIngestResponse{
		Accepted: accepted,
		Failed:   failed,
	}, nil
}

func (l *TrafficIngestLogic) ingestRecord(record types.KernelTrafficRecord) error {
	subscription, err := l.resolveSubscription(record)
	if err != nil {
		return err
	}

	protocol := strings.ToLower(strings.TrimSpace(record.Protocol))
	multiplier := l.resolveMultiplier(subscription.PlanName, protocol)
	raw := maxInt64(record.BytesUp+record.BytesDown, 0)
	charged := int64(math.Round(float64(raw) * multiplier))

	observedAt := time.Now().UTC()
	if record.ObservedAt > 0 {
		observedAt = time.Unix(record.ObservedAt, 0).UTC()
	}

	usage := repository.TrafficUsageRecord{
		UserID:            subscription.UserID,
		SubscriptionID:    subscription.ID,
		ProtocolBindingID: record.ProtocolBindingID,
		NodeID:            record.NodeID,
		Protocol:          protocol,
		BytesUp:           record.BytesUp,
		BytesDown:         record.BytesDown,
		RawBytes:          raw,
		ChargedBytes:      charged,
		Multiplier:        multiplier,
		ObservedAt:        observedAt,
	}

	return l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		if _, err := txRepos.TrafficUsage.Create(l.ctx, usage); err != nil {
			return err
		}
		_, err := txRepos.Subscription.IncrementTrafficUsage(l.ctx, subscription.ID, charged)
		return err
	})
}

func (l *TrafficIngestLogic) resolveSubscription(record types.KernelTrafficRecord) (repository.Subscription, error) {
	if record.SubscriptionID != 0 {
		return l.svcCtx.Repositories.Subscription.Get(l.ctx, record.SubscriptionID)
	}
	if record.UserID != 0 {
		return l.svcCtx.Repositories.Subscription.GetActiveByUser(l.ctx, record.UserID)
	}
	return repository.Subscription{}, repository.ErrInvalidArgument
}

func (l *TrafficIngestLogic) resolveMultiplier(planName, protocol string) float64 {
	if planName == "" {
		return 1
	}

	plan, err := l.svcCtx.Repositories.Plan.GetByName(l.ctx, planName)
	if err != nil {
		return 1
	}

	if plan.TrafficMultipliers == nil {
		return 1
	}
	if protocol == "" {
		return 1
	}

	if multiplier, ok := plan.TrafficMultipliers[protocol]; ok && multiplier > 0 {
		return multiplier
	}
	return 1
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
