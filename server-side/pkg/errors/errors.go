package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a unique error code
type ErrorCode string

const (
	// General errors
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnprocessable   ErrorCode = "UNPROCESSABLE_ENTITY"

	// Authentication errors
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeEmailAlreadyExists ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrCodeInvalidToken       ErrorCode = "INVALID_TOKEN"
	ErrCodeExpiredToken       ErrorCode = "EXPIRED_TOKEN"

	// Resource errors
	ErrCodeUserNotFound     ErrorCode = "USER_NOT_FOUND"
	ErrCodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCodeVersionConflict  ErrorCode = "VERSION_CONFLICT"

	// Business logic errors
	ErrCodeInvalidInput       ErrorCode = "INVALID_INPUT"
	ErrCodeInsufficientFunds  ErrorCode = "INSUFFICIENT_FUNDS"
	ErrCodeOperationNotAllowed ErrorCode = "OPERATION_NOT_ALLOWED"
)

// AppError represents an application error with code and HTTP status
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// New creates a new AppError
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with AppError
func Wrap(err error, code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}

	appErr, ok := err.(*AppError)
	return appErr, ok
}

// GetHTTPStatus extracts HTTP status from error, defaults to 500
func GetHTTPStatus(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// GetErrorCode extracts error code from error
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Code
	}
	return ErrCodeInternal
}

// Predefined errors - General
var (
	ErrInternal = New(
		ErrCodeInternal,
		"An internal error occurred",
		http.StatusInternalServerError,
	)

	ErrBadRequest = New(
		ErrCodeBadRequest,
		"Invalid request",
		http.StatusBadRequest,
	)

	ErrUnauthorized = New(
		ErrCodeUnauthorized,
		"Unauthorized access",
		http.StatusUnauthorized,
	)

	ErrForbidden = New(
		ErrCodeForbidden,
		"Access forbidden",
		http.StatusForbidden,
	)

	ErrNotFound = New(
		ErrCodeNotFound,
		"Resource not found",
		http.StatusNotFound,
	)

	ErrConflict = New(
		ErrCodeConflict,
		"Resource conflict",
		http.StatusConflict,
	)

	ErrValidation = New(
		ErrCodeValidation,
		"Validation failed",
		http.StatusBadRequest,
	)
)

// Predefined errors - Authentication
var (
	ErrInvalidCredentials = New(
		ErrCodeInvalidCredentials,
		"Invalid email or password",
		http.StatusUnauthorized,
	)

	ErrEmailAlreadyExists = New(
		ErrCodeEmailAlreadyExists,
		"Email already registered",
		http.StatusConflict,
	)

	ErrInvalidToken = New(
		ErrCodeInvalidToken,
		"Invalid authentication token",
		http.StatusUnauthorized,
	)

	ErrExpiredToken = New(
		ErrCodeExpiredToken,
		"Authentication token has expired",
		http.StatusUnauthorized,
	)
)

// Predefined errors - Resources
var (
	ErrUserNotFound = New(
		ErrCodeUserNotFound,
		"User not found",
		http.StatusNotFound,
	)

	ErrResourceNotFound = New(
		ErrCodeResourceNotFound,
		"Resource not found",
		http.StatusNotFound,
	)

	ErrVersionConflict = New(
		ErrCodeVersionConflict,
		"Resource version conflict",
		http.StatusConflict,
	)
)

// Predefined errors - Business Logic
var (
	ErrInvalidInput = New(
		ErrCodeInvalidInput,
		"Invalid input provided",
		http.StatusBadRequest,
	)

	ErrOperationNotAllowed = New(
		ErrCodeOperationNotAllowed,
		"Operation not allowed",
		http.StatusForbidden,
	)
)
