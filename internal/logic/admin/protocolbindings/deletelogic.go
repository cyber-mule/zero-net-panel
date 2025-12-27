package protocolbindings

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DeleteLogic handles binding deletion.
type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteLogic constructs DeleteLogic.
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Delete removes a protocol binding.
func (l *DeleteLogic) Delete(req *types.AdminDeleteProtocolBindingRequest) error {
	return l.svcCtx.Repositories.ProtocolBinding.Delete(l.ctx, req.BindingID)
}
