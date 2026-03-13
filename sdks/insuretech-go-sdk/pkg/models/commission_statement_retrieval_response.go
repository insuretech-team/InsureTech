package models


// CommissionStatementRetrievalResponse represents a commission_statement_retrieval_response
type CommissionStatementRetrievalResponse struct {
	PendingAmount *Money `json:"pending_amount,omitempty"`
	ByType []*CommissionSummary `json:"by_type,omitempty"`
	Error *Error `json:"error,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	PeriodStart string `json:"period_start,omitempty"`
	PeriodEnd string `json:"period_end,omitempty"`
	TotalEarned *Money `json:"total_earned,omitempty"`
	TotalPaid *Money `json:"total_paid,omitempty"`
}
