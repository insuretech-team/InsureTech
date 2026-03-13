package models


// CommissionStructureUpdateRequest represents a commission_structure_update_request
type CommissionStructureUpdateRequest struct {
	PartnerId string `json:"partner_id"`
	CommissionRates map[string]interface{} `json:"commission_rates,omitempty"`
}
