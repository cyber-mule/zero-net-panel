package nodes

import (
	"sort"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapNodeSummary(node repository.Node) types.NodeSummary {
	return types.NodeSummary{
		ID:                node.ID,
		Name:              node.Name,
		Region:            node.Region,
		Country:           node.Country,
		ISP:               node.ISP,
		Status:            node.Status,
		Tags:              append([]string(nil), node.Tags...),
		CapacityMbps:      node.CapacityMbps,
		Description:       node.Description,
		AccessAddress:     node.AccessAddress,
		ControlEndpoint:   node.ControlEndpoint,
		StatusSyncEnabled: node.StatusSyncEnabled,
		LastSyncedAt:      toUnixOrZero(node.LastSyncedAt),
		UpdatedAt:         toUnixOrZero(node.UpdatedAt),
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

func normalizeNodeStatus(status string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "" {
		return "", repository.ErrInvalidArgument
	}
	switch normalized {
	case "online", "offline", "maintenance", "disabled":
		return normalized, nil
	default:
		return "", repository.ErrInvalidArgument
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
