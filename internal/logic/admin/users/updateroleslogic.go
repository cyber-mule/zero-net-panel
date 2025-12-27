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

// UpdateRolesLogic handles role updates.
type UpdateRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateRolesLogic constructs UpdateRolesLogic.
func NewUpdateRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolesLogic {
	return &UpdateRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update updates user roles.
func (l *UpdateRolesLogic) Update(req *types.AdminUpdateUserRolesRequest) (*types.AdminUserResponse, error) {
	roles := normalizeRoles(req.Roles)
	if len(roles) == 0 {
		return nil, repository.ErrInvalidArgument
	}

	var updated repository.User
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		user, err := txRepos.User.UpdateRoles(l.ctx, req.UserID, roles)
		if err != nil {
			return err
		}
		updated = user

		now := time.Now().UTC()
		if err := txRepos.User.UpdateTokenInvalidBefore(l.ctx, req.UserID, now); err != nil {
			return err
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
			Action:       "admin.user.roles",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", req.UserID),
			Metadata: map[string]any{
				"roles": roles,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminUserResponse{User: toAdminUserSummary(updated)}, nil
}
