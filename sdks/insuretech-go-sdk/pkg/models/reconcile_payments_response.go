package models


// ReconcilePaymentsResponse represents a reconcile_payments_response
type ReconcilePaymentsResponse struct {
	MismatchPaymentIds []string `json:"mismatch_payment_ids,omitempty"`
	Error *Error `json:"error,omitempty"`
	TotalPayments int `json:"total_payments,omitempty"`
	ReconciledCount int `json:"reconciled_count,omitempty"`
	MismatchCount int `json:"mismatch_count,omitempty"`
}
