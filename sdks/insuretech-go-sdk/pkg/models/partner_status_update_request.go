package models


// PartnerStatusUpdateRequest represents a partner_status_update_request
type PartnerStatusUpdateRequest struct {
	PartnerId string `json:"partner_id"`
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
}
