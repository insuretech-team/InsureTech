package models


// CommissionStatementRetrievalRequest represents a commission_statement_retrieval_request
type CommissionStatementRetrievalRequest struct {
	RecipientId string `json:"recipient_id"`
	PeriodStart string `json:"period_start,omitempty"`
	PeriodEnd string `json:"period_end,omitempty"`
}
