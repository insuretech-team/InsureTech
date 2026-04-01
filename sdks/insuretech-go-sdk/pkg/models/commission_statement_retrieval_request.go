package models


// CommissionStatementRetrievalRequest represents a commission_statement_retrieval_request
type CommissionStatementRetrievalRequest struct {
	PeriodEnd string `json:"period_end,omitempty"`
	PeriodStart string `json:"period_start,omitempty"`
	RecipientId string `json:"recipient_id"`
}
