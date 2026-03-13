package client

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// buildPath replaces path parameters with actual values
// Example: "/v1/policies/{policy_id}" with {"policy_id": "123"} -> "/v1/policies/123"
func buildPath(template string, params map[string]string) string {
	path := template
	for key, value := range params {
		placeholder := "{" + key + "}"
		path = strings.ReplaceAll(path, placeholder, url.PathEscape(value))
	}
	return path
}

// buildPathWithQuery builds a URL path with query parameters
func buildPathWithQuery(basePath string, params interface{}) string {
	if params == nil {
		return basePath
	}

	query := buildQueryString(params)
	if query == "" {
		return basePath
	}

	return basePath + "?" + query
}

// buildQueryString converts a struct to URL query parameters
func buildQueryString(params interface{}) string {
	if params == nil {
		return ""
	}

	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	t := v.Type()
	var parts []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get JSON tag
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON tag (handle omitempty)
		tagParts := strings.Split(jsonTag, ",")
		paramName := tagParts[0]
		omitEmpty := len(tagParts) > 1 && tagParts[1] == "omitempty"

		// Skip if field is zero value and omitempty is set
		if omitEmpty && isZeroValue(field) {
			continue
		}

		// Convert field value to string
		paramValue := formatQueryValue(field)
		if paramValue != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(paramName), url.QueryEscape(paramValue)))
		}
	}

	return strings.Join(parts, "&")
}

// isZeroValue checks if a reflect.Value is the zero value for its type
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}

// formatQueryValue converts a reflect.Value to a query parameter string
func formatQueryValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return formatQueryValue(v.Elem())
	case reflect.Slice, reflect.Array:
		// Handle slices as comma-separated values
		var values []string
		for i := 0; i < v.Len(); i++ {
			val := formatQueryValue(v.Index(i))
			if val != "" {
				values = append(values, val)
			}
		}
		return strings.Join(values, ",")
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// PaginationHelper provides utilities for paginated requests
type PaginationHelper struct {
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
	HasMore    bool
}

// NewPaginationHelper creates a new pagination helper
func NewPaginationHelper(page, pageSize int) *PaginationHelper {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return &PaginationHelper{
		Page:     page,
		PageSize: pageSize,
	}
}

// NextPage returns the next page number
func (p *PaginationHelper) NextPage() int {
	if p.HasMore {
		return p.Page + 1
	}
	return p.Page
}

// PrevPage returns the previous page number
func (p *PaginationHelper) PrevPage() int {
	if p.Page > 1 {
		return p.Page - 1
	}
	return 1
}

// UpdateFromResponse updates pagination info from a response
func (p *PaginationHelper) UpdateFromResponse(totalItems, totalPages int) {
	p.TotalItems = totalItems
	p.TotalPages = totalPages
	p.HasMore = p.Page < p.TotalPages
}

// IsFirstPage returns true if current page is the first page
func (p *PaginationHelper) IsFirstPage() bool {
	return p.Page == 1
}

// IsLastPage returns true if current page is the last page
func (p *PaginationHelper) IsLastPage() bool {
	return !p.HasMore
}

// Offset returns the offset for database queries
func (p *PaginationHelper) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the limit for database queries
func (p *PaginationHelper) Limit() int {
	return p.PageSize
}
