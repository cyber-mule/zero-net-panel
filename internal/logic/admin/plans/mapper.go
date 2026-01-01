package plans

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toPlanSummary(plan repository.Plan, options []repository.PlanBillingOption, bindingIDs []uint64) types.PlanSummary {
	return types.PlanSummary{
		ID:                 plan.ID,
		Name:               plan.Name,
		Slug:               plan.Slug,
		Description:        plan.Description,
		Tags:               append([]string(nil), plan.Tags...),
		Features:           append([]string(nil), plan.Features...),
		BindingIDs:         append([]uint64(nil), bindingIDs...),
		BillingOptions:     toPlanBillingOptionSummaries(options),
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

func toPlanBillingOptionSummaries(options []repository.PlanBillingOption) []types.PlanBillingOptionSummary {
	if len(options) == 0 {
		return []types.PlanBillingOptionSummary{}
	}
	result := make([]types.PlanBillingOptionSummary, 0, len(options))
	for _, option := range options {
		result = append(result, types.PlanBillingOptionSummary{
			ID:            option.ID,
			PlanID:        option.PlanID,
			Name:          option.Name,
			DurationValue: option.DurationValue,
			DurationUnit:  option.DurationUnit,
			PriceCents:    option.PriceCents,
			Currency:      option.Currency,
			SortOrder:     option.SortOrder,
			Status:        option.Status,
			Visible:       option.Visible,
			CreatedAt:     option.CreatedAt.Unix(),
			UpdatedAt:     option.UpdatedAt.Unix(),
		})
	}
	return result
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
