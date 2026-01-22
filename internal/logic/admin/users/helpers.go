package users

import (
	"net/mail"
	"strings"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

func normalizeRoles(roles []string) []string {
	seen := make(map[string]struct{}, len(roles))
	var normalized []string
	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role == "" {
			continue
		}
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		normalized = append(normalized, role)
	}
	return normalized
}

func normalizeStatus(statusCode int) (int, error) {
	switch statusCode {
	case status.UserStatusActive, status.UserStatusDisabled, status.UserStatusPending:
		return statusCode, nil
	case 0:
		return 0, repository.ErrInvalidArgument
	default:
		return 0, repository.ErrInvalidArgument
	}
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func normalizeDisplayName(name, email string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		return name
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return email
}
