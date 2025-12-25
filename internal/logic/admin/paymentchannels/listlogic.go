package paymentchannels

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic 管理端支付通道列表。
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

// List 返回支付通道列表。
func (l *ListLogic) List(req *types.AdminListPaymentChannelsRequest) (*types.AdminPaymentChannelListResponse, error) {
	opts := repository.ListPaymentChannelsOptions{
		Page:      req.Page,
		PerPage:   req.PerPage,
		Query:     req.Query,
		Provider:  req.Provider,
		Enabled:   req.Enabled,
		Sort:      req.Sort,
		Direction: req.Direction,
	}

	channels, total, err := l.svcCtx.Repositories.PaymentChannel.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make([]types.PaymentChannelSummary, 0, len(channels))
	for _, channel := range channels {
		result = append(result, toPaymentChannelSummary(channel))
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.AdminPaymentChannelListResponse{
		Channels:   result,
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
