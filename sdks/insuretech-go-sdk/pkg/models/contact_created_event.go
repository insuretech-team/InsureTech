package models

import (
	"time"
)

// ContactCreatedEvent represents a contact_created_event
type ContactCreatedEvent struct {
	ContactId string `json:"contact_id,omitempty"`
	ContactType *ContactType `json:"contact_type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}
