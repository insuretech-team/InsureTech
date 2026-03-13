package models


// RemoveOrgMemberRequest represents a remove_org_member_request
type RemoveOrgMemberRequest struct {
	OrganisationId string `json:"organisation_id"`
	MemberId string `json:"member_id"`
}
