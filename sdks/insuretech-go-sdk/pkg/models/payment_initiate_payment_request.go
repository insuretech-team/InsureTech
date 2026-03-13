package models


// PaymentInitiatePaymentRequest represents a payment_initiate_payment_request
type PaymentInitiatePaymentRequest struct {
	Currency string `json:"currency,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	TenantId string `json:"tenant_id"`
	CustomerId string `json:"customer_id"`
	CustomerAddressLine1 string `json:"customer_address_line1,omitempty"`
	CustomerCountry string `json:"customer_country,omitempty"`
	PolicyId string `json:"policy_id"`
	CallbackUrl string `json:"callback_url,omitempty"`
	OrderId string `json:"order_id"`
	OrganisationId string `json:"organisation_id"`
	CustomerName string `json:"customer_name,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	UserId string `json:"user_id"`
	InvoiceId string `json:"invoice_id"`
	PurchaseOrderId string `json:"purchase_order_id"`
	CustomerCity string `json:"customer_city,omitempty"`
	CustomerPostcode string `json:"customer_postcode,omitempty"`
}
