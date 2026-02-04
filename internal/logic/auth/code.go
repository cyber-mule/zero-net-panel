package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/pkg/cache"
)

const (
	codePurposeVerify      = "verify"
	codePurposeReset       = "reset"
	codePurposeEmailChange = "email_change"
)

type codePolicy struct {
	CodeLength       int
	CodeTTL          time.Duration
	SendCooldown     time.Duration
	SendLimitPerHour int
}

type codeRecord struct {
	Hash   string `json:"hash"`
	SentAt int64  `json:"sent_at"`
}

type counterRecord struct {
	Count int `json:"count"`
}

// IssueEmailChangeCode issues a verification code for updating user email.
func IssueEmailChangeCode(ctx context.Context, c cache.Cache, cfg config.AuthVerificationConfig, email string) (string, error) {
	policy := normalizeCodePolicy(codePolicy{
		CodeLength:       cfg.CodeLength,
		CodeTTL:          cfg.CodeTTL,
		SendCooldown:     cfg.SendCooldown,
		SendLimitPerHour: cfg.SendLimitPerHour,
	})
	return issueAuthCode(ctx, c, policy, codePurposeEmailChange, email)
}

// VerifyEmailChangeCode verifies the email change code.
func VerifyEmailChangeCode(ctx context.Context, c cache.Cache, email, code string) error {
	return verifyAuthCode(ctx, c, codePurposeEmailChange, email, code)
}

func issueAuthCode(ctx context.Context, c cache.Cache, policy codePolicy, purpose, email string) (string, error) {
	if c == nil {
		return "", errors.New("auth: cache is required")
	}

	emailKey := normalizeEmailInput(email)
	if emailKey == "" {
		return "", repository.NewInvalidArgument("email is required")
	}

	policy = normalizeCodePolicy(policy)
	if policy.SendCooldown > 0 {
		if err := c.Get(ctx, cooldownKey(purpose, emailKey), nil); err == nil {
			return "", repository.ErrTooManyRequests
		} else if err != cache.ErrNotFound {
			return "", err
		}
	}

	if policy.SendLimitPerHour > 0 {
		count, err := incrementCounter(ctx, c, countKey(purpose, emailKey), time.Hour)
		if err != nil {
			return "", err
		}
		if count > policy.SendLimitPerHour {
			return "", repository.ErrTooManyRequests
		}
	}

	code, err := generateNumericCode(policy.CodeLength)
	if err != nil {
		return "", err
	}

	record := codeRecord{
		Hash:   hashCode(code),
		SentAt: time.Now().UTC().Unix(),
	}
	if err := c.Set(ctx, codeKey(purpose, emailKey), record, policy.CodeTTL); err != nil {
		return "", err
	}
	if policy.SendCooldown > 0 {
		_ = c.Set(ctx, cooldownKey(purpose, emailKey), map[string]any{"sent_at": record.SentAt}, policy.SendCooldown)
	}

	return code, nil
}

func verifyAuthCode(ctx context.Context, c cache.Cache, purpose, email, code string) error {
	if c == nil {
		return errors.New("auth: cache is required")
	}
	emailKey := normalizeEmailInput(email)
	if emailKey == "" || strings.TrimSpace(code) == "" {
		return repository.NewInvalidArgument("email and code are required")
	}

	var record codeRecord
	if err := c.Get(ctx, codeKey(purpose, emailKey), &record); err != nil {
		if err == cache.ErrNotFound {
			return repository.NewInvalidArgument("invalid or expired code")
		}
		return err
	}

	if hashCode(code) != record.Hash {
		return repository.NewInvalidArgument("invalid or expired code")
	}

	_ = c.Del(ctx, codeKey(purpose, emailKey))
	return nil
}

func normalizeCodePolicy(policy codePolicy) codePolicy {
	if policy.CodeLength <= 0 {
		policy.CodeLength = 6
	}
	if policy.CodeTTL <= 0 {
		policy.CodeTTL = 15 * time.Minute
	}
	if policy.SendCooldown < 0 {
		policy.SendCooldown = 0
	}
	if policy.SendLimitPerHour < 0 {
		policy.SendLimitPerHour = 0
	}
	return policy
}

func generateNumericCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	var builder strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate code: %w", err)
		}
		builder.WriteByte(byte('0' + n.Int64()))
	}
	return builder.String(), nil
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(code)))
	return hex.EncodeToString(sum[:])
}

func codeKey(purpose, email string) string {
	return fmt.Sprintf("auth:code:%s:%s", strings.ToLower(strings.TrimSpace(purpose)), normalizeEmailInput(email))
}

func cooldownKey(purpose, email string) string {
	return fmt.Sprintf("auth:code:cooldown:%s:%s", strings.ToLower(strings.TrimSpace(purpose)), normalizeEmailInput(email))
}

func countKey(purpose, email string) string {
	return fmt.Sprintf("auth:code:count:%s:%s", strings.ToLower(strings.TrimSpace(purpose)), normalizeEmailInput(email))
}

func incrementCounter(ctx context.Context, c cache.Cache, key string, ttl time.Duration) (int, error) {
	var counter counterRecord
	err := c.Get(ctx, key, &counter)
	if err != nil && err != cache.ErrNotFound {
		return 0, err
	}

	counter.Count++
	if err := c.Set(ctx, key, counter, ttl); err != nil {
		return 0, err
	}
	return counter.Count, nil
}
