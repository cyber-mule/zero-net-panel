package paymentchannels

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic 管理端创建支付通道。
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic 构造函数。
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create 创建支付通道。
func (l *CreateLogic) Create(req *types.AdminCreatePaymentChannelRequest) (*types.PaymentChannelSummary, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" || code == "" {
		return nil, repository.ErrInvalidArgument
	}

	channel := repository.PaymentChannel{
		Name:      name,
		Code:      code,
		Provider:  strings.TrimSpace(req.Provider),
		Enabled:   req.Enabled,
		SortOrder: req.SortOrder,
		Config:    req.Config,
	}

	created, err := l.svcCtx.Repositories.PaymentChannel.Create(l.ctx, channel)
	if err != nil {
		return nil, err
	}

	summary := toPaymentChannelSummary(created)
	return &summary, nil
}
