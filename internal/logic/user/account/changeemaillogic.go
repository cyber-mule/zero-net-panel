package account

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"

	authlogic "github.com/zero-net-panel/zero-net-panel/internal/logic/auth"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ChangeEmailLogic handles user email updates.
type ChangeEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangeEmailLogic constructs ChangeEmailLogic.
func NewChangeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeEmailLogic {
	return &ChangeEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Change updates the user's email after verification.
func (l *ChangeEmailLogic) Change(req *types.UserChangeEmailRequest) (*types.UserChangeEmailResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok || actor.ID == 0 {
		return nil, repository.ErrUnauthorized
	}
	if req == nil {
		return nil, repository.ErrInvalidArgument
	}

	email := normalizeEmailInput(req.Email)
	code := strings.TrimSpace(req.Code)
	password := strings.TrimSpace(req.Password)
	if email == "" || code == "" || password == "" || !isValidEmail(email) {
		return nil, repository.ErrInvalidArgument
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	if user.Status != status.UserStatusActive {
		return nil, repository.ErrForbidden
	}
	if strings.EqualFold(user.Email, email) {
		return nil, fmt.Errorf("email unchanged: %w", repository.ErrConflict)
	}

	if existing, err := l.svcCtx.Repositories.User.GetByEmail(l.ctx, email); err == nil && existing.ID != actor.ID {
		return nil, fmt.Errorf("email already in use: %w", repository.ErrConflict)
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, repository.ErrUnauthorized
	}

	if err := authlogic.VerifyEmailChangeCode(l.ctx, l.svcCtx.Cache, email, code); err != nil {
		return nil, err
	}

	oldEmail := user.Email
	now := time.Now().UTC()
	var updated repository.User
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		user, err := txRepos.User.UpdateEmail(l.ctx, actor.ID, email, &now)
		if err != nil {
			return err
		}
		updated = user

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "user.email.change",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", actor.ID),
			Metadata: map[string]any{
				"from": oldEmail,
				"to":   email,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.UserChangeEmailResponse{Profile: mapUserProfile(updated)}, nil
}
