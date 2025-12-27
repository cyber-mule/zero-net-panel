package subscriptions

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// GetLogic handles subscription lookup.
type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetLogic constructs GetLogic.
func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Get returns subscription detail.
func (l *GetLogic) Get(req *types.AdminGetSubscriptionRequest) (*types.AdminSubscriptionResponse, error) {
	sub, err := l.svcCtx.Repositories.Subscription.Get(l.ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, sub.UserID)
	if err != nil {
		return nil, err
	}

	return &types.AdminSubscriptionResponse{
		Subscription: toAdminSubscriptionSummary(sub, user),
	}, nil
}
