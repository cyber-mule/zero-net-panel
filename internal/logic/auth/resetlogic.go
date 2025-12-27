package auth

import (
	"context"
	"errors"
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

// ResetLogic handles password resets.
type ResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewResetLogic constructs ResetLogic.
func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Reset verifies code and updates password.
func (l *ResetLogic) Reset(req *types.AuthResetPasswordRequest) (*types.AuthResetPasswordResponse, error) {
	email := normalizeEmailInput(req.Email)
	code := strings.TrimSpace(req.Code)
	password := strings.TrimSpace(req.Password)
	if email == "" || code == "" || password == "" || !isValidEmail(email) {
		return nil, repository.ErrInvalidArgument
	}

	user, err := l.svcCtx.Repositories.User.GetByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrInvalidArgument
		}
		return nil, err
	}

	if err := verifyAuthCode(l.ctx, l.svcCtx.Cache, codePurposeReset, email, code); err != nil {
		return nil, err
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
		if err := txRepos.User.UpdatePassword(l.ctx, user.ID, string(passwordHash), &now); err != nil {
			return err
		}
		if err := txRepos.User.UpdateTokenInvalidBefore(l.ctx, user.ID, now); err != nil {
			return err
		}
		_, err := txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &user.ID,
			ActorEmail:   user.Email,
			ActorRoles:   user.Roles,
			Action:       "auth.reset",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", user.ID),
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AuthResetPasswordResponse{Message: "password reset successful"}, nil
}
