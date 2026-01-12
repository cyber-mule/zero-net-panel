package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingLogic) Ping() (*types.PingResponse, error) {
	defaults := repository.SiteSettingDefaults{
		Name:                                 l.svcCtx.Config.Site.Name,
		LogoURL:                              l.svcCtx.Config.Site.LogoURL,
		KernelOfflineProbeMaxIntervalSeconds: int(l.svcCtx.Config.Kernel.OfflineProbeMaxInterval / time.Second),
	}
	setting, err := l.svcCtx.Repositories.Site.GetSiteSetting(l.ctx, defaults)
	if err != nil {
		return nil, err
	}

	resp := &types.PingResponse{
		Status:    "ok",
		Service:   l.svcCtx.Config.Project.Name,
		Version:   l.svcCtx.Config.Project.Version,
		SiteName:  setting.Name,
		LogoURL:   setting.LogoURL,
		Timestamp: time.Now().Unix(),
	}

	return resp, nil
}
