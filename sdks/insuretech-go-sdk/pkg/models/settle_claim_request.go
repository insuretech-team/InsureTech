package models


// SettleClaimRequest represents a settle_claim_request
type SettleClaimRequest struct {
	ClaimId string `json:"claim_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
}
