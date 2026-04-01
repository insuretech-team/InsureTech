package models


// OrgAdminAssignmentRequest represents a org_admin_assignment_request
type OrgAdminAssignmentRequest struct {
	MemberId string `json:"member_id"`
	OrganisationId string `json:"organisation_id"`
}
