package models

import (
	"time"
)

// Organisation represents a organisation
type Organisation struct {
	Address string `json:"address,omitempty"`
	Code string `json:"code,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Industry string `json:"industry,omitempty"`
	Name string `json:"name,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Status *OrganisationStatus `json:"status,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TotalEmployees int `json:"total_employees,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
