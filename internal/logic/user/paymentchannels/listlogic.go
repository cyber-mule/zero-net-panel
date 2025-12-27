package paymentchannels

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles user payment channel listing.
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

// List returns enabled payment channels for users.
func (l *ListLogic) List(req *types.UserPaymentChannelListRequest) (*types.UserPaymentChannelListResponse, error) {
	enabled := true
	opts := repository.ListPaymentChannelsOptions{
		Page:     1,
		PerPage:  200,
		Provider: req.Provider,
		Enabled:  &enabled,
	}

	channels, _, err := l.svcCtx.Repositories.PaymentChannel.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make([]types.UserPaymentChannelSummary, 0, len(channels))
	for _, channel := range channels {
		result = append(result, toUserPaymentChannelSummary(channel))
	}

	return &types.UserPaymentChannelListResponse{Channels: result}, nil
}

func toUserPaymentChannelSummary(channel repository.PaymentChannel) types.UserPaymentChannelSummary {
	return types.UserPaymentChannelSummary{
		ID:        channel.ID,
		Name:      channel.Name,
		Code:      channel.Code,
		Provider:  channel.Provider,
		SortOrder: channel.SortOrder,
		Config:    channel.Config,
	}
}
