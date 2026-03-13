package models

import (
	"time"
)

// PartnerOnboardedEvent represents a partner_onboarded_event
type PartnerOnboardedEvent struct {
	PartnerType string `json:"partner_type,omitempty"`
	FocalPersonId string `json:"focal_person_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
}
