package models

import (
	"time"
)

// OrgMemberAddedEvent represents a org_member_added_event
type OrgMemberAddedEvent struct {
	AddedBy string `json:"added_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MemberId string `json:"member_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Role *OrgMemberRole `json:"role,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
