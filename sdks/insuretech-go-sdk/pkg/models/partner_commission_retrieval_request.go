package models

import (
	"time"
)

// PartnerCommissionRetrievalRequest represents a partner_commission_retrieval_request
type PartnerCommissionRetrievalRequest struct {
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate time.Time `json:"end_date,omitempty"`
	PartnerId string `json:"partner_id"`
}
