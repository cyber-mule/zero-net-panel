package protocolconfigs

import (
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

var supportedProtocols = map[string]struct{}{
	"ss":      {},
	"socks":   {},
	"http":    {},
	"vless":   {},
	"reality": {},
}

func normalizeProtocol(input string) (string, bool) {
	protocol := strings.ToLower(strings.TrimSpace(input))
	if protocol == "" {
		return "", false
	}
	_, ok := supportedProtocols[protocol]
	return protocol, ok
}

func mapProtocolConfigSummary(cfg repository.ProtocolConfig) types.ProtocolConfigSummary {
	return types.ProtocolConfigSummary{
		ID:          cfg.ID,
		Name:        cfg.Name,
		Protocol:    cfg.Protocol,
		Status:      cfg.Status,
		Tags:        append([]string(nil), cfg.Tags...),
		Description: cfg.Description,
		Profile:     cfg.Profile,
		CreatedAt:   toUnixOrZero(cfg.CreatedAt),
		UpdatedAt:   toUnixOrZero(cfg.UpdatedAt),
	}
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}
