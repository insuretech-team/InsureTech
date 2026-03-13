package models


// UnderwritingApprovalResponse represents a underwriting_approval_response
type UnderwritingApprovalResponse struct {
	DecisionId string `json:"decision_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
