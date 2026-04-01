package models

import (
	"time"
)

// OrgMemberRoleChangedEvent represents a org_member_role_changed_event
type OrgMemberRoleChangedEvent struct {
	ChangedBy string `json:"changed_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MemberId string `json:"member_id,omitempty"`
	NewRole *OrgMemberRole `json:"new_role,omitempty"`
	OldRole *OrgMemberRole `json:"old_role,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
