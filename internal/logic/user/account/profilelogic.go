package account

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ProfileLogic handles user profile retrieval.
type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewProfileLogic constructs ProfileLogic.
func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Get returns the authenticated user's profile.
func (l *ProfileLogic) Get(_ *types.UserProfileRequest) (*types.UserProfileResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok || actor.ID == 0 {
		return nil, repository.ErrUnauthorized
	}

	user, err := l.svcCtx.Repositories.User.Get(l.ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(user.Status, "active") {
		return nil, repository.ErrForbidden
	}

	return &types.UserProfileResponse{Profile: mapUserProfile(user)}, nil
}
