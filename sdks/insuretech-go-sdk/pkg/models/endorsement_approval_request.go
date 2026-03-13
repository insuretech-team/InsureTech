package models


// EndorsementApprovalRequest represents a endorsement_approval_request
type EndorsementApprovalRequest struct {
	EndorsementId string `json:"endorsement_id"`
	ApprovedBy string `json:"approved_by,omitempty"`
	Comments string `json:"comments,omitempty"`
}
