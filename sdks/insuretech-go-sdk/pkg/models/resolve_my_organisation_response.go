package models


// ResolveMyOrganisationResponse represents a resolve_my_organisation_response
type ResolveMyOrganisationResponse struct {
	OrganisationId string `json:"organisation_id,omitempty"`
	OrganisationName string `json:"organisation_name,omitempty"`
	Role *OrgMemberRole `json:"role,omitempty"`
}
