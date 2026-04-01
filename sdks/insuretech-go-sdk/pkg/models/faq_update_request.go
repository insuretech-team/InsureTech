package models


// FAQUpdateRequest represents a faq_update_request
type FAQUpdateRequest struct {
	Faq *FAQ `json:"faq,omitempty"`
	FaqId string `json:"faq_id"`
}
