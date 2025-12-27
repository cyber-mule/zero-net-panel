package subscriptions

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

func normalizeStatus(status string) (string, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		return "", repository.ErrInvalidArgument
	}
	switch status {
	case "active", "disabled", "expired", "pending":
		return status, nil
	default:
		return "", repository.ErrInvalidArgument
	}
}

func normalizeOptionalStatus(status *string) (*string, error) {
	if status == nil {
		return nil, nil
	}
	normalized, err := normalizeStatus(*status)
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
