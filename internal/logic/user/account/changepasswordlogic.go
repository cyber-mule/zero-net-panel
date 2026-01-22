package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ChangePasswordLogic handles user password changes.
type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangePasswordLogic constructs ChangePasswordLogic.
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Change updates the password for the authenticated user.
func (l *ChangePasswordLogic) Change(req *types.UserChangePasswordRequest) (*types.UserChangePasswordResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok || actor.ID == 0 {
		return nil, repository.ErrUnauthorized
	}
	if req == nil {
		return nil, repository.ErrInvalidArgument
	}

	current := strings.TrimSpace(req.CurrentPassword)
	next := strings.TrimSpace(req.NewPassword)
	if current == "" || next == "" {
		return nil, repository.ErrInvalidArgument
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	if user.Status != status.UserStatusActive {
		return nil, repository.ErrForbidden
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(current)); err != nil {
		return nil, repository.ErrUnauthorized
	}

	if err := security.ValidatePasswordPolicy(next, l.svcCtx.Config.Auth.PasswordPolicy); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		if err := txRepos.User.UpdatePassword(l.ctx, actor.ID, string(passwordHash), nil); err != nil {
			return err
		}
		if err := txRepos.User.UpdateTokenInvalidBefore(l.ctx, actor.ID, now); err != nil {
			return err
		}
		_, err := txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "user.password.change",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", actor.ID),
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.UserChangePasswordResponse{Message: "password updated"}, nil
}
