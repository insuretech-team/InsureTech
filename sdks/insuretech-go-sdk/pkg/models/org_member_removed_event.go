package models

import (
	"time"
)

// OrgMemberRemovedEvent represents a org_member_removed_event
type OrgMemberRemovedEvent struct {
	UserId string `json:"user_id,omitempty"`
	RemovedBy string `json:"removed_by,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MemberId string `json:"member_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
}
