package users

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toAdminUserSummary(user repository.User) types.AdminUserSummary {
	summary := types.AdminUserSummary{
		ID:                  user.ID,
		Email:               user.Email,
		DisplayName:         user.DisplayName,
		Roles:               append([]string(nil), user.Roles...),
		Status:              user.Status,
		FailedLoginAttempts: user.FailedLoginAttempts,
		CreatedAt:           user.CreatedAt.Unix(),
		UpdatedAt:           user.UpdatedAt.Unix(),
	}

	if !repository.IsZeroTime(user.EmailVerifiedAt) {
		ts := user.EmailVerifiedAt.Unix()
		summary.EmailVerifiedAt = &ts
	}
	if !repository.IsZeroTime(user.LockedUntil) {
		ts := user.LockedUntil.Unix()
		summary.LockedUntil = &ts
	}
	if !repository.IsZeroTime(user.LastLoginAt) {
		ts := user.LastLoginAt.Unix()
		summary.LastLoginAt = &ts
	}

	return summary
}
