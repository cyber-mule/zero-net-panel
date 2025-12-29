package planbillingoptions

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toPlanBillingOptionSummary(option repository.PlanBillingOption) types.PlanBillingOptionSummary {
	return types.PlanBillingOptionSummary{
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
	}
}
