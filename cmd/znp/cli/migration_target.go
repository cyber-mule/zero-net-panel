package cli

import (
	"strconv"
	"strings"
)

func parseMigrationTarget(value string) (uint64, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || strings.EqualFold(trimmed, "latest") {
		return 0, nil
	}
	return strconv.ParseUint(trimmed, 10, 64)
}
