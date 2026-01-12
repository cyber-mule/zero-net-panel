package site

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toSiteSetting(setting repository.SiteSetting) types.SiteSetting {
	return types.SiteSetting{
		ID:                                   setting.ID,
		Name:                                 setting.Name,
		LogoURL:                              setting.LogoURL,
		AccessDomain:                         setting.AccessDomain,
		KernelOfflineProbeMaxIntervalSeconds: setting.KernelOfflineProbeMaxIntervalSeconds,
		CreatedAt:                            setting.CreatedAt.Unix(),
		UpdatedAt:                            setting.UpdatedAt.Unix(),
	}
}
