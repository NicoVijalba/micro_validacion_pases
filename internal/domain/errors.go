// Package domain contains core business entities and domain-level errors.
package domain

import "errors"

var (
	// ErrUnauthorized indicates the current principal is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates the current principal is authenticated but lacks permissions.
	ErrForbidden    = errors.New("forbidden")
	// ErrInvalidInput indicates the request payload violates domain validation rules.
	ErrInvalidInput = errors.New("invalid input")
	// ErrConflict indicates the resource state conflicts with the requested operation.
	ErrConflict     = errors.New("conflict")
)
