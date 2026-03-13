package models


// PayoutCreationRequest represents a payout_creation_request
type PayoutCreationRequest struct {
	RecipientType string `json:"recipient_type,omitempty"`
	RecipientId string `json:"recipient_id"`
	PeriodStart string `json:"period_start,omitempty"`
	PeriodEnd string `json:"period_end,omitempty"`
	CommissionIds []string `json:"commission_ids,omitempty"`
}
