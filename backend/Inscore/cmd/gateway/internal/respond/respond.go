// Package respond provides the unified ApiResponse<T> envelope for every HTTP
// response produced by the InsureTech API Gateway — success AND error alike.
//
// Contract (matches api/openapi.yaml ApiResponse schema):
//
//	{
//	  "success": true|false,
//	  "data":    <T> | null,
//	  "error":   null | { code, message, error_id, http_status_code, retryable, field_violations },
//	  "meta":    { request_id, timestamp, pagination? }
//	}
//
// All handlers MUST call respond.JSON / respond.Error — never write raw JSON
// or call http.Error() directly so that every response is shape-identical.
package respond

import (
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ──────────────────────────────────────────────────────────────────────────────
// Core envelope types
// ──────────────────────────────────────────────────────────────────────────────

// ApiResponse is the top-level envelope returned for every endpoint.
// It is generic in Go 1.21+ — callers supply a concrete data type T.
// When the handler writes raw JSON (proto marshalled bytes) use RawJSON.
type ApiResponse[T any] struct {
	Success bool         `json:"success"`
	Data    T            `json:"data"`
	Error   *ApiError    `json:"error"`
	Meta    ResponseMeta `json:"meta"`
}

// ApiError holds structured error information for failed responses.
type ApiError struct {
	Code            string           `json:"code"`
	Message         string           `json:"message"`
	ErrorID         string           `json:"error_id"`
	HTTPStatusCode  int              `json:"http_status_code"`
	Retryable       bool             `json:"retryable"`
	FieldViolations []FieldViolation `json:"field_violations"`
}

// FieldViolation describes a single field-level validation failure.
type FieldViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

// ResponseMeta carries per-request metadata attached to every response.
type ResponseMeta struct {
	RequestID  string          `json:"request_id"`
	Timestamp  string          `json:"timestamp"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// PaginationMeta is included on list responses.
type PaginationMeta struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalCount int  `json:"total_count"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// rawEnvelope is an internal type used when data is already marshalled JSON
// (e.g. from protojson) so we avoid double-marshalling.
type rawEnvelope struct {
	Success bool             `json:"success"`
	Data    json.RawMessage  `json:"data"`
	Error   *ApiError        `json:"error"`
	Meta    ResponseMeta     `json:"meta"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Meta helpers
// ──────────────────────────────────────────────────────────────────────────────

func buildMeta(r *http.Request) ResponseMeta {
	reqID := ""
	if r != nil {
		reqID = r.Header.Get("X-Request-ID")
	}
	return ResponseMeta{
		RequestID: reqID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Success writers
// ──────────────────────────────────────────────────────────────────────────────

// JSON writes a 200 OK response with data as the envelope payload.
// Use this when you have a typed Go struct to return.
func JSON[T any](w http.ResponseWriter, r *http.Request, data T) {
	writeJSON(w, http.StatusOK, ApiResponse[T]{
		Success: true,
		Data:    data,
		Error:   nil,
		Meta:    buildMeta(r),
	})
}

// Created writes a 201 Created response with data as the envelope payload.
func Created[T any](w http.ResponseWriter, r *http.Request, data T) {
	writeJSON(w, http.StatusCreated, ApiResponse[T]{
		Success: true,
		Data:    data,
		Error:   nil,
		Meta:    buildMeta(r),
	})
}

// NoContent writes a 200 OK with data: null for operations that return no body
// (e.g. DELETE, logout, revoke). We use 200 not 204 so the envelope is always
// present — 204 drops the body and breaks the unified contract.
func NoContent(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, ApiResponse[any]{
		Success: true,
		Data:    nil,
		Error:   nil,
		Meta:    buildMeta(r),
	})
}

// RawProtoJSON writes a success response where data is already marshalled JSON
// bytes (e.g. from protojson.Marshal). Wraps them in the envelope without
// double-marshalling.
func RawProtoJSON(w http.ResponseWriter, r *http.Request, dataJSON []byte, statusCode int) {
	if len(dataJSON) == 0 || string(dataJSON) == "null" || string(dataJSON) == "{}" {
		dataJSON = []byte("null")
	}
	env := rawEnvelope{
		Success: true,
		Data:    json.RawMessage(dataJSON),
		Error:   nil,
		Meta:    buildMeta(r),
	}
	writeRaw(w, statusCode, env)
}

// ──────────────────────────────────────────────────────────────────────────────
// Error writers
// ──────────────────────────────────────────────────────────────────────────────

// Error writes a structured error response using the unified ApiResponse envelope.
// data is always null for error responses.
func Error(w http.ResponseWriter, r *http.Request, httpStatus int, code, message string) {
	ErrorWithViolations(w, r, httpStatus, code, message, false, nil)
}

// ErrorRetryable writes a structured error response with retryable=true.
func ErrorRetryable(w http.ResponseWriter, r *http.Request, httpStatus int, code, message string) {
	ErrorWithViolations(w, r, httpStatus, code, message, true, nil)
}

// ErrorWithViolations writes a structured 422 / 400 error with field-level details.
func ErrorWithViolations(w http.ResponseWriter, r *http.Request, httpStatus int, code, message string, retryable bool, violations []FieldViolation) {
	reqID := ""
	if r != nil {
		reqID = r.Header.Get("X-Request-ID")
	}
	if violations == nil {
		violations = []FieldViolation{}
	}
	env := ApiResponse[any]{
		Success: false,
		Data:    nil,
		Error: &ApiError{
			Code:            code,
			Message:         message,
			ErrorID:         reqID,
			HTTPStatusCode:  httpStatus,
			Retryable:       retryable,
			FieldViolations: violations,
		},
		Meta: buildMeta(r),
	}
	writeJSON(w, httpStatus, env)
}

// GRPCError converts a gRPC status error into the unified ApiResponse envelope.
// It maps gRPC codes to HTTP status codes and produces structured error output.
func GRPCError(w http.ResponseWriter, r *http.Request, err error) {
	st, _ := status.FromError(err)
	httpStatus := GRPCCodeToHTTP(st.Code())
	errCode := grpcCodeToErrorCode(st.Code())
	retryable := isRetryable(st.Code())
	ErrorWithViolations(w, r, httpStatus, errCode, st.Message(), retryable, nil)
}

// ──────────────────────────────────────────────────────────────────────────────
// gRPC → HTTP mapping
// ──────────────────────────────────────────────────────────────────────────────

// GRPCCodeToHTTP maps a gRPC status code to its canonical HTTP status code.
func GRPCCodeToHTTP(c codes.Code) int {
	switch c {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499 // client closed request
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusUnprocessableEntity
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusUnprocessableEntity
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func grpcCodeToErrorCode(c codes.Code) string {
	switch c {
	case codes.InvalidArgument:
		return "INVALID_ARGUMENT"
	case codes.NotFound:
		return "NOT_FOUND"
	case codes.AlreadyExists:
		return "ALREADY_EXISTS"
	case codes.PermissionDenied:
		return "PERMISSION_DENIED"
	case codes.Unauthenticated:
		return "UNAUTHENTICATED"
	case codes.ResourceExhausted:
		return "RATE_LIMITED"
	case codes.FailedPrecondition:
		return "FAILED_PRECONDITION"
	case codes.Unimplemented:
		return "NOT_IMPLEMENTED"
	case codes.Unavailable:
		return "SERVICE_UNAVAILABLE"
	case codes.DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case codes.Internal:
		return "INTERNAL_ERROR"
	case codes.Aborted:
		return "CONFLICT"
	default:
		return "INTERNAL_ERROR"
	}
}

func isRetryable(c codes.Code) bool {
	switch c {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
		return true
	default:
		return false
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Internal writers
// ──────────────────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		// Last-resort fallback — must never happen in practice
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"success":false,"data":null,"error":{"code":"MARSHAL_ERROR","message":"failed to encode response","error_id":"","http_status_code":500,"retryable":false,"field_violations":[]},"meta":{"request_id":"","timestamp":""}}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}

func writeRaw(w http.ResponseWriter, statusCode int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse[any]{
			Success: false,
			Data:    nil,
			Error: &ApiError{
				Code:            "MARSHAL_ERROR",
				Message:         "failed to encode response",
				HTTPStatusCode:  http.StatusInternalServerError,
				Retryable:       false,
				FieldViolations: []FieldViolation{},
			},
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}
