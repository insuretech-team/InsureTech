package models


// PayoutCreationRequest represents a payout_creation_request
type PayoutCreationRequest struct {
	CommissionIds []string `json:"commission_ids,omitempty"`
	PeriodEnd string `json:"period_end,omitempty"`
	PeriodStart string `json:"period_start,omitempty"`
	RecipientId string `json:"recipient_id"`
	RecipientType string `json:"recipient_type,omitempty"`
}
