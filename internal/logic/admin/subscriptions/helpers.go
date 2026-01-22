package subscriptions

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

func normalizeStatus(statusCode int) (int, error) {
	if statusCode == 0 {
		return 0, repository.ErrInvalidArgument
	}
	switch statusCode {
	case status.SubscriptionStatusActive,
		status.SubscriptionStatusDisabled,
		status.SubscriptionStatusExpired:
		return statusCode, nil
	default:
		return 0, repository.ErrInvalidArgument
	}
}

func normalizeOptionalStatus(statusCode *int) (*int, error) {
	if statusCode == nil {
		return nil, nil
	}
	normalized, err := normalizeStatus(*statusCode)
	if err != nil {
		return nil, err
	}
	return &normalized, nil
}

func generateToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func validateExpiry(ts int64) (time.Time, error) {
	if ts <= 0 {
		return time.Time{}, repository.ErrInvalidArgument
	}
	return time.Unix(ts, 0).UTC(), nil
}

func validateExtendRequest(extendDays, extendHours int, expiresAt *int64) (time.Duration, *time.Time, error) {
	if expiresAt != nil {
		if *expiresAt <= 0 {
			return 0, nil, repository.ErrInvalidArgument
		}
		value := time.Unix(*expiresAt, 0).UTC()
		return 0, &value, nil
	}

	if extendDays == 0 && extendHours == 0 {
		return 0, nil, repository.ErrInvalidArgument
	}
	if extendDays < 0 || extendHours < 0 {
		return 0, nil, repository.ErrInvalidArgument
	}
	duration := time.Duration(extendDays)*24*time.Hour + time.Duration(extendHours)*time.Hour
	if duration <= 0 {
		return 0, nil, repository.ErrInvalidArgument
	}
	return duration, nil, nil
}
