package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

// ValidationRule defines a validation rule
type ValidationRule struct {
	Field     string
	Required  bool
	Type      string // "string", "int", "float", "bool", "email", "uuid"
	MinLength int
	MaxLength int
	Min       float64
	Max       float64
	Pattern   string
	Enum      []string
	compiled  *regexp.Regexp
}

// PathValidation defines validation rules for a specific path
type PathValidation struct {
	Method string
	Path   string
	Rules  []ValidationRule
}

// RequestValidator validates incoming requests
type RequestValidator struct {
	validations map[string]map[string]*PathValidation // path -> method -> validation

	// Statistics
	totalRequests   atomic.Int64
	validRequests   atomic.Int64
	invalidRequests atomic.Int64
}

// NewRequestValidator creates a new request validator
func NewRequestValidator() *RequestValidator {
	rv := &RequestValidator{
		validations: make(map[string]map[string]*PathValidation),
	}

	// Register default validations
	rv.registerDefaultValidations()

	return rv
}

// registerDefaultValidations registers common validation rules
func (rv *RequestValidator) registerDefaultValidations() {
	// Order creation validation
	rv.AddValidation(&PathValidation{
		Method: "POST",
		Path:   "/v1/orders",
		Rules: []ValidationRule{
			{Field: "customer_id", Required: true, Type: "uuid"},
			{Field: "items", Required: true, Type: "array"},
			{Field: "total", Required: true, Type: "float", Min: 0.01},
			{Field: "currency", Required: false, Type: "string", MaxLength: 3},
		},
	})

	// Product creation validation
	rv.AddValidation(&PathValidation{
		Method: "POST",
		Path:   "/v1/products",
		Rules: []ValidationRule{
			{Field: "name", Required: true, Type: "string", MinLength: 1, MaxLength: 200},
			{Field: "sku", Required: true, Type: "string", MinLength: 1, MaxLength: 50},
			{Field: "price", Required: true, Type: "float", Min: 0},
			{Field: "category_id", Required: false, Type: "uuid"},
		},
	})

	// Shipping rate request validation
	rv.AddValidation(&PathValidation{
		Method: "POST",
		Path:   "/v1/shipping/rates",
		Rules: []ValidationRule{
			{Field: "origin", Required: true, Type: "object"},
			{Field: "destination", Required: true, Type: "object"},
			{Field: "weight", Required: true, Type: "float", Min: 0.01},
			{Field: "length", Required: false, Type: "float", Min: 0},
			{Field: "width", Required: false, Type: "float", Min: 0},
			{Field: "height", Required: false, Type: "float", Min: 0},
		},
	})
}

// AddValidation adds a validation rule for a path
func (rv *RequestValidator) AddValidation(validation *PathValidation) error {
	// Compile regex patterns
	for i := range validation.Rules {
		if validation.Rules[i].Pattern != "" {
			compiled, err := regexp.Compile(validation.Rules[i].Pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern for field %s: %w", validation.Rules[i].Field, err)
			}
			validation.Rules[i].compiled = compiled
		}
	}

	if rv.validations[validation.Path] == nil {
		rv.validations[validation.Path] = make(map[string]*PathValidation)
	}

	rv.validations[validation.Path][validation.Method] = validation
	return nil
}

// validateValue validates a single value against a rule
func validateValue(value interface{}, rule ValidationRule) error {
	// Check required
	if rule.Required && value == nil {
		return fmt.Errorf("field '%s' is required", rule.Field)
	}

	if value == nil {
		return nil // Optional field not provided
	}

	// Type-specific validation
	switch rule.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field '%s' must be a string", rule.Field)
		}

		if rule.MinLength > 0 && len(str) < rule.MinLength {
			return fmt.Errorf("field '%s' must be at least %d characters", rule.Field, rule.MinLength)
		}

		if rule.MaxLength > 0 && len(str) > rule.MaxLength {
			return fmt.Errorf("field '%s' must be at most %d characters", rule.Field, rule.MaxLength)
		}

		if rule.Pattern != "" && rule.compiled != nil {
			if !rule.compiled.MatchString(str) {
				return fmt.Errorf("field '%s' does not match required pattern", rule.Field)
			}
		}

		if len(rule.Enum) > 0 {
			valid := false
			for _, e := range rule.Enum {
				if str == e {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("field '%s' must be one of: %v", rule.Field, rule.Enum)
			}
		}

	case "int":
		var num int64
		switch v := value.(type) {
		case float64:
			num = int64(v)
		case int:
			num = int64(v)
		case int64:
			num = v
		default:
			return fmt.Errorf("field '%s' must be an integer", rule.Field)
		}

		if rule.Min != 0 && float64(num) < rule.Min {
			return fmt.Errorf("field '%s' must be at least %v", rule.Field, rule.Min)
		}

		if rule.Max != 0 && float64(num) > rule.Max {
			return fmt.Errorf("field '%s' must be at most %v", rule.Field, rule.Max)
		}

	case "float":
		var num float64
		switch v := value.(type) {
		case float64:
			num = v
		case int:
			num = float64(v)
		default:
			return fmt.Errorf("field '%s' must be a number", rule.Field)
		}

		if rule.Min != 0 && num < rule.Min {
			return fmt.Errorf("field '%s' must be at least %v", rule.Field, rule.Min)
		}

		if rule.Max != 0 && num > rule.Max {
			return fmt.Errorf("field '%s' must be at most %v", rule.Field, rule.Max)
		}

	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field '%s' must be a boolean", rule.Field)
		}

	case "email":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field '%s' must be a string", rule.Field)
		}

		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(str) {
			return fmt.Errorf("field '%s' must be a valid email address", rule.Field)
		}

	case "uuid":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field '%s' must be a string", rule.Field)
		}

		uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
		if !uuidRegex.MatchString(strings.ToLower(str)) {
			return fmt.Errorf("field '%s' must be a valid UUID", rule.Field)
		}

	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("field '%s' must be an array", rule.Field)
		}

	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("field '%s' must be an object", rule.Field)
		}
	}

	return nil
}

// validateRequest validates a request body against rules
func (rv *RequestValidator) validateRequest(body map[string]interface{}, validation *PathValidation) []string {
	var errors []string

	for _, rule := range validation.Rules {
		value := body[rule.Field]
		if err := validateValue(value, rule); err != nil {
			errors = append(errors, err.Error())
		}
	}

	return errors
}

// Middleware implements request validation
func (rv *RequestValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv.totalRequests.Add(1)

		// Get validation rules for this path
		pathValidations, hasPath := rv.validations[r.URL.Path]
		if !hasPath {
			// No validation rules for this path
			rv.validRequests.Add(1)
			next.ServeHTTP(w, r)
			return
		}

		validation, hasMethod := pathValidations[r.Method]
		if !hasMethod {
			// No validation rules for this method
			rv.validRequests.Add(1)
			next.ServeHTTP(w, r)
			return
		}

		// Skip validation for multipart/form-data requests (e.g., file uploads)
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			// Pass through to handler without validation
			next.ServeHTTP(w, r)
			return
		}

		// Read request body
		if r.Body == nil {
			rv.invalidRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Request body is required",
			})
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			rv.invalidRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Failed to read request body",
			})
			return
		}

		// Restore body for next handler
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse JSON
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			rv.invalidRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Invalid JSON",
			})
			return
		}

		// Validate
		validationErrors := rv.validateRequest(body, validation)
		if len(validationErrors) > 0 {
			rv.invalidRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":  "Validation failed",
				"errors": validationErrors,
			})
			return
		}

		// Restore body again for next handler
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		rv.validRequests.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Stats returns validation statistics
func (rv *RequestValidator) Stats() map[string]interface{} {
	total := rv.totalRequests.Load()
	valid := rv.validRequests.Load()
	invalid := rv.invalidRequests.Load()

	var validRate float64
	if total > 0 {
		validRate = float64(valid) / float64(total) * 100
	}

	return map[string]interface{}{
		"total_requests":   total,
		"valid_requests":   valid,
		"invalid_requests": invalid,
		"validation_rate":  validRate,
		"registered_paths": len(rv.validations),
	}
}

// ValidateQueryParam validates a query parameter
func ValidateQueryParam(r *http.Request, param string, required bool, validator func(string) error) error {
	value := r.URL.Query().Get(param)

	if value == "" && required {
		return fmt.Errorf("query parameter '%s' is required", param)
	}

	if value != "" && validator != nil {
		return validator(value)
	}

	return nil
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(s string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(s))
}

// IsValidEmail checks if a string is a valid email
func IsValidEmail(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

// IsValidInt checks if a string is a valid integer in range
func IsValidInt(s string, min, max int) bool {
	val, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return val >= min && val <= max
}
