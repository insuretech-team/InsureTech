package models


// PartnerCommissionRetrievalResponse represents a partner_commission_retrieval_response
type PartnerCommissionRetrievalResponse struct {
	TotalCommission *Money `json:"total_commission,omitempty"`
	Currency string `json:"currency,omitempty"`
	Details []*CommissionDetail `json:"details,omitempty"`
	Error *Error `json:"error,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
}
