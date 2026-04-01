package models


// QuoteStatisticsRetrievalResponse represents a quote_statistics_retrieval_response
type QuoteStatisticsRetrievalResponse struct {
	AcceptedQuotes int `json:"accepted_quotes,omitempty"`
	AveragePremium *Money `json:"average_premium,omitempty"`
	ConversionRate float64 `json:"conversion_rate,omitempty"`
	ConvertedQuotes int `json:"converted_quotes,omitempty"`
	DeclinedQuotes int `json:"declined_quotes,omitempty"`
	DraftQuotes int `json:"draft_quotes,omitempty"`
	ExpiredQuotes int `json:"expired_quotes,omitempty"`
	SentQuotes int `json:"sent_quotes,omitempty"`
	TotalPremiumValue *Money `json:"total_premium_value,omitempty"`
	TotalQuotes int `json:"total_quotes,omitempty"`
}
