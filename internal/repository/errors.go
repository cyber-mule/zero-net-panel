package repository

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound            = errors.New("repository: resource not found")
	ErrConflict            = errors.New("repository: conflict")
	ErrInvalidArgument     = &InvalidArgumentError{}
	ErrForbidden           = errors.New("repository: forbidden")
	ErrUnauthorized        = errors.New("repository: unauthorized")
	ErrInsufficientBalance = errors.New("repository: insufficient balance")
	ErrInvalidState        = errors.New("repository: invalid state")
	ErrTooManyRequests     = errors.New("repository: too many requests")
	ErrInviteCodeRequired  = errors.New("repository: invite code required")
	ErrInviteCodeInvalid   = errors.New("repository: invite code invalid")
)

// InvalidArgumentError captures a user-facing validation message while remaining
// compatible with errors.Is(err, ErrInvalidArgument).
type InvalidArgumentError struct {
	msg string
}

// Error returns a user-facing validation message.
func (e *InvalidArgumentError) Error() string {
	if e == nil || strings.TrimSpace(e.msg) == "" {
		return "invalid argument"
	}
	return e.msg
}

// Is allows errors.Is to match any InvalidArgumentError.
func (e *InvalidArgumentError) Is(target error) bool {
	_, ok := target.(*InvalidArgumentError)
	return ok
}

// NewInvalidArgument constructs a new invalid-argument error with details.
func NewInvalidArgument(message string) error {
	return &InvalidArgumentError{msg: strings.TrimSpace(message)}
}

// InvalidArgumentf constructs a formatted invalid-argument error.
func InvalidArgumentf(format string, args ...any) error {
	return &InvalidArgumentError{msg: fmt.Sprintf(format, args...)}
}
