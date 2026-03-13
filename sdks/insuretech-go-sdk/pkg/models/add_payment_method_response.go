package models


// AddPaymentMethodResponse represents a add_payment_method_response
type AddPaymentMethodResponse struct {
	MethodId string `json:"method_id,omitempty"`
	Success bool `json:"success,omitempty"`
	Error *Error `json:"error,omitempty"`
}
