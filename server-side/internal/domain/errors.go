package domain

import "errors"

var (
	// ErrNotFound indicates that the requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrConflict indicates a version conflict (optimistic locking)
	ErrConflict = errors.New("resource conflict: version mismatch")

	// ErrInvalidInput indicates invalid input data
	ErrInvalidInput = errors.New("invalid input")

	// ErrAlreadyDeleted indicates the resource is already soft deleted
	ErrAlreadyDeleted = errors.New("resource already deleted")

	// ErrDuplicatePhoneNumber indicates a phone number already exists
	ErrDuplicatePhoneNumber = errors.New("phone number already exists")
)
