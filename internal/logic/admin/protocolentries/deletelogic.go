package protocolentries

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DeleteLogic handles protocol entry deletions.
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

// Delete deletes a protocol entry.
func (l *DeleteLogic) Delete(req *types.AdminDeleteProtocolEntryRequest) error {
	if req.EntryID == 0 {
		return repository.ErrInvalidArgument
	}
	return l.svcCtx.Repositories.ProtocolEntry.Delete(l.ctx, req.EntryID)
}
