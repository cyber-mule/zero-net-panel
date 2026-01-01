package subscriptionutil

import (
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

// IsSubscriptionEffective reports whether a subscription is active and not expired.
func IsSubscriptionEffective(sub repository.Subscription, now time.Time) bool {
	if !strings.EqualFold(strings.TrimSpace(sub.Status), "active") {
		return false
	}
	if sub.ExpiresAt.IsZero() {
		return true
	}
	return sub.ExpiresAt.After(now)
}
