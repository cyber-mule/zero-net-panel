package subscriptions

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	subscriptionutil "github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ExtendLogic handles subscription extension.
type ExtendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewExtendLogic constructs ExtendLogic.
func NewExtendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExtendLogic {
	return &ExtendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Extend extends subscription expiry.
func (l *ExtendLogic) Extend(req *types.AdminExtendSubscriptionRequest) (*types.AdminSubscriptionResponse, error) {
	duration, absolute, err := validateExtendRequest(req.ExtendDays, req.ExtendHours, req.ExpiresAt)
	if err != nil {
		return nil, err
	}

	sub, err := l.svcCtx.Repositories.Subscription.Get(l.ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}

	target := time.Now().UTC()
	if absolute != nil {
		target = *absolute
	} else {
		base := sub.ExpiresAt
		if base.IsZero() {
			base = time.Now().UTC()
		}
		target = base.Add(duration)
	}

	updateStatus := ""
	if sub.Status == "expired" && target.After(time.Now().UTC()) {
		updateStatus = "active"
	}

	input := repository.UpdateSubscriptionInput{
		ExpiresAt: &target,
	}
	if updateStatus != "" {
		input.Status = &updateStatus
	}

	now := time.Now().UTC()
	var updated repository.Subscription
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		result, err := txRepos.Subscription.Update(l.ctx, req.SubscriptionID, input)
		if err != nil {
			return err
		}
		updated = result
		if subscriptionutil.IsSubscriptionEffective(updated, now) {
			if err := txRepos.Subscription.DisableOtherActive(l.ctx, updated.UserID, updated.ID); err != nil {
				return err
			}
		}

		actor, ok := security.UserFromContext(l.ctx)
		var actorID *uint64
		if ok && actor.ID != 0 {
			actorID = &actor.ID
		}

		meta := map[string]any{
			"expires_at": target.Unix(),
		}
		if absolute == nil {
			meta["extend_seconds"] = int64(duration.Seconds())
		}

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.subscription.extend",
			ResourceType: "subscription",
			ResourceID:   fmt.Sprintf("%d", updated.ID),
			Metadata:     meta,
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
