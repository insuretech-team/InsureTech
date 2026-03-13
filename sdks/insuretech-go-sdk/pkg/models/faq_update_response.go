package models


// FAQUpdateResponse represents a faq_update_response
type FAQUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
