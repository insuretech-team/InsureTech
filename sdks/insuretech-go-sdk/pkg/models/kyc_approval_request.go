package models


// KYCApprovalRequest represents a kyc_approval_request
type KYCApprovalRequest struct {
	KycId string `json:"kyc_id"`
	ReviewerId string `json:"reviewer_id"`
}
