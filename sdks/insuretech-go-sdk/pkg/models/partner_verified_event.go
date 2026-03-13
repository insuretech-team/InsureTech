package models

import (
	"time"
)

// PartnerVerifiedEvent represents a partner_verified_event
type PartnerVerifiedEvent struct {
	PartnerId string `json:"partner_id,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
