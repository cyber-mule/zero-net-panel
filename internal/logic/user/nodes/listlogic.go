package nodes

import (
	"context"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	subscriptionutil "github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
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
	user, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrForbidden
	}

	sub, err := l.svcCtx.Repositories.Subscription.GetActiveByUser(l.ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return emptyNodeResponse(req.Page, req.PerPage), nil
		}
		return nil, err
	}

	bindings, err := subscriptionutil.LoadSubscriptionBindings(l.ctx, l.svcCtx.Repositories, sub)
	if err != nil {
		return nil, err
	}

	filterProtocol := strings.ToLower(strings.TrimSpace(req.Protocol))
	bindingsByNode := make(map[uint64][]repository.ProtocolBinding)
	for _, binding := range bindings {
		if filterProtocol != "" && !strings.EqualFold(binding.Protocol, filterProtocol) {
			continue
		}
		bindingsByNode[binding.NodeID] = append(bindingsByNode[binding.NodeID], binding)
	}

	if len(bindingsByNode) == 0 {
		return emptyNodeResponse(req.Page, req.PerPage), nil
	}

	opts := repository.ListNodesOptions{
		Page:     req.Page,
		PerPage:  req.PerPage,
		Status:   req.Status,
		Protocol: req.Protocol,
	}

	nodeIDs := make([]uint64, 0, len(bindingsByNode))
	for nodeID := range bindingsByNode {
		nodeIDs = append(nodeIDs, nodeID)
	}
	opts.NodeIDs = nodeIDs

	nodes, total, err := l.svcCtx.Repositories.Node.List(l.ctx, opts)
	if err != nil {
		return nil, err
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

func emptyNodeResponse(page, perPage int) *types.UserNodeStatusListResponse {
	page, perPage = normalizePage(page, perPage)
	return &types.UserNodeStatusListResponse{
		Nodes: []types.UserNodeStatusSummary{},
		Pagination: types.PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			TotalCount: 0,
			HasNext:    false,
			HasPrev:    page > 1,
		},
	}
}
