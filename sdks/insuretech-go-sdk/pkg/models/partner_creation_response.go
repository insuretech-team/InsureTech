package models


// PartnerCreationResponse represents a partner_creation_response
type PartnerCreationResponse struct {
	PartnerId string `json:"partner_id,omitempty"`
	Partner *Partner `json:"partner,omitempty"`
	Error *Error `json:"error,omitempty"`
}
