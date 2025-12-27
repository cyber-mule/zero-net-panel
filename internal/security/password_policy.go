package security

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

// ValidatePasswordPolicy checks password complexity against configured policy.
func ValidatePasswordPolicy(password string, policy config.AuthPasswordPolicyConfig) error {
	policy.Normalize()
	length := utf8.RuneCountInString(password)
	if length < policy.MinLength {
		return fmt.Errorf("password policy: minimum length is %d: %w", policy.MinLength, repository.ErrInvalidArgument)
	}
	if policy.MaxLength > 0 && length > policy.MaxLength {
		return fmt.Errorf("password policy: maximum length is %d: %w", policy.MaxLength, repository.ErrInvalidArgument)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if policy.RequireUpper && !hasUpper {
		return fmt.Errorf("password policy: require uppercase letters: %w", repository.ErrInvalidArgument)
	}
	if policy.RequireLower && !hasLower {
		return fmt.Errorf("password policy: require lowercase letters: %w", repository.ErrInvalidArgument)
	}
	if policy.RequireDigit && !hasDigit {
		return fmt.Errorf("password policy: require digits: %w", repository.ErrInvalidArgument)
	}
	if policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password policy: require special characters: %w", repository.ErrInvalidArgument)
	}

	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password policy: password cannot be blank: %w", repository.ErrInvalidArgument)
	}

	return nil
}
