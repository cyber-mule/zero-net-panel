package protocolconfigs

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DeleteLogic handles protocol config deletion.
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

// Delete removes a protocol config.
func (l *DeleteLogic) Delete(req *types.AdminDeleteProtocolConfigRequest) error {
	return l.svcCtx.Repositories.ProtocolConfig.Delete(l.ctx, req.ConfigID)
}
