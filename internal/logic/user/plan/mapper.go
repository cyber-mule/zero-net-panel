package plan

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toUserPlanSummary(plan repository.Plan, options []repository.PlanBillingOption) types.UserPlanSummary {
	return types.UserPlanSummary{
		ID:                plan.ID,
		Name:              plan.Name,
		Description:       plan.Description,
		Features:          append([]string(nil), plan.Features...),
		BillingOptions:    toPlanBillingOptionSummaries(options),
		PriceCents:        plan.PriceCents,
		Currency:          plan.Currency,
		DurationDays:      plan.DurationDays,
		TrafficLimitBytes: plan.TrafficLimitBytes,
		DevicesLimit:      plan.DevicesLimit,
		Tags:              append([]string(nil), plan.Tags...),
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
