package models


// PremiumCalculationsRetrievalRequest represents a premium_calculations_retrieval_request
type PremiumCalculationsRetrievalRequest struct {
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	RegistrationId string `json:"registration_id"`
}
