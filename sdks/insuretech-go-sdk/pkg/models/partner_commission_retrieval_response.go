package models


// PartnerCommissionRetrievalResponse represents a partner_commission_retrieval_response
type PartnerCommissionRetrievalResponse struct {
	Currency string `json:"currency,omitempty"`
	Details []*CommissionDetail `json:"details,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	TotalCommission *Money `json:"total_commission,omitempty"`
}
