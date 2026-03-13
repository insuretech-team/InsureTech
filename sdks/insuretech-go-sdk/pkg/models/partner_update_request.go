package models


// PartnerUpdateRequest represents a partner_update_request
type PartnerUpdateRequest struct {
	Partner *Partner `json:"partner,omitempty"`
	UpdateMask []string `json:"update_mask,omitempty"`
	PartnerId string `json:"partner_id"`
}
