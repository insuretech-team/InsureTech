package models


// FAQUpdateRequest represents a faq_update_request
type FAQUpdateRequest struct {
	FaqId string `json:"faq_id"`
	Faq *FAQ `json:"faq,omitempty"`
}
