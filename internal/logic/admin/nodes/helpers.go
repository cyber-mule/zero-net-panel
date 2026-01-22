package nodes

import (
	"sort"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/nodecfg"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapNodeSummary(node repository.Node) types.NodeSummary {
	return types.NodeSummary{
		ID:                              node.ID,
		Name:                            node.Name,
		Region:                          node.Region,
		Country:                         node.Country,
		ISP:                             node.ISP,
		Status:                          node.Status,
		Tags:                            append([]string(nil), node.Tags...),
		CapacityMbps:                    node.CapacityMbps,
		Description:                     node.Description,
		AccessAddress:                   node.AccessAddress,
		ControlEndpoint:                 node.ControlEndpoint,
		KernelDefaultProtocol:           node.KernelDefaultProtocol,
		KernelHTTPTimeoutSeconds:        node.KernelHTTPTimeoutSeconds,
		KernelStatusPollIntervalSeconds: node.KernelStatusPollIntervalSeconds,
		KernelStatusPollBackoffEnabled:  node.KernelStatusPollBackoffEnabled,
		KernelStatusPollBackoffMaxIntervalSeconds: node.KernelStatusPollBackoffMaxIntervalSeconds,
		KernelStatusPollBackoffMultiplier:         node.KernelStatusPollBackoffMultiplier,
		KernelStatusPollBackoffJitter:             node.KernelStatusPollBackoffJitter,
		KernelOfflineProbeMaxIntervalSeconds:      node.KernelOfflineProbeMaxIntervalSeconds,
		StatusSyncEnabled:                         node.StatusSyncEnabled,
		LastSyncedAt:                              toUnixOrZero(node.LastSyncedAt),
		UpdatedAt:                                 toUnixOrZero(node.UpdatedAt),
	}
}

func mapKernelSummary(kernel repository.NodeKernel) types.NodeKernelSummary {
	return types.NodeKernelSummary{
		Protocol:     kernel.Protocol,
		Endpoint:     kernel.Endpoint,
		Revision:     kernel.Revision,
		Status:       kernel.Status,
		Config:       kernel.Config,
		LastSyncedAt: toUnixOrZero(kernel.LastSyncedAt),
	}
}

func normalizeNodeStatus(statusCode int) (int, error) {
	if statusCode == 0 {
		return 0, repository.ErrInvalidArgument
	}
	switch statusCode {
	case status.NodeStatusOnline,
		status.NodeStatusOffline,
		status.NodeStatusMaintenance,
		status.NodeStatusDisabled:
		return statusCode, nil
	default:
		return 0, repository.ErrInvalidArgument
	}
}

func normalizeKernelStatus(statusCode int) (int, error) {
	if statusCode == 0 {
		return 0, repository.ErrInvalidArgument
	}
	switch statusCode {
	case status.NodeKernelStatusConfigured,
		status.NodeKernelStatusSynced:
		return statusCode, nil
	default:
		return 0, repository.ErrInvalidArgument
	}
}

func normalizeTags(tags []string) []string {
	return normalizeStringSet(tags, false)
}

func normalizeStringSet(values []string, lower bool) []string {
	if values == nil {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" {
			continue
		}
		if lower {
			item = strings.ToLower(item)
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	sort.Strings(result)
	return result
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}

func resolveKernelHTTPTimeout(node repository.Node) time.Duration {
	timeoutSeconds := node.KernelHTTPTimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = nodecfg.DefaultKernelHTTPTimeoutSeconds
	}
	return time.Duration(timeoutSeconds) * time.Second
}
