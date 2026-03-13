package models


// GatewayWebhookHandlingResponse represents a gateway_webhook_handling_response
type GatewayWebhookHandlingResponse struct {
	Accepted bool `json:"accepted,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	Status string `json:"status,omitempty"`
	Error *Error `json:"error,omitempty"`
}
