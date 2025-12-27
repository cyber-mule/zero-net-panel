package account

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/logic/credentialutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// RotateCredentialLogic handles user credential rotation.
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

// Rotate rotates the current user's credential.
func (l *RotateCredentialLogic) Rotate(req *types.UserRotateCredentialRequest) (*types.UserRotateCredentialResponse, error) {
	_ = req

	user, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrForbidden
	}

	credential, err := credentialutil.RotateCredential(l.ctx, l.svcCtx.Repositories, l.svcCtx.Credentials, user.ID)
	if err != nil {
		return nil, err
	}

	return &types.UserRotateCredentialResponse{
		Credential: mapCredentialSummary(credential),
	}, nil
}
