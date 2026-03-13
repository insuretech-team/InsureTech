package models


// PartnerRetrievalResponse represents a partner_retrieval_response
type PartnerRetrievalResponse struct {
	Partner *Partner `json:"partner,omitempty"`
	Error *Error `json:"error,omitempty"`
}
