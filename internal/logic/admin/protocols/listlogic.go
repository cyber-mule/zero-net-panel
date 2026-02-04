package protocols

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles protocol list queries.
type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListLogic constructs ListLogic.
func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// List returns distinct protocols configured in protocol bindings.
func (l *ListLogic) List() (*types.AdminProtocolListResponse, error) {
	protocols, err := l.svcCtx.Repositories.ProtocolBinding.ListProtocols(l.ctx)
	if err != nil {
		return nil, err
	}

	return &types.AdminProtocolListResponse{
		Protocols: protocols,
	}, nil
}
