package models

import (
	"time"
)

// QuoteStatisticsRetrievalRequest represents a quote_statistics_retrieval_request
type QuoteStatisticsRetrievalRequest struct {
	AgentId string `json:"agent_id"`
	CustomerId string `json:"customer_id"`
	EndDate time.Time `json:"end_date,omitempty"`
	ProductId string `json:"product_id"`
	StartDate time.Time `json:"start_date,omitempty"`
}
