package domain

import "errors"

// Business errors; the handler layer maps each one to an HTTP status
var (
	ErrNotFound           = errors.New("not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTodoDone           = errors.New("completed tasks can't be edited; mark it as not done first")
)
