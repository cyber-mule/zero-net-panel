package subscription

import (
	"strings"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toUserSummary(sub repository.Subscription, subscriptionBase string) types.UserSubscriptionSummary {
	summary := types.UserSubscriptionSummary{
		ID:                   sub.ID,
		Name:                 sub.Name,
		PlanName:             sub.PlanName,
		PlanID:               sub.PlanID,
		Status:               sub.Status,
		TemplateID:           sub.TemplateID,
		Token:                sub.Token,
		SubscriptionURL:      buildSubscriptionURL(subscriptionBase, sub.Token),
		AvailableTemplateIDs: append([]uint64(nil), sub.AvailableTemplateIDs...),
		ExpiresAt:            sub.ExpiresAt.Unix(),
		TrafficTotalBytes:    sub.TrafficTotalBytes,
		TrafficUsedBytes:     sub.TrafficUsedBytes,
		DevicesLimit:         sub.DevicesLimit,
		LastRefreshedAt:      sub.LastRefreshedAt.Unix(),
	}
	return summary
}

func buildSubscriptionURL(base, token string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	token = strings.TrimSpace(token)
	if base == "" || token == "" {
		return ""
	}
	return base + "/api/v1/subscriptions/" + token
}
