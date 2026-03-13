package models


// OrgMembersListingResponse represents a org_members_listing_response
type OrgMembersListingResponse struct {
	Members []*OrgMember `json:"members,omitempty"`
	Error *Error `json:"error,omitempty"`
}
