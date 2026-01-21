package protocolbindings

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/nodecfg"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

var allowedRoles = map[string]struct{}{
	"listener":  {},
	"connector": {},
}

func normalizeRole(role string) (string, bool) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return "", false
	}
	_, ok := allowedRoles[role]
	return role, ok
}

func mapProtocolBindingSummary(binding repository.ProtocolBinding) types.ProtocolBindingSummary {
	return types.ProtocolBindingSummary{
		ID:              binding.ID,
		Name:            binding.Name,
		NodeID:          binding.NodeID,
		NodeName:        binding.Node.Name,
		Protocol:        normalizeBindingProtocol(binding),
		Role:            binding.Role,
		Listen:          binding.Listen,
		Connect:         binding.Connect,
		AccessPort:      binding.AccessPort,
		Status:          binding.Status,
		KernelID:        binding.KernelID,
		SyncStatus:      binding.SyncStatus,
		HealthStatus:    binding.HealthStatus,
		LastSyncedAt:    toUnixOrZero(binding.LastSyncedAt),
		LastHeartbeatAt: toUnixOrZero(binding.LastHeartbeatAt),
		LastSyncError:   binding.LastSyncError,
		Tags:            append([]string(nil), binding.Tags...),
		Description:     binding.Description,
		Profile:         cloneBindingProfile(binding.Profile),
		Metadata:        binding.Metadata,
		CreatedAt:       toUnixOrZero(binding.CreatedAt),
		UpdatedAt:       toUnixOrZero(binding.UpdatedAt),
	}
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}

func extractHostPort(address string) (string, int) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", 0
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return address, 0
	}
	portNum, _ := strconv.Atoi(port)
	return host, portNum
}

func normalizeBindingProtocol(binding repository.ProtocolBinding) string {
	return strings.ToLower(strings.TrimSpace(binding.Protocol))
}

func resolveKernelHTTPTimeout(node repository.Node) time.Duration {
	timeoutSeconds := node.KernelHTTPTimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = nodecfg.DefaultKernelHTTPTimeoutSeconds
	}
	return time.Duration(timeoutSeconds) * time.Second
}
