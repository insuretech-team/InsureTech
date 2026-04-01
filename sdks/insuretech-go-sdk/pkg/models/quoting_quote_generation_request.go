package models


// QuotingQuoteGenerationRequest represents a quoting_quote_generation_request
type QuotingQuoteGenerationRequest struct {
	AgentId string `json:"agent_id"`
	CustomerId string `json:"customer_id"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Parameters *QuoteParameters `json:"parameters,omitempty"`
	ProductId string `json:"product_id"`
	ValidityDays int `json:"validity_days,omitempty"`
}
