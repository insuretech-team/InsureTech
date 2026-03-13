package models


// OrgAdminAssignmentRequest represents a org_admin_assignment_request
type OrgAdminAssignmentRequest struct {
	OrganisationId string `json:"organisation_id"`
	MemberId string `json:"member_id"`
}
