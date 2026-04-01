package models


// RemoveOrgMemberRequest represents a remove_org_member_request
type RemoveOrgMemberRequest struct {
	MemberId string `json:"member_id"`
	OrganisationId string `json:"organisation_id"`
}
