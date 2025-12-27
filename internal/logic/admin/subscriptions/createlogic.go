package subscriptions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles subscription creation.
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic constructs CreateLogic.
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create provisions a new subscription.
func (l *CreateLogic) Create(req *types.AdminCreateSubscriptionRequest) (*types.AdminSubscriptionResponse, error) {
	if req.UserID == 0 || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.PlanName) == "" || req.TemplateID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	expiresAt, err := validateExpiry(req.ExpiresAt)
	if err != nil {
		return nil, err
	}

	status := "active"
	if req.Status != nil {
		normalized, err := normalizeStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		status = normalized
	}

	available := append([]uint64(nil), req.AvailableTemplateIDs...)
	if len(available) == 0 {
		available = []uint64{req.TemplateID}
	}

	token := ""
	if req.Token != nil {
		token = strings.TrimSpace(*req.Token)
	}
	if token == "" {
		token, err = generateToken()
		if err != nil {
			return nil, err
		}
	}

	_, err = l.svcCtx.Repositories.User.Get(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := l.ensureTemplates(req.TemplateID, available); err != nil {
		return nil, err
	}

	usedBytes := int64(0)
	if req.TrafficUsedBytes != nil {
		usedBytes = *req.TrafficUsedBytes
	}
	if req.TrafficTotalBytes < 0 || usedBytes < 0 {
		return nil, repository.ErrInvalidArgument
	}
	if req.TrafficTotalBytes > 0 && usedBytes > req.TrafficTotalBytes {
		return nil, repository.ErrInvalidArgument
	}
	if req.DevicesLimit <= 0 {
		return nil, repository.ErrInvalidArgument
	}

	now := time.Now().UTC()
	subscription := repository.Subscription{
		UserID:               req.UserID,
		Name:                 strings.TrimSpace(req.Name),
		PlanName:             strings.TrimSpace(req.PlanName),
		Status:               status,
		TemplateID:           req.TemplateID,
		AvailableTemplateIDs: available,
		Token:                token,
		ExpiresAt:            expiresAt,
		TrafficTotalBytes:    req.TrafficTotalBytes,
		TrafficUsedBytes:     usedBytes,
		DevicesLimit:         req.DevicesLimit,
		LastRefreshedAt:      now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	var created repository.Subscription
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		sub, err := txRepos.Subscription.Create(l.ctx, subscription)
		if err != nil {
			return err
		}
		created = sub

		actor, ok := security.UserFromContext(l.ctx)
		var actorID *uint64
		if ok && actor.ID != 0 {
			actorID = &actor.ID
		}

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.subscription.create",
			ResourceType: "subscription",
			ResourceID:   fmt.Sprintf("%d", created.ID),
			Metadata: map[string]any{
				"user_id":  req.UserID,
				"plan":     subscription.PlanName,
				"status":   subscription.Status,
				"template": subscription.TemplateID,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, created.UserID)
	if err != nil {
		return nil, err
	}

	return &types.AdminSubscriptionResponse{
		Subscription: toAdminSubscriptionSummary(created, user),
	}, nil
}

func (l *CreateLogic) ensureTemplates(primary uint64, available []uint64) error {
	if primary == 0 {
		return repository.ErrInvalidArgument
	}

	if _, err := l.svcCtx.Repositories.SubscriptionTemplate.Get(l.ctx, primary); err != nil {
		return err
	}
	for _, id := range available {
		if id == 0 {
			return repository.ErrInvalidArgument
		}
		if _, err := l.svcCtx.Repositories.SubscriptionTemplate.Get(l.ctx, id); err != nil {
			return err
		}
	}
	return nil
}
