package models

import (
	"time"
)

// OrgMemberRoleChangedEvent represents a org_member_role_changed_event
type OrgMemberRoleChangedEvent struct {
	UserId string `json:"user_id,omitempty"`
	OldRole *OrgMemberRole `json:"old_role,omitempty"`
	NewRole *OrgMemberRole `json:"new_role,omitempty"`
	ChangedBy string `json:"changed_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MemberId string `json:"member_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
}
