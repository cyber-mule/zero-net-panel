package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ResetPasswordLogic handles admin password resets.
type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewResetPasswordLogic constructs ResetPasswordLogic.
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Reset updates the user password.
func (l *ResetPasswordLogic) Reset(req *types.AdminResetUserPasswordRequest) (*types.AdminResetUserPasswordResponse, error) {
	password := strings.TrimSpace(req.Password)
	if password == "" {
		return nil, repository.ErrInvalidArgument
	}

	if err := security.ValidatePasswordPolicy(password, l.svcCtx.Config.Auth.PasswordPolicy); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		if _, err := txRepos.User.Get(l.ctx, req.UserID); err != nil {
			return err
		}
		if err := txRepos.User.UpdatePassword(l.ctx, req.UserID, string(passwordHash), &now); err != nil {
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
			Action:       "admin.user.reset_password",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", req.UserID),
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminResetUserPasswordResponse{Message: "password reset"}, nil
}
