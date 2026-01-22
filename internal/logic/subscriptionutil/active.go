package subscriptionutil

import (
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

// IsSubscriptionEffective reports whether a subscription is active and not expired.
func IsSubscriptionEffective(sub repository.Subscription, now time.Time) bool {
	if sub.Status != status.SubscriptionStatusActive {
		return false
	}
	if sub.ExpiresAt.IsZero() {
		return true
	}
	return sub.ExpiresAt.After(now)
}
