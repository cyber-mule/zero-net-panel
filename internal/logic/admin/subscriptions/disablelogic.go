package subscriptions

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DisableLogic handles subscription disabling.
type DisableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDisableLogic constructs DisableLogic.
func NewDisableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableLogic {
	return &DisableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Disable sets the subscription status to disabled.
func (l *DisableLogic) Disable(req *types.AdminDisableSubscriptionRequest) (*types.AdminSubscriptionResponse, error) {
	statusCode := status.SubscriptionStatusDisabled

	var updated repository.Subscription
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		result, err := txRepos.Subscription.Update(l.ctx, req.SubscriptionID, repository.UpdateSubscriptionInput{
			Status: &statusCode,
		})
		if err != nil {
			return err
		}
		updated = result

		actor, ok := security.UserFromContext(l.ctx)
		var actorID *uint64
		if ok && actor.ID != 0 {
			actorID = &actor.ID
		}

		metadata := map[string]any{}
		if req.Reason != nil && *req.Reason != "" {
			metadata["reason"] = *req.Reason
		}

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.subscription.disable",
			ResourceType: "subscription",
			ResourceID:   fmt.Sprintf("%d", updated.ID),
			Metadata:     metadata,
		})
		return err
	}); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, updated.UserID)
	if err != nil {
		return nil, err
	}

	return &types.AdminSubscriptionResponse{
		Subscription: toAdminSubscriptionSummary(updated, user),
	}, nil
}
