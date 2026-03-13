package models


// AddOrgMemberResponse represents a add_org_member_response
type AddOrgMemberResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	Member *OrgMember `json:"member,omitempty"`
}
