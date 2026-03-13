package models

import (
	"time"
)

// Organisation represents a organisation
type Organisation struct {
	Code string `json:"code,omitempty"`
	Industry string `json:"industry,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Name string `json:"name,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Address string `json:"address,omitempty"`
	Status *OrganisationStatus `json:"status,omitempty"`
	TotalEmployees int `json:"total_employees,omitempty"`
}
