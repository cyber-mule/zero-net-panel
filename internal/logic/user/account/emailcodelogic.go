package account

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	authlogic "github.com/zero-net-panel/zero-net-panel/internal/logic/auth"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// EmailChangeCodeLogic handles email change verification code delivery.
type EmailChangeCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewEmailChangeCodeLogic constructs EmailChangeCodeLogic.
func NewEmailChangeCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailChangeCodeLogic {
	return &EmailChangeCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Send issues a verification code for updating email.
func (l *EmailChangeCodeLogic) Send(req *types.UserEmailChangeCodeRequest) (*types.UserEmailChangeCodeResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok || actor.ID == 0 {
		return nil, repository.ErrUnauthorized
	}
	if req == nil {
		return nil, repository.NewInvalidArgument("invalid request")
	}

	email := normalizeEmailInput(req.Email)
	if email == "" || !isValidEmail(email) {
		return nil, repository.NewInvalidArgument("email is invalid")
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

	code, err := authlogic.IssueEmailChangeCode(l.ctx, l.svcCtx.Cache, l.svcCtx.Config.Auth.Verification, email)
	if err != nil {
		return nil, err
	}

	projectName := strings.TrimSpace(l.svcCtx.Config.Project.Name)
	if projectName == "" {
		projectName = "ZNP"
	}
	subject := fmt.Sprintf("%s email change code", projectName)
	body := fmt.Sprintf("Your email change code is %s. It expires in %s.", code, l.svcCtx.Config.Auth.Verification.CodeTTL)
	if err := authlogic.SendAuthEmail(l.Logger, l.svcCtx.Config.Auth.Email, email, subject, body); err != nil {
		return nil, err
	}

	_, _ = l.svcCtx.Repositories.AuditLog.Create(l.ctx, repository.AuditLog{
		ActorID:      &actor.ID,
		ActorEmail:   actor.Email,
		ActorRoles:   actor.Roles,
		Action:       "user.email.change_request",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", actor.ID),
		Metadata: map[string]any{
			"email": email,
		},
	})

	return &types.UserEmailChangeCodeResponse{Message: "verification code sent"}, nil
}
