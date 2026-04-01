package models


// PartnerStatusUpdateRequest represents a partner_status_update_request
type PartnerStatusUpdateRequest struct {
	PartnerId string `json:"partner_id"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status,omitempty"`
}
