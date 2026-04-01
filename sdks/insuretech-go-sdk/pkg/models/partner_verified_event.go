package models

import (
	"time"
)

// PartnerVerifiedEvent represents a partner_verified_event
type PartnerVerifiedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
