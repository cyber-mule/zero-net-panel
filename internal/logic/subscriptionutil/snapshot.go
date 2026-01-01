package subscriptionutil

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

// BuildPlanSnapshot freezes key plan attributes for a subscription.
func BuildPlanSnapshot(plan repository.Plan, bindingIDs []uint64) map[string]any {
	snapshot := map[string]any{
		"id":                  plan.ID,
		"name":                plan.Name,
		"slug":                plan.Slug,
		"description":         plan.Description,
		"price_cents":         plan.PriceCents,
		"currency":            plan.Currency,
		"duration_days":       plan.DurationDays,
		"traffic_limit_bytes": plan.TrafficLimitBytes,
		"traffic_multipliers": cloneTrafficMultipliers(plan.TrafficMultipliers),
		"devices_limit":       plan.DevicesLimit,
		"features":            append([]string(nil), plan.Features...),
		"tags":                append([]string(nil), plan.Tags...),
	}
	if len(bindingIDs) > 0 {
		snapshot["binding_ids"] = append([]uint64(nil), bindingIDs...)
	}
	return snapshot
}

// ClonePlanSnapshot returns a shallow clone for safe reuse.
func ClonePlanSnapshot(snapshot map[string]any) map[string]any {
	if snapshot == nil {
		return nil
	}
	clone := make(map[string]any, len(snapshot))
	for key, value := range snapshot {
		clone[key] = value
	}
	return clone
}

// ExtractBindingIDs reads binding_ids from a snapshot payload.
func ExtractBindingIDs(snapshot map[string]any) []uint64 {
	if snapshot == nil {
		return nil
	}
	raw, ok := snapshot["binding_ids"]
	if !ok {
		return nil
	}
	switch value := raw.(type) {
	case []uint64:
		return append([]uint64(nil), value...)
	case []int:
		result := make([]uint64, 0, len(value))
		for _, item := range value {
			if item > 0 {
				result = append(result, uint64(item))
			}
		}
		return result
	case []any:
		result := make([]uint64, 0, len(value))
		for _, item := range value {
			switch v := item.(type) {
			case int:
				if v > 0 {
					result = append(result, uint64(v))
				}
			case int64:
				if v > 0 {
					result = append(result, uint64(v))
				}
			case float64:
				if v > 0 {
					result = append(result, uint64(v))
				}
			case uint64:
				if v > 0 {
					result = append(result, v)
				}
			}
		}
		return result
	default:
		return nil
	}
}

func cloneTrafficMultipliers(input map[string]float64) map[string]float64 {
	if input == nil {
		return map[string]float64{}
	}
	result := make(map[string]float64, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}
