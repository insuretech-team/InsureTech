package models


// ReconcilePaymentsResponse represents a reconcile_payments_response
type ReconcilePaymentsResponse struct {
	MismatchCount int `json:"mismatch_count,omitempty"`
	MismatchPaymentIds []string `json:"mismatch_payment_ids,omitempty"`
	ReconciledCount int `json:"reconciled_count,omitempty"`
	TotalPayments int `json:"total_payments,omitempty"`
}
