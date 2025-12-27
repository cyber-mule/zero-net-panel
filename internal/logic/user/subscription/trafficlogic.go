package subscription

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// TrafficLogic handles user traffic usage queries.
type TrafficLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewTrafficLogic constructs TrafficLogic.
func NewTrafficLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrafficLogic {
	return &TrafficLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Traffic returns traffic usage details for a subscription.
func (l *TrafficLogic) Traffic(req *types.UserSubscriptionTrafficRequest) (*types.UserSubscriptionTrafficResponse, error) {
	user, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrForbidden
	}

	sub, err := l.svcCtx.Repositories.Subscription.Get(l.ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}
	if sub.UserID != user.ID {
		return nil, repository.ErrForbidden
	}

	var from *time.Time
	if req.From != nil && *req.From > 0 {
		ts := time.Unix(*req.From, 0).UTC()
		from = &ts
	}
	var to *time.Time
	if req.To != nil && *req.To > 0 {
		ts := time.Unix(*req.To, 0).UTC()
		to = &ts
	}

	opts := repository.ListTrafficUsageOptions{
		Page:              req.Page,
		PerPage:           req.PerPage,
		Protocol:          req.Protocol,
		NodeID:            req.NodeID,
		ProtocolBindingID: req.ProtocolBindingID,
		From:              from,
		To:                to,
	}

	records, total, err := l.svcCtx.Repositories.TrafficUsage.ListBySubscription(l.ctx, req.SubscriptionID, opts)
	if err != nil {
		return nil, err
	}

	rawTotal, chargedTotal, err := l.svcCtx.Repositories.TrafficUsage.SumBySubscription(l.ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}

	items := make([]types.UserTrafficUsageRecord, 0, len(records))
	for _, record := range records {
		items = append(items, types.UserTrafficUsageRecord{
			ID:                record.ID,
			Protocol:          record.Protocol,
			NodeID:            record.NodeID,
			ProtocolBindingID: record.ProtocolBindingID,
			BytesUp:           record.BytesUp,
			BytesDown:         record.BytesDown,
			RawBytes:          record.RawBytes,
			ChargedBytes:      record.ChargedBytes,
			Multiplier:        record.Multiplier,
			ObservedAt:        toUnixOrZero(record.ObservedAt),
		})
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.UserSubscriptionTrafficResponse{
		Summary: types.UserSubscriptionTrafficSummary{
			RawBytes:     rawTotal,
			ChargedBytes: chargedTotal,
		},
		Records:    items,
		Pagination: pagination,
	}, nil
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}
