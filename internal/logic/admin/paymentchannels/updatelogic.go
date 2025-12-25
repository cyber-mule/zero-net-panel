package paymentchannels

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic 管理端更新支付通道。
type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLogic 构造函数。
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update 更新支付通道配置。
func (l *UpdateLogic) Update(req *types.AdminUpdatePaymentChannelRequest) (*types.PaymentChannelSummary, error) {
	channel, err := l.svcCtx.Repositories.PaymentChannel.Get(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		channel.Name = strings.TrimSpace(*req.Name)
	}
	if req.Code != nil {
		channel.Code = strings.TrimSpace(*req.Code)
	}
	if req.Provider != nil {
		channel.Provider = strings.TrimSpace(*req.Provider)
	}
	if req.Enabled != nil {
		channel.Enabled = *req.Enabled
	}
	if req.SortOrder != nil {
		channel.SortOrder = *req.SortOrder
	}
	if req.Config != nil {
		channel.Config = req.Config
	}

	updated, err := l.svcCtx.Repositories.PaymentChannel.Update(l.ctx, channel.ID, channel)
	if err != nil {
		return nil, err
	}

	summary := toPaymentChannelSummary(updated)
	return &summary, nil
}
