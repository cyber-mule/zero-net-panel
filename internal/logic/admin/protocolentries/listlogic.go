package protocolentries

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles protocol entry listing.
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

// List returns entry summaries.
func (l *ListLogic) List(req *types.AdminListProtocolEntriesRequest) (*types.AdminProtocolEntryListResponse, error) {
	opts := repository.ListProtocolEntriesOptions{
		Page:      req.Page,
		PerPage:   req.PerPage,
		Sort:      req.Sort,
		Direction: req.Direction,
		Query:     req.Query,
		Status:    req.Status,
		Protocol:  req.Protocol,
	}
	if req.BindingID != nil && *req.BindingID > 0 {
		opts.BindingID = req.BindingID
	}

	entries, total, err := l.svcCtx.Repositories.ProtocolEntry.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	summaries := make([]types.ProtocolEntrySummary, 0, len(entries))
	for _, entry := range entries {
		summaries = append(summaries, mapProtocolEntrySummary(entry))
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.AdminProtocolEntryListResponse{
		Entries:    summaries,
		Pagination: pagination,
	}, nil
}

func normalizePage(page, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}
