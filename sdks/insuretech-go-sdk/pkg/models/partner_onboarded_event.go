package models

import (
	"time"
)

// PartnerOnboardedEvent represents a partner_onboarded_event
type PartnerOnboardedEvent struct {
	EventId string `json:"event_id,omitempty"`
	FocalPersonId string `json:"focal_person_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	PartnerType string `json:"partner_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
