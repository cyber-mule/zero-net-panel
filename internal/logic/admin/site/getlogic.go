package site

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// GetLogic 查询站点配置。
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

// Get 返回当前站点配置。
func (l *GetLogic) Get() (*types.AdminSiteSettingResponse, error) {
	defaults := repository.SiteSettingDefaults{
		Name:                                 l.svcCtx.Config.Site.Name,
		LogoURL:                              l.svcCtx.Config.Site.LogoURL,
		KernelOfflineProbeMaxIntervalSeconds: int(l.svcCtx.Config.Kernel.OfflineProbeMaxInterval / time.Second),
	}
	setting, err := l.svcCtx.Repositories.Site.GetSiteSetting(l.ctx, defaults)
	if err != nil {
		return nil, err
	}

	resp := &types.AdminSiteSettingResponse{
		Setting: toSiteSetting(setting),
	}
	return resp, nil
}
