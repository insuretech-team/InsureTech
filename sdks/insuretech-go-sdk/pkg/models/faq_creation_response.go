package models


// FAQCreationResponse represents a faq_creation_response
type FAQCreationResponse struct {
	FaqId string `json:"faq_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
