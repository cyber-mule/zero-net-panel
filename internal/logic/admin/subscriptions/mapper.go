package subscriptions

import (
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toAdminSubscriptionSummary(sub repository.Subscription, user repository.User) types.AdminSubscriptionSummary {
	return types.AdminSubscriptionSummary{
		ID: sub.ID,
		User: types.AdminSubscriptionUserSummary{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
		Name:                 sub.Name,
		PlanName:             sub.PlanName,
		PlanID:               sub.PlanID,
		PlanSnapshot:         sub.PlanSnapshot,
		Status:               sub.Status,
		TemplateID:           sub.TemplateID,
		AvailableTemplateIDs: append([]uint64(nil), sub.AvailableTemplateIDs...),
		Token:                sub.Token,
		ExpiresAt:            toUnixOrZero(sub.ExpiresAt),
		TrafficTotalBytes:    sub.TrafficTotalBytes,
		TrafficUsedBytes:     sub.TrafficUsedBytes,
		DevicesLimit:         sub.DevicesLimit,
		LastRefreshedAt:      toUnixOrZero(sub.LastRefreshedAt),
		CreatedAt:            toUnixOrZero(sub.CreatedAt),
		UpdatedAt:            toUnixOrZero(sub.UpdatedAt),
	}
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}
