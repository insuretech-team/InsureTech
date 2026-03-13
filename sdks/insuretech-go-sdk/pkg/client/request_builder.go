package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestBuilder provides a fluent interface for building HTTP requests
type RequestBuilder struct {
	client      *Client
	method      string
	path        string
	pathParams  map[string]string
	queryParams url.Values
	body        interface{}
	headers     map[string]string
	ctx         context.Context
}

// NewRequestBuilder creates a new request builder
func (c *Client) NewRequestBuilder(method, path string) *RequestBuilder {
	return &RequestBuilder{
		client:      c,
		method:      method,
		path:        path,
		pathParams:  make(map[string]string),
		queryParams: make(url.Values),
		headers:     make(map[string]string),
		ctx:         context.Background(),
	}
}

// WithContext sets the context for the request
func (rb *RequestBuilder) WithContext(ctx context.Context) *RequestBuilder {
	rb.ctx = ctx
	return rb
}

// WithPathParam adds a path parameter
func (rb *RequestBuilder) WithPathParam(key, value string) *RequestBuilder {
	rb.pathParams[key] = value
	return rb
}

// WithPathParams adds multiple path parameters
func (rb *RequestBuilder) WithPathParams(params map[string]string) *RequestBuilder {
	for k, v := range params {
		rb.pathParams[k] = v
	}
	return rb
}

// WithQueryParam adds a query parameter
func (rb *RequestBuilder) WithQueryParam(key, value string) *RequestBuilder {
	rb.queryParams.Add(key, value)
	return rb
}

// WithQueryParams adds query parameters from a struct
func (rb *RequestBuilder) WithQueryParams(params interface{}) *RequestBuilder {
	if params != nil {
		queryString := buildQueryString(params)
		if queryString != "" {
			values, _ := url.ParseQuery(queryString)
			for k, v := range values {
				for _, val := range v {
					rb.queryParams.Add(k, val)
				}
			}
		}
	}
	return rb
}

// WithBody sets the request body
func (rb *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
	rb.body = body
	return rb
}

// WithHeader adds a header
func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

// WithHeaders adds multiple headers
func (rb *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	for k, v := range headers {
		rb.headers[k] = v
	}
	return rb
}

// Do executes the request and unmarshals the response into result
func (rb *RequestBuilder) Do(result interface{}) error {
	// Build the path with path parameters
	path := rb.path
	for key, value := range rb.pathParams {
		placeholder := "{" + key + "}"
		path = strings.ReplaceAll(path, placeholder, url.PathEscape(value))
	}

	// Add query parameters
	if len(rb.queryParams) > 0 {
		path = path + "?" + rb.queryParams.Encode()
	}

	// Build URL
	u, err := url.Parse(rb.client.baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Prepare request body
	var bodyReader io.Reader
	if rb.body != nil {
		bodyBytes, err := json.Marshal(rb.body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(rb.ctx, rb.method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "InsureTech-Go-SDK/1.0.0")

	// Set authentication
	if rb.client.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+rb.client.apiKey)
	}

	// Set custom headers
	for key, value := range rb.headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := rb.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	// Parse response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// DoRaw executes the request and returns the raw response
func (rb *RequestBuilder) DoRaw() (*http.Response, error) {
	// Build the path with path parameters
	path := rb.path
	for key, value := range rb.pathParams {
		placeholder := "{" + key + "}"
		path = strings.ReplaceAll(path, placeholder, url.PathEscape(value))
	}

	// Add query parameters
	if len(rb.queryParams) > 0 {
		path = path + "?" + rb.queryParams.Encode()
	}

	// Build URL
	u, err := url.Parse(rb.client.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Prepare request body
	var bodyReader io.Reader
	if rb.body != nil {
		bodyBytes, err := json.Marshal(rb.body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(rb.ctx, rb.method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "InsureTech-Go-SDK/1.0.0")

	// Set authentication
	if rb.client.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+rb.client.apiKey)
	}

	// Set custom headers
	for key, value := range rb.headers {
		req.Header.Set(key, value)
	}

	// Execute request
	return rb.client.httpClient.Do(req)
}

// DoBytes executes the request and returns the response body as bytes
func (rb *RequestBuilder) DoBytes() ([]byte, error) {
	resp, err := rb.DoRaw()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp.StatusCode, body)
	}

	return body, nil
}

// DoString executes the request and returns the response body as a string
func (rb *RequestBuilder) DoString() (string, error) {
	body, err := rb.DoBytes()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// parseAPIError parses an API error from response
func parseAPIError(statusCode int, body []byte) error {
	// Try to parse as JSON error
	var errResp struct {
		Error struct {
			Code      string                 `json:"code"`
			Message   string                 `json:"message"`
			Details   map[string]interface{} `json:"details"`
			RequestID string                 `json:"request_id"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return &APIError{
			StatusCode: statusCode,
			Code:       errResp.Error.Code,
			Message:    errResp.Error.Message,
			Details:    errResp.Error.Details,
			RequestID:  errResp.Error.RequestID,
		}
	}

	// Fallback to simple error
	return &APIError{
		StatusCode: statusCode,
		Message:    string(body),
	}
}

// APIError represents an API error
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
