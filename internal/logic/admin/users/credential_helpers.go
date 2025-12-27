package users

import (
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapCredentialSummary(credential repository.UserCredential) types.CredentialSummary {
	return types.CredentialSummary{
		Version:      credential.Version,
		Status:       credential.Status,
		IssuedAt:     toUnixOrZero(credential.IssuedAt),
		DeprecatedAt: toUnixPtr(credential.DeprecatedAt),
		RevokedAt:    toUnixPtr(credential.RevokedAt),
		LastSeenAt:   toUnixPtr(credential.LastSeenAt),
	}
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}

func toUnixPtr(ts time.Time) *int64 {
	if ts.IsZero() {
		return nil
	}
	value := ts.Unix()
	return &value
}
