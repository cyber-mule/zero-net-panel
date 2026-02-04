package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// RegisterLogic handles account registration.
type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic constructs RegisterLogic.
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register creates a new user and optionally issues tokens.
func (l *RegisterLogic) Register(req *types.AuthRegisterRequest) (*types.AuthRegisterResponse, error) {
	email := normalizeEmailInput(req.Email)
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" || !isValidEmail(email) {
		return nil, repository.NewInvalidArgument("email or password is invalid")
	}

	authCfg := l.svcCtx.Config.Auth
	if !authCfg.Registration.Enabled {
		return nil, repository.ErrForbidden
	}

	if authCfg.Registration.InviteOnly {
		inviteCode := normalizeInviteCode(req.InviteCode)
		if inviteCode == "" {
			return nil, repository.ErrInviteCodeRequired
		}
		if !inviteAllowed(inviteCode, authCfg.Registration.InviteCodes) {
			return nil, repository.ErrInviteCodeInvalid
		}
	}

	if err := security.ValidatePasswordPolicy(password, authCfg.PasswordPolicy); err != nil {
		return nil, err
	}

	existing, err := l.svcCtx.Repositories.User.GetByEmail(l.ctx, email)
	if err == nil {
		if existing.Status == status.UserStatusPending && repository.IsZeroTime(existing.EmailVerifiedAt) {
			if err := l.sendVerificationCode(email, authCfg.Verification); err != nil {
				return nil, err
			}
			return &types.AuthRegisterResponse{
				RequiresVerification: true,
				User:                 toAuthenticatedUser(existing),
			}, nil
		}
		return nil, repository.ErrConflict
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	displayName := normalizeDisplayName("", email)
	if req.DisplayName != nil {
		displayName = normalizeDisplayName(*req.DisplayName, email)
	}

	roles := normalizeRolesInput(authCfg.Registration.DefaultRoles)
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	requiresVerification := authCfg.Registration.RequireEmailVerification
	statusCode := status.UserStatusActive
	verifiedAt := repository.ZeroTime()
	if requiresVerification {
		statusCode = status.UserStatusPending
	} else {
		verifiedAt = time.Now().UTC()
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := repository.User{
		Email:           email,
		DisplayName:     displayName,
		PasswordHash:    string(passwordHash),
		Roles:           roles,
		Status:          statusCode,
		EmailVerifiedAt: verifiedAt,
	}

	created, err := l.svcCtx.Repositories.User.Create(l.ctx, user)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.Repositories.AuditLog.Create(l.ctx, repository.AuditLog{
		ActorEmail:   email,
		Action:       "auth.register",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", created.ID),
		Metadata: map[string]any{
			"status": statusCode,
		},
	}); err != nil {
		return nil, err
	}

	if requiresVerification {
		if err := l.sendVerificationCode(email, authCfg.Verification); err != nil {
			return nil, err
		}
		return &types.AuthRegisterResponse{
			RequiresVerification: true,
			User:                 toAuthenticatedUser(created),
		}, nil
	}

	resp, err := l.issueTokens(created)
	if err != nil {
		return nil, err
	}
	resp.RequiresVerification = false
	return resp, nil
}

func (l *RegisterLogic) sendVerificationCode(email string, cfg config.AuthVerificationConfig) error {
	policy := normalizeCodePolicy(codePolicy{
		CodeLength:       cfg.CodeLength,
		CodeTTL:          cfg.CodeTTL,
		SendCooldown:     cfg.SendCooldown,
		SendLimitPerHour: cfg.SendLimitPerHour,
	})
	code, err := issueAuthCode(l.ctx, l.svcCtx.Cache, policy, codePurposeVerify, email)
	if err != nil {
		return err
	}

	projectName := strings.TrimSpace(l.svcCtx.Config.Project.Name)
	if projectName == "" {
		projectName = "ZNP"
	}

	subject := fmt.Sprintf("%s verification code", projectName)
	body := fmt.Sprintf("Your verification code is %s. It expires in %s.", code, policy.CodeTTL)
	return sendAuthEmail(l.Logger, l.svcCtx.Config.Auth.Email, email, subject, body)
}

func (l *RegisterLogic) issueTokens(user repository.User) (*types.AuthRegisterResponse, error) {
	audience := l.svcCtx.Config.Project.Name
	if audience == "" {
		audience = "znp"
	}

	pair, err := l.svcCtx.Auth.GenerateTokenPair(fmt.Sprintf("%d", user.ID), user.Roles, audience)
	if err != nil {
		return nil, err
	}

	return &types.AuthRegisterResponse{
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        computeTTL(pair.AccessExpire),
		RefreshExpiresIn: computeTTL(pair.RefreshExpire),
		User:             toAuthenticatedUser(user),
	}, nil
}

func normalizeInviteCode(code *string) string {
	if code == nil {
		return ""
	}
	return strings.TrimSpace(*code)
}

func inviteAllowed(code string, allowed []string) bool {
	if len(allowed) == 0 {
		return false
	}
	if code == "" {
		return false
	}
	for _, item := range allowed {
		if strings.EqualFold(strings.TrimSpace(item), code) {
			return true
		}
	}
	return false
}
