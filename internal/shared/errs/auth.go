package errs

import "errors"

// Sentinel errors - wrap these in repo/service to control what TUI sees
var (
	// User-facing errors (TUI sees err.Error())
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("already exists")

	// Auth errors
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("email already exists")
)
