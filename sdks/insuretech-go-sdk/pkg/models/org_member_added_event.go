package models

import (
	"time"
)

// OrgMemberAddedEvent represents a org_member_added_event
type OrgMemberAddedEvent struct {
	Role *OrgMemberRole `json:"role,omitempty"`
	AddedBy string `json:"added_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MemberId string `json:"member_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
