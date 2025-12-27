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

// ForceLogoutLogic handles forced logout.
type ForceLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewForceLogoutLogic constructs ForceLogoutLogic.
func NewForceLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutLogic {
	return &ForceLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Force invalidates existing tokens.
func (l *ForceLogoutLogic) Force(req *types.AdminForceLogoutRequest) (*types.AdminForceLogoutResponse, error) {
	now := time.Now().UTC()
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		if _, err := txRepos.User.Get(l.ctx, req.UserID); err != nil {
			return err
		}
		if err := txRepos.User.UpdateTokenInvalidBefore(l.ctx, req.UserID, now); err != nil {
			return err
		}

		var actorID *uint64
		actor, ok := security.UserFromContext(l.ctx)
		if ok && actor.ID != 0 {
			actorID = &actor.ID
		}
		_, err := txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.user.force_logout",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", req.UserID),
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminForceLogoutResponse{Message: "logout forced"}, nil
}
