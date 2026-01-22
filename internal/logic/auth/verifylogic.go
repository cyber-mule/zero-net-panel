package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// VerifyLogic handles email verification.
type VerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewVerifyLogic constructs VerifyLogic.
func NewVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyLogic {
	return &VerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Verify validates verification code and activates the user.
func (l *VerifyLogic) Verify(req *types.AuthVerifyRequest) (*types.AuthVerifyResponse, error) {
	email := normalizeEmailInput(req.Email)
	code := strings.TrimSpace(req.Code)
	if email == "" || code == "" || !isValidEmail(email) {
		return nil, repository.ErrInvalidArgument
	}

	user, err := l.svcCtx.Repositories.User.GetByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrInvalidArgument
		}
		return nil, err
	}

	if !repository.IsZeroTime(user.EmailVerifiedAt) {
		return nil, repository.ErrConflict
	}

	if err := verifyAuthCode(l.ctx, l.svcCtx.Cache, codePurposeVerify, email, code); err != nil {
		return nil, err
	}

	statusCode := user.Status
	if statusCode == status.UserStatusPending || statusCode == 0 {
		statusCode = status.UserStatusActive
	}

	verifiedAt := time.Now().UTC()
	user, err = l.svcCtx.Repositories.User.UpdateVerification(l.ctx, user.ID, verifiedAt, statusCode)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.Repositories.AuditLog.Create(l.ctx, repository.AuditLog{
		ActorID:      &user.ID,
		ActorEmail:   user.Email,
		ActorRoles:   user.Roles,
		Action:       "auth.verify",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", user.ID),
	}); err != nil {
		return nil, err
	}

	if user.Status != status.UserStatusActive {
		return nil, repository.ErrForbidden
	}

	audience := l.svcCtx.Config.Project.Name
	if audience == "" {
		audience = "znp"
	}
	pair, err := l.svcCtx.Auth.GenerateTokenPair(fmt.Sprintf("%d", user.ID), user.Roles, audience)
	if err != nil {
		return nil, err
	}

	return &types.AuthVerifyResponse{
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        computeTTL(pair.AccessExpire),
		RefreshExpiresIn: computeTTL(pair.RefreshExpire),
		User:             toAuthenticatedUser(user),
	}, nil
}
