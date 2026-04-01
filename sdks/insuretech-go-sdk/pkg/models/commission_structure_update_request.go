package models


// CommissionStructureUpdateRequest represents a commission_structure_update_request
type CommissionStructureUpdateRequest struct {
	CommissionRates map[string]interface{} `json:"commission_rates,omitempty"`
	PartnerId string `json:"partner_id"`
}
