// Package apierr defines structured domain error types for the authn microservice.
// Using typed errors instead of string matching allows callers (gRPC handlers,
// tests) to map errors to the correct gRPC status codes without fragile string
// comparisons.
package apierr

import "fmt"

// Code is a domain error classification.
type Code string

const (
	CodeNotFound           Code = "NOT_FOUND"
	CodeAlreadyExists      Code = "ALREADY_EXISTS"
	CodeInvalidCredentials Code = "INVALID_CREDENTIALS"
	CodeInvalidArgument    Code = "INVALID_ARGUMENT"
	CodeExpired            Code = "EXPIRED"
	CodeRateLimited        Code = "RATE_LIMITED"
	CodePermissionDenied   Code = "PERMISSION_DENIED"
	CodeInternal           Code = "INTERNAL"
	CodeUnauthenticated    Code = "UNAUTHENTICATED"
)

// DomainError is a structured service-layer error.
type DomainError struct {
	Code    Code
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error { return e.Cause }

// Constructors

func NotFound(msg string, cause error) *DomainError {
	return &DomainError{Code: CodeNotFound, Message: msg, Cause: cause}
}

func AlreadyExists(msg string, cause error) *DomainError {
	return &DomainError{Code: CodeAlreadyExists, Message: msg, Cause: cause}
}

func InvalidCredentials(msg string) *DomainError {
	return &DomainError{Code: CodeInvalidCredentials, Message: msg}
}

func InvalidArgument(msg string) *DomainError {
	return &DomainError{Code: CodeInvalidArgument, Message: msg}
}

func Expired(msg string) *DomainError {
	return &DomainError{Code: CodeExpired, Message: msg}
}

func RateLimited(msg string) *DomainError {
	return &DomainError{Code: CodeRateLimited, Message: msg}
}

func PermissionDenied(msg string) *DomainError {
	return &DomainError{Code: CodePermissionDenied, Message: msg}
}

func Internal(msg string, cause error) *DomainError {
	return &DomainError{Code: CodeInternal, Message: msg, Cause: cause}
}

func Unauthenticated(msg string) *DomainError {
	return &DomainError{Code: CodeUnauthenticated, Message: msg}
}

// Is allows errors.Is to match by Code.
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Sentinel errors for errors.Is comparisons.
var (
	ErrNotFound           = &DomainError{Code: CodeNotFound}
	ErrAlreadyExists      = &DomainError{Code: CodeAlreadyExists}
	ErrInvalidCredentials = &DomainError{Code: CodeInvalidCredentials}
	ErrInvalidArgument    = &DomainError{Code: CodeInvalidArgument}
	ErrExpired            = &DomainError{Code: CodeExpired}
	ErrRateLimited        = &DomainError{Code: CodeRateLimited}
	ErrPermissionDenied   = &DomainError{Code: CodePermissionDenied}
	ErrInternal           = &DomainError{Code: CodeInternal}
	ErrUnauthenticated    = &DomainError{Code: CodeUnauthenticated}
)
