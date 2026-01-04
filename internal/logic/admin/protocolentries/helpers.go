package protocolentries

import (
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapProtocolEntrySummary(entry repository.ProtocolEntry) types.ProtocolEntrySummary {
	binding := entry.Binding
	node := binding.Node
	return types.ProtocolEntrySummary{
		ID:            entry.ID,
		Name:          entry.Name,
		BindingID:     binding.ID,
		BindingName:   binding.Name,
		NodeID:        binding.NodeID,
		NodeName:      node.Name,
		Protocol:      normalizeEntryProtocol(entry, binding),
		Status:        entry.Status,
		BindingStatus: binding.Status,
		HealthStatus:  binding.HealthStatus,
		EntryAddress:  entry.EntryAddress,
		EntryPort:     entry.EntryPort,
		Tags:          append([]string(nil), entry.Tags...),
		Description:   entry.Description,
		Profile:       cloneEntryProfile(entry.Profile),
		CreatedAt:     toUnixOrZero(entry.CreatedAt),
		UpdatedAt:     toUnixOrZero(entry.UpdatedAt),
	}
}

func normalizeEntryProtocol(entry repository.ProtocolEntry, binding repository.ProtocolBinding) string {
	if value := strings.ToLower(strings.TrimSpace(entry.Protocol)); value != "" {
		return value
	}
	return strings.ToLower(strings.TrimSpace(binding.Protocol))
}

func cloneEntryProfile(profile map[string]any) map[string]any {
	if profile == nil {
		return nil
	}
	cloned := make(map[string]any, len(profile))
	for key, value := range profile {
		cloned[key] = value
	}
	return cloned
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}
