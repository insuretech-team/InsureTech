package models


// PartnerStatusUpdateResponse represents a partner_status_update_response
type PartnerStatusUpdateResponse struct {
	Partner *Partner `json:"partner,omitempty"`
	Error *Error `json:"error,omitempty"`
}
