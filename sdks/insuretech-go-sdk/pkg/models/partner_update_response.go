package models


// PartnerUpdateResponse represents a partner_update_response
type PartnerUpdateResponse struct {
	Partner *Partner `json:"partner,omitempty"`
	Error *Error `json:"error,omitempty"`
}
