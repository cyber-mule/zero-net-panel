package users

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/logic/credentialutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// RotateCredentialLogic handles admin credential rotation.
type RotateCredentialLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRotateCredentialLogic constructs RotateCredentialLogic.
func NewRotateCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RotateCredentialLogic {
	return &RotateCredentialLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Rotate rotates a user's credential.
func (l *RotateCredentialLogic) Rotate(req *types.AdminRotateUserCredentialRequest) (*types.AdminRotateUserCredentialResponse, error) {
	if req.UserID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	if _, err := l.svcCtx.Repositories.User.Get(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	credential, err := credentialutil.RotateCredential(l.ctx, l.svcCtx.Repositories, l.svcCtx.Credentials, req.UserID)
	if err != nil {
		return nil, err
	}

	return &types.AdminRotateUserCredentialResponse{
		UserID:     req.UserID,
		Credential: mapCredentialSummary(credential),
	}, nil
}
