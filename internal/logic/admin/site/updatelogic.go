package site

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic 更新站点配置。
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

// Update 根据请求更新站点配置。
func (l *UpdateLogic) Update(req *types.AdminUpdateSiteSettingRequest) (*types.AdminSiteSettingResponse, error) {
	defaults := repository.SiteSettingDefaults{
		Name:                                 l.svcCtx.Config.Site.Name,
		LogoURL:                              l.svcCtx.Config.Site.LogoURL,
		KernelOfflineProbeMaxIntervalSeconds: int(l.svcCtx.Config.Kernel.OfflineProbeMaxInterval / time.Second),
	}
	setting, err := l.svcCtx.Repositories.Site.GetSiteSetting(l.ctx, defaults)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		setting.Name = strings.TrimSpace(*req.Name)
	}
	if req.LogoURL != nil {
		setting.LogoURL = strings.TrimSpace(*req.LogoURL)
	}
	if req.AccessDomain != nil {
		setting.AccessDomain = strings.TrimSpace(*req.AccessDomain)
	}
	if req.KernelOfflineProbeMaxIntervalSeconds != nil {
		value := *req.KernelOfflineProbeMaxIntervalSeconds
		if value < 0 {
			value = 0
		}
		setting.KernelOfflineProbeMaxIntervalSeconds = value
	}

	updated, err := l.svcCtx.Repositories.Site.UpsertSiteSetting(l.ctx, setting)
	if err != nil {
		return nil, err
	}

	resp := &types.AdminSiteSettingResponse{
		Setting: toSiteSetting(updated),
	}
	return resp, nil
}
