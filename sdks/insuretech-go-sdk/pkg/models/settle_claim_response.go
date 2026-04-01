package models


// SettleClaimResponse represents a settle_claim_response
type SettleClaimResponse struct {
	PaymentId string `json:"payment_id,omitempty"`
	SettledAmount *Money `json:"settled_amount,omitempty"`
}
