package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateProfileLogic handles user profile updates.
type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateProfileLogic constructs UpdateProfileLogic.
func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update applies display name changes.
func (l *UpdateProfileLogic) Update(req *types.UserUpdateProfileRequest) (*types.UserProfileResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok || actor.ID == 0 {
		return nil, repository.ErrUnauthorized
	}
	if req == nil || req.DisplayName == nil {
		return nil, repository.ErrInvalidArgument
	}

	displayName := strings.TrimSpace(*req.DisplayName)
	if displayName == "" {
		return nil, repository.ErrInvalidArgument
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(user.Status, "active") {
		return nil, repository.ErrForbidden
	}

	var updated repository.User
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		user, err := txRepos.User.UpdateProfile(l.ctx, actor.ID, repository.UpdateUserProfileInput{
			DisplayName: &displayName,
		})
		if err != nil {
			return err
		}
		updated = user

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "user.profile.update",
			ResourceType: "user",
			ResourceID:   fmt.Sprintf("%d", actor.ID),
			Metadata: map[string]any{
				"display_name": displayName,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.UserProfileResponse{Profile: mapUserProfile(updated)}, nil
}
