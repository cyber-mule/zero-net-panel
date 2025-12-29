package planbillingoptions

import (
	"fmt"
	"strings"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

func normalizeDurationUnit(unit string) string {
	unit = strings.TrimSpace(strings.ToLower(unit))
	switch unit {
	case "hours":
		return repository.DurationUnitHour
	case "days":
		return repository.DurationUnitDay
	case "months":
		return repository.DurationUnitMonth
	case "years":
		return repository.DurationUnitYear
	default:
		return unit
	}
}

func isValidDurationUnit(unit string) bool {
	switch normalizeDurationUnit(unit) {
	case repository.DurationUnitHour,
		repository.DurationUnitDay,
		repository.DurationUnitMonth,
		repository.DurationUnitYear:
		return true
	default:
		return false
	}
}

func formatDurationLabel(value int, unit string) string {
	if value <= 0 {
		return ""
	}
	normalized := normalizeDurationUnit(unit)
	label := normalized
	if value != 1 {
		label = fmt.Sprintf("%ss", normalized)
	}
	return fmt.Sprintf("%d %s", value, label)
}
