package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ForgotLogic handles password reset code requests.
type ForgotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewForgotLogic constructs ForgotLogic.
func NewForgotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgotLogic {
	return &ForgotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Forgot triggers a password reset code email.
func (l *ForgotLogic) Forgot(req *types.AuthForgotPasswordRequest) (*types.AuthForgotPasswordResponse, error) {
	email := normalizeEmailInput(req.Email)
	if email == "" || !isValidEmail(email) {
		return nil, repository.ErrInvalidArgument
	}

	user, err := l.svcCtx.Repositories.User.GetByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &types.AuthForgotPasswordResponse{Message: "if the account exists, a reset code has been sent"}, nil
		}
		return nil, err
	}
	if !strings.EqualFold(user.Status, "active") {
		return &types.AuthForgotPasswordResponse{Message: "if the account exists, a reset code has been sent"}, nil
	}

	if err := l.sendResetCode(email, l.svcCtx.Config.Auth.PasswordReset); err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.Repositories.AuditLog.Create(l.ctx, repository.AuditLog{
		ActorEmail:   email,
		Action:       "auth.forgot",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", user.ID),
	}); err != nil {
		return nil, err
	}

	return &types.AuthForgotPasswordResponse{Message: "if the account exists, a reset code has been sent"}, nil
}

func (l *ForgotLogic) sendResetCode(email string, cfg config.AuthPasswordResetConfig) error {
	policy := normalizeCodePolicy(codePolicy{
		CodeLength:       cfg.CodeLength,
		CodeTTL:          cfg.CodeTTL,
		SendCooldown:     cfg.SendCooldown,
		SendLimitPerHour: cfg.SendLimitPerHour,
	})
	code, err := issueAuthCode(l.ctx, l.svcCtx.Cache, policy, codePurposeReset, email)
	if err != nil {
		return err
	}

	projectName := strings.TrimSpace(l.svcCtx.Config.Project.Name)
	if projectName == "" {
		projectName = "ZNP"
	}
	subject := fmt.Sprintf("%s password reset", projectName)
	body := fmt.Sprintf("Your reset code is %s. It expires in %s.", code, policy.CodeTTL)
	return sendAuthEmail(l.Logger, l.svcCtx.Config.Auth.Email, email, subject, body)
}
