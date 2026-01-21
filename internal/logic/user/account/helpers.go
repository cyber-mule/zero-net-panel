package account

import (
	"net/mail"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapUserProfile(user repository.User) types.UserProfile {
	return types.UserProfile{
		ID:              user.ID,
		Email:           user.Email,
		DisplayName:     user.DisplayName,
		Status:          user.Status,
		EmailVerifiedAt: toUnixPtr(user.EmailVerifiedAt),
		CreatedAt:       toUnixOrZero(user.CreatedAt),
		UpdatedAt:       toUnixOrZero(user.UpdatedAt),
	}
}

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
	if repository.IsZeroTime(ts) {
		return nil
	}
	value := ts.Unix()
	return &value
}

func normalizeEmailInput(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}
