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
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles admin user creation.
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic constructs CreateLogic.
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create provisions a new user.
func (l *CreateLogic) Create(req *types.AdminCreateUserRequest) (*types.AdminUserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" || !isValidEmail(email) {
		return nil, repository.ErrInvalidArgument
	}

	displayName := normalizeDisplayName("", email)
	if req.DisplayName != nil {
		displayName = normalizeDisplayName(*req.DisplayName, email)
	}

	roles := normalizeRoles(req.Roles)
	if len(roles) == 0 {
		roles = normalizeRoles(l.svcCtx.Config.Auth.Registration.DefaultRoles)
		if len(roles) == 0 {
			roles = []string{"user"}
		}
	}

	statusCode := status.UserStatusActive
	if req.Status != nil {
		normalized, err := normalizeStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		statusCode = normalized
	}

	now := time.Now().UTC()
	verifiedAt := repository.ZeroTime()
	if req.EmailVerified != nil {
		if *req.EmailVerified {
			verifiedAt = now
		}
	} else if statusCode == status.UserStatusActive {
		verifiedAt = now
	}

	if err := security.ValidatePasswordPolicy(password, l.svcCtx.Config.Auth.PasswordPolicy); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var created repository.User
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		user, err := txRepos.User.Create(l.ctx, repository.User{
			Email:           email,
			DisplayName:     displayName,
			PasswordHash:    string(passwordHash),
			Roles:           roles,
			Status:          statusCode,
			EmailVerifiedAt: verifiedAt,
		})
		if err != nil {
			return err
		}
		created = user

		var actorID *uint64
		actor, ok := security.UserFromContext(l.ctx)
		if ok && actor.ID != 0 {
			actorID = &actor.ID
		}

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      actorID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.user.create",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", user.ID),
			Metadata: map[string]any{
				"status": statusCode,
				"roles":  roles,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminUserResponse{User: toAdminUserSummary(created)}, nil
}
