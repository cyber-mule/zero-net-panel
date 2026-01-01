package subscriptions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	subscriptionutil "github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic handles subscription updates.
type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLogic constructs UpdateLogic.
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update patches subscription fields.
func (l *UpdateLogic) Update(req *types.AdminUpdateSubscriptionRequest) (*types.AdminSubscriptionResponse, error) {
	status, err := normalizeOptionalStatus(req.Status)
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		value, err := validateExpiry(*req.ExpiresAt)
		if err != nil {
			return nil, err
		}
		expiresAt = &value
	}

	if req.Token != nil && strings.TrimSpace(*req.Token) == "" {
		return nil, repository.ErrInvalidArgument
	}

	if req.TemplateID != nil && *req.TemplateID == 0 {
		return nil, repository.ErrInvalidArgument
	}
	if req.PlanID != nil && *req.PlanID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	if req.AvailableTemplateIDs != nil {
		for _, id := range *req.AvailableTemplateIDs {
			if id == 0 {
				return nil, repository.ErrInvalidArgument
			}
		}
	}

	sub, err := l.svcCtx.Repositories.Subscription.Get(l.ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	var planName *string
	var planSnapshot *map[string]any
	if req.PlanID != nil {
		plan, err := l.svcCtx.Repositories.Plan.Get(l.ctx, *req.PlanID)
		if err != nil {
			return nil, err
		}
		bindingIDs, err := l.svcCtx.Repositories.PlanProtocolBinding.ListBindingIDs(l.ctx, plan.ID)
		if err != nil {
			return nil, err
		}
		name := strings.TrimSpace(plan.Name)
		if req.PlanName != nil && strings.TrimSpace(*req.PlanName) != "" {
			name = strings.TrimSpace(*req.PlanName)
		}
		if name == "" {
			return nil, repository.ErrInvalidArgument
		}
		planName = &name
		snapshot := subscriptionutil.BuildPlanSnapshot(plan, bindingIDs)
		planSnapshot = &snapshot
	} else if req.PlanName != nil {
		if strings.TrimSpace(*req.PlanName) == "" {
			return nil, repository.ErrInvalidArgument
		}
		trimmed := strings.TrimSpace(*req.PlanName)
		planName = &trimmed
	}

	if req.TemplateID != nil {
		if _, err := l.svcCtx.Repositories.SubscriptionTemplate.Get(l.ctx, *req.TemplateID); err != nil {
			return nil, err
		}
	}
	if req.AvailableTemplateIDs != nil {
		for _, id := range *req.AvailableTemplateIDs {
			if _, err := l.svcCtx.Repositories.SubscriptionTemplate.Get(l.ctx, id); err != nil {
				return nil, err
			}
		}
	}

	if req.TrafficTotalBytes != nil {
		if *req.TrafficTotalBytes < 0 {
			return nil, repository.ErrInvalidArgument
		}
		used := sub.TrafficUsedBytes
		if req.TrafficUsedBytes != nil {
			used = *req.TrafficUsedBytes
		}
		if *req.TrafficTotalBytes > 0 && used > *req.TrafficTotalBytes {
			return nil, repository.ErrInvalidArgument
		}
	}
	if req.TrafficUsedBytes != nil && *req.TrafficUsedBytes < 0 {
		return nil, repository.ErrInvalidArgument
	}
	if req.DevicesLimit != nil && *req.DevicesLimit <= 0 {
		return nil, repository.ErrInvalidArgument
	}

	input := repository.UpdateSubscriptionInput{
		Name:                 req.Name,
		PlanName:             planName,
		PlanID:               req.PlanID,
		PlanSnapshot:         planSnapshot,
		Status:               status,
		TemplateID:           req.TemplateID,
		AvailableTemplateIDs: req.AvailableTemplateIDs,
		Token:                req.Token,
		ExpiresAt:            expiresAt,
		TrafficTotalBytes:    req.TrafficTotalBytes,
		TrafficUsedBytes:     req.TrafficUsedBytes,
		DevicesLimit:         req.DevicesLimit,
	}

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

		metadata := map[string]any{}
		if req.Status != nil {
			metadata["status"] = updated.Status
		}
		if req.PlanID != nil {
			metadata["plan_id"] = updated.PlanID
		}
		if planName != nil {
			metadata["plan_name"] = updated.PlanName
		}
		if req.TemplateID != nil {
			metadata["template_id"] = updated.TemplateID
		}
		if req.ExpiresAt != nil {
			metadata["expires_at"] = updated.ExpiresAt.Unix()
		}

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.subscription.update",
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
