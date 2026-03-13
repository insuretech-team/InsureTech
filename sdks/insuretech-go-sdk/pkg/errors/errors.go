package errors

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error returned by the InsureTech API
type APIError struct {
	StatusCode int                    `json:"status_code"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("API error (status: %d, code: %s, request_id: %s): %s",
			e.StatusCode, e.Code, e.RequestID, e.Message)
	}
	return fmt.Sprintf("API error (status: %d, code: %s): %s",
		e.StatusCode, e.Code, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found error
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if the error is a 403 Forbidden error
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsValidationError returns true if the error is a 400 Bad Request error
func (e *APIError) IsValidationError() bool {
	return e.StatusCode == 400
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == 429
}

// IsServerError returns true if the error is a 5xx Server Error
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// parseAPIError parses an API error response
func parseAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{
		StatusCode: statusCode,
	}
	
	// Try to parse JSON error response
	if len(body) > 0 {
		var errResp struct {
			Error struct {
				Code      string                 `json:"code"`
				Message   string                 `json:"message"`
				Details   map[string]interface{} `json:"details"`
				RequestID string                 `json:"request_id"`
			} `json:"error"`
		}
		
		if err := json.Unmarshal(body, &errResp); err == nil {
			apiErr.Code = errResp.Error.Code
			apiErr.Message = errResp.Error.Message
			apiErr.Details = errResp.Error.Details
			apiErr.RequestID = errResp.Error.RequestID
		} else {
			// Fallback to plain text message
			apiErr.Message = string(body)
		}
	}
	
	// Set default message based on status code if not provided
	if apiErr.Message == "" {
		apiErr.Message = getDefaultErrorMessage(statusCode)
	}
	
	return apiErr
}

// getDefaultErrorMessage returns a default error message for a status code
func getDefaultErrorMessage(statusCode int) string {
	switch statusCode {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 409:
		return "Conflict"
	case 422:
		return "Unprocessable Entity"
	case 429:
		return "Too Many Requests"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	case 504:
		return "Gateway Timeout"
	default:
		return fmt.Sprintf("HTTP Error %d", statusCode)
	}
}

// Common error types
var (
	// ErrInvalidRequest indicates an invalid request
	ErrInvalidRequest = &APIError{StatusCode: 400, Code: "invalid_request", Message: "Invalid request"}
	
	// ErrUnauthorized indicates authentication is required or failed
	ErrUnauthorized = &APIError{StatusCode: 401, Code: "unauthorized", Message: "Unauthorized"}
	
	// ErrForbidden indicates the user doesn't have permission
	ErrForbidden = &APIError{StatusCode: 403, Code: "forbidden", Message: "Forbidden"}
	
	// ErrNotFound indicates the resource was not found
	ErrNotFound = &APIError{StatusCode: 404, Code: "not_found", Message: "Not Found"}
	
	// ErrConflict indicates a conflict with existing data
	ErrConflict = &APIError{StatusCode: 409, Code: "conflict", Message: "Conflict"}
	
	// ErrRateLimited indicates the rate limit has been exceeded
	ErrRateLimited = &APIError{StatusCode: 429, Code: "rate_limited", Message: "Too Many Requests"}
	
	// ErrServerError indicates a server error
	ErrServerError = &APIError{StatusCode: 500, Code: "server_error", Message: "Internal Server Error"}
)

// ValidationError represents a validation error with field-level details
type ValidationError struct {
	APIError
	Fields map[string][]string `json:"fields"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s (fields: %v)", e.Message, e.Fields)
}

// AddFieldError adds a field-level validation error
func (e *ValidationError) AddFieldError(field, message string) {
	if e.Fields == nil {
		e.Fields = make(map[string][]string)
	}
	e.Fields[field] = append(e.Fields[field], message)
}

// HasFieldErrors returns true if there are field-level validation errors
func (e *ValidationError) HasFieldErrors() bool {
	return len(e.Fields) > 0
}
