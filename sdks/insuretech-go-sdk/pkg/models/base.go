package models

import (
	"time"
)

// ModelModifier provides methods to modify and extend generated models
type ModelModifier interface {
	Validate() error
	Transform() error
}

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// IsDeleted returns true if the model has been soft-deleted
func (m *BaseModel) IsDeleted() bool {
	return m.DeletedAt != nil
}

// Metadata represents additional metadata for models
type Metadata map[string]interface{}

// Get retrieves a metadata value by key
func (m Metadata) Get(key string) (interface{}, bool) {
	val, ok := m[key]
	return val, ok
}

// GetString retrieves a string metadata value
func (m Metadata) GetString(key string) (string, bool) {
	val, ok := m[key]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// GetInt retrieves an integer metadata value
func (m Metadata) GetInt(key string) (int, bool) {
	val, ok := m[key]
	if !ok {
		return 0, false
	}
	
	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// Set sets a metadata value
func (m Metadata) Set(key string, value interface{}) {
	m[key] = value
}

// Note: Common types like Money, Address, ContactInfo are generated from OpenAPI spec
// They will be in their own files (money.go, address.go, contact_info.go)

// DateRange represents a date range
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// IsActive checks if the date range is currently active
func (dr *DateRange) IsActive() bool {
	now := time.Now()
	return now.After(dr.StartDate) && now.Before(dr.EndDate)
}

// Duration returns the duration of the date range
func (dr *DateRange) Duration() time.Duration {
	return dr.EndDate.Sub(dr.StartDate)
}
