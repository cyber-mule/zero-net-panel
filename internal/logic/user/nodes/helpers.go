package nodes

import (
	"sort"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapUserNodeStatus(node repository.Node, kernels []repository.NodeKernel, bindings []repository.ProtocolBinding) types.UserNodeStatusSummary {
	statuses := make([]types.UserNodeKernelStatusSummary, 0, len(kernels))
	for _, kernel := range kernels {
		statuses = append(statuses, types.UserNodeKernelStatusSummary{
			Protocol:     kernel.Protocol,
			Status:       kernel.Status,
			LastSyncedAt: toUnixOrZero(kernel.LastSyncedAt),
		})
	}

	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Protocol < statuses[j].Protocol
	})

	protocolStatuses := make([]types.UserNodeProtocolStatusSummary, 0, len(bindings))
	for _, binding := range bindings {
		protocolStatuses = append(protocolStatuses, types.UserNodeProtocolStatusSummary{
			BindingID:       binding.ID,
			Protocol:        binding.Protocol,
			Role:            binding.Role,
			Status:          binding.Status,
			HealthStatus:    binding.HealthStatus,
			LastHeartbeatAt: toUnixOrZero(binding.LastHeartbeatAt),
		})
	}

	sort.Slice(protocolStatuses, func(i, j int) bool {
		if protocolStatuses[i].Protocol == protocolStatuses[j].Protocol {
			return protocolStatuses[i].BindingID < protocolStatuses[j].BindingID
		}
		return protocolStatuses[i].Protocol < protocolStatuses[j].Protocol
	})

	return types.UserNodeStatusSummary{
		ID:               node.ID,
		Name:             node.Name,
		Region:           node.Region,
		Country:          node.Country,
		ISP:              node.ISP,
		Status:           node.Status,
		Tags:             append([]string(nil), node.Tags...),
		Protocols:        append([]string(nil), node.Protocols...),
		CapacityMbps:     node.CapacityMbps,
		Description:      node.Description,
		LastSyncedAt:     toUnixOrZero(node.LastSyncedAt),
		UpdatedAt:        toUnixOrZero(node.UpdatedAt),
		KernelStatuses:   statuses,
		ProtocolStatuses: protocolStatuses,
	}
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

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}
