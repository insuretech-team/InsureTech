package client

// ---------------------------------------------------------------------------
// Rule 01: Standard API response envelope — used by ALL endpoints
// ---------------------------------------------------------------------------

// ApiResponse is the standard envelope wrapping every API response.
// success=true  → Data has payload, Err is nil
// success=false → Data is nil, Err has details
type ApiResponse[T any] struct {
	Success bool         `json:"success"`
	Data    *T           `json:"data"`
	Error   *ApiError    `json:"error"`
	Meta    *ResponseMeta `json:"meta"`
}

// ApiError is returned in the envelope when success=false.
// Rule 03: Errors never live inside success response schemas.
type ApiError struct {
	Code               string           `json:"code"`
	Message            string           `json:"message"`
	FieldViolations    []FieldViolation `json:"field_violations,omitempty"`
	ErrorID            string           `json:"error_id,omitempty"`
	Retryable          bool             `json:"retryable,omitempty"`
	RetryAfterSeconds  int              `json:"retry_after_seconds,omitempty"`
	HTTPStatusCode     int              `json:"http_status_code,omitempty"`
	DocumentationURL   string           `json:"documentation_url,omitempty"`
}

// Error implements the error interface so ApiError can be returned as error.
func (e *ApiError) Error() string {
	if e == nil {
		return ""
	}
	return e.Code + ": " + e.Message
}

// FieldViolation describes a field-level validation error (Rule 03 — 422 responses).
type FieldViolation struct {
	Field         string `json:"field"`
	Message       string `json:"message"`
	Code          string `json:"code,omitempty"`
	RejectedValue string `json:"rejected_value,omitempty"`
}

// ResponseMeta is attached to every API response for tracing and pagination.
type ResponseMeta struct {
	RequestID  string          `json:"request_id,omitempty"`
	Timestamp  string          `json:"timestamp,omitempty"`
	APIVersion string          `json:"api_version,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// PaginationMeta is the single standard pagination schema (Rule 05).
// Replaces the deprecated PageResponse and PaginationResponse schemas.
type PaginationMeta struct {
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	TotalPages      int    `json:"total_pages"`
	TotalItems      int64  `json:"total_items"`
	HasNext         bool   `json:"has_next"`
	HasPrevious     bool   `json:"has_previous"`
	NextPageToken   string `json:"next_page_token,omitempty"`
}

// ListResult is a convenience type for list endpoint responses (Rule 05).
type ListResult[T any] struct {
	Items      []T
	Pagination *PaginationMeta
}

// ---------------------------------------------------------------------------
// InsureTechError — typed error returned by DoRequestEnvelope
// ---------------------------------------------------------------------------

// InsureTechError is returned when the API responds with success=false.
// Implements the error interface and carries structured ApiError data.
type InsureTechError struct {
	APIError   *ApiError
	HTTPStatus int
}

func (e *InsureTechError) Error() string {
	if e.APIError == nil {
		return "InsureTech API error"
	}
	return e.APIError.Error()
}

// IsValidationError returns true if this is a 422 validation error.
func (e *InsureTechError) IsValidationError() bool {
	return e.HTTPStatus == 422
}

// IsAuthError returns true if this is a 401 authentication error.
func (e *InsureTechError) IsAuthError() bool {
	return e.HTTPStatus == 401
}

// IsNotFound returns true if this is a 404 not found error.
func (e *InsureTechError) IsNotFound() bool {
	return e.HTTPStatus == 404
}

// IsRateLimited returns true if this is a 429 rate limit error.
func (e *InsureTechError) IsRateLimited() bool {
	return e.HTTPStatus == 429
}

// FieldErrors returns a map of field → messages from field_violations (Rule 03).
func (e *InsureTechError) FieldErrors() map[string]string {
	result := make(map[string]string)
	if e.APIError == nil {
		return result
	}
	for _, v := range e.APIError.FieldViolations {
		result[v.Field] = v.Message
	}
	return result
}
