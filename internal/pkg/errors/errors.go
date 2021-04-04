package errors

import (
	"errors"
)

// General service errors
var (
	// Common errors
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternal       = errors.New("internal service error")
	// 404
	ErrNotFound = errors.New("not found")
)
