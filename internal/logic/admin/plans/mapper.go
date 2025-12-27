package plans

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toPlanSummary(plan repository.Plan) types.PlanSummary {
	return types.PlanSummary{
		ID:                 plan.ID,
		Name:               plan.Name,
		Slug:               plan.Slug,
		Description:        plan.Description,
		Tags:               append([]string(nil), plan.Tags...),
		Features:           append([]string(nil), plan.Features...),
		PriceCents:         plan.PriceCents,
		Currency:           plan.Currency,
		DurationDays:       plan.DurationDays,
		TrafficLimitBytes:  plan.TrafficLimitBytes,
		TrafficMultipliers: cloneTrafficMultipliers(plan.TrafficMultipliers),
		DevicesLimit:       plan.DevicesLimit,
		SortOrder:          plan.SortOrder,
		Status:             plan.Status,
		Visible:            plan.Visible,
		CreatedAt:          plan.CreatedAt.Unix(),
		UpdatedAt:          plan.UpdatedAt.Unix(),
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
