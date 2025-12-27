package users

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateStatusLogic handles user status updates.
type UpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateStatusLogic constructs UpdateStatusLogic.
func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update updates the user status.
func (l *UpdateStatusLogic) Update(req *types.AdminUpdateUserStatusRequest) (*types.AdminUserResponse, error) {
	status := normalizeStatus(req.Status)
	if status != "active" && status != "disabled" && status != "pending" {
		return nil, repository.ErrInvalidArgument
	}

	var updated repository.User
	now := time.Now().UTC()
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		user, err := txRepos.User.UpdateStatus(l.ctx, req.UserID, status)
		if err != nil {
			return err
		}
		updated = user

		if status == "disabled" {
			if err := txRepos.User.UpdateTokenInvalidBefore(l.ctx, req.UserID, now); err != nil {
				return err
			}
		}

		var actorID *uint64
		actor, ok := security.UserFromContext(l.ctx)
		if ok && actor.ID != 0 {
			actorID = &actor.ID
		}
		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.user.status",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", req.UserID),
			Metadata: map[string]any{
				"status": status,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminUserResponse{User: toAdminUserSummary(updated)}, nil
}
