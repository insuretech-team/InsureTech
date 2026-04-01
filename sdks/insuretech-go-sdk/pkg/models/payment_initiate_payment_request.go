package models


// PaymentInitiatePaymentRequest represents a payment_initiate_payment_request
type PaymentInitiatePaymentRequest struct {
	Amount *Money `json:"amount,omitempty"`
	CallbackUrl string `json:"callback_url,omitempty"`
	Currency string `json:"currency,omitempty"`
	CustomerAddressLine1 string `json:"customer_address_line1,omitempty"`
	CustomerCity string `json:"customer_city,omitempty"`
	CustomerCountry string `json:"customer_country,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerId string `json:"customer_id"`
	CustomerName string `json:"customer_name,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	CustomerPostcode string `json:"customer_postcode,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	InvoiceId string `json:"invoice_id"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	OrderId string `json:"order_id"`
	OrganisationId string `json:"organisation_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PolicyId string `json:"policy_id"`
	PurchaseOrderId string `json:"purchase_order_id"`
	TenantId string `json:"tenant_id"`
	UserId string `json:"user_id"`
}
