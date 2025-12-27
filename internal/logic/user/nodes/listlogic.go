package nodes

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic 用户端节点状态列表。
type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListLogic 构造函数。
func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// List 返回用户可见节点状态。
func (l *ListLogic) List(req *types.UserNodeStatusListRequest) (*types.UserNodeStatusListResponse, error) {
	opts := repository.ListNodesOptions{
		Page:     req.Page,
		PerPage:  req.PerPage,
		Status:   req.Status,
		Protocol: req.Protocol,
	}

	nodes, total, err := l.svcCtx.Repositories.Node.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	nodeIDs := make([]uint64, 0, len(nodes))
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}

	bindingsByNode := make(map[uint64][]repository.ProtocolBinding, len(nodes))
	if len(nodeIDs) > 0 {
		bindings, err := l.svcCtx.Repositories.ProtocolBinding.ListByNodeIDs(l.ctx, nodeIDs)
		if err != nil {
			return nil, err
		}
		for _, binding := range bindings {
			bindingsByNode[binding.NodeID] = append(bindingsByNode[binding.NodeID], binding)
		}
	}

	result := make([]types.UserNodeStatusSummary, 0, len(nodes))
	for _, node := range nodes {
		kernels, err := l.svcCtx.Repositories.Node.GetKernels(l.ctx, node.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, mapUserNodeStatus(node, kernels, bindingsByNode[node.ID]))
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.UserNodeStatusListResponse{
		Nodes:      result,
		Pagination: pagination,
	}, nil
}
