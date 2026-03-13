package models


// AddOrgMemberRequest represents a add_org_member_request
type AddOrgMemberRequest struct {
	Role *OrgMemberRole `json:"role,omitempty"`
	OrganisationId string `json:"organisation_id"`
	UserId string `json:"user_id"`
}
