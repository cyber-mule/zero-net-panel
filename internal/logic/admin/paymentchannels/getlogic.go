package paymentchannels

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// GetLogic 管理端查询支付通道。
type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetLogic 构造函数。
func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Get 返回支付通道详情。
func (l *GetLogic) Get(req *types.AdminGetPaymentChannelRequest) (*types.PaymentChannelSummary, error) {
	channel, err := l.svcCtx.Repositories.PaymentChannel.Get(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	summary := toPaymentChannelSummary(channel)
	return &summary, nil
}
