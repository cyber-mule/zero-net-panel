package auth

import (
	"net/mail"
	"strings"
)

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

func normalizeRolesInput(roles []string) []string {
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
