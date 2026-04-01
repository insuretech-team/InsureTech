package models


// AddOrgMemberRequest represents a add_org_member_request
type AddOrgMemberRequest struct {
	OrganisationId string `json:"organisation_id"`
	Role *OrgMemberRole `json:"role,omitempty"`
	UserId string `json:"user_id"`
}
