package models


// CommissionStatementRetrievalResponse represents a commission_statement_retrieval_response
type CommissionStatementRetrievalResponse struct {
	ByType []*CommissionSummary `json:"by_type,omitempty"`
	PendingAmount *Money `json:"pending_amount,omitempty"`
	PeriodEnd string `json:"period_end,omitempty"`
	PeriodStart string `json:"period_start,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	TotalEarned *Money `json:"total_earned,omitempty"`
	TotalPaid *Money `json:"total_paid,omitempty"`
}
