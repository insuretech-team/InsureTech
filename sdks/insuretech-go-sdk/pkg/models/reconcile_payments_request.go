package models

import (
	"time"
)

// ReconcilePaymentsRequest represents a reconcile_payments_request
type ReconcilePaymentsRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Provider string `json:"provider,omitempty"`
}
