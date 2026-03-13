package models

import (
	"time"
)

// OrgMember represents a org_member
type OrgMember struct {
	JoinedAt time.Time `json:"joined_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Role *OrgMemberRole `json:"role,omitempty"`
	Status *OrgMemberStatus `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	MemberId string `json:"member_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
