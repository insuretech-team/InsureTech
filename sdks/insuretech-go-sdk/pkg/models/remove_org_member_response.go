package models


// RemoveOrgMemberResponse represents a remove_org_member_response
type RemoveOrgMemberResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
