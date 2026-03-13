package models


// PayoutCreationResponse represents a payout_creation_response
type PayoutCreationResponse struct {
	Error *Error `json:"error,omitempty"`
	PayoutId string `json:"payout_id,omitempty"`
	PayoutNumber string `json:"payout_number,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	CommissionCount int `json:"commission_count,omitempty"`
}
