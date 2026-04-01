package models

import (
	"time"
)

// Payment represents a payment
type Payment struct {
	Amount *Money `json:"amount"`
	BankTranId string `json:"bank_tran_id,omitempty"`
	CallbackReceivedAt time.Time `json:"callback_received_at,omitempty"`
	CardBrand string `json:"card_brand,omitempty"`
	CardIssuer string `json:"card_issuer,omitempty"`
	CardIssuerCountry string `json:"card_issuer_country,omitempty"`
	CardType string `json:"card_type,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Currency string `json:"currency"`
	CustomerId string `json:"customer_id,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	GatewayResponse string `json:"gateway_response,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	InitiatedAt time.Time `json:"initiated_at"`
	InvoiceId string `json:"invoice_id,omitempty"`
	IpnReceivedAt time.Time `json:"ipn_received_at,omitempty"`
	LedgerTransactionId string `json:"ledger_transaction_id,omitempty"`
	ManualProofFileId string `json:"manual_proof_file_id,omitempty"`
	ManualReviewStatus interface{} `json:"manual_review_status"`
	Method *PaymentMethod `json:"method"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PayeeId string `json:"payee_id,omitempty"`
	PayerId string `json:"payer_id"`
	PaymentId string `json:"payment_id"`
	PolicyId string `json:"policy_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	ProviderReference string `json:"provider_reference,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	ReceiptDocumentId string `json:"receipt_document_id,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	ReceiptUrl string `json:"receipt_url,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	RetryCount int `json:"retry_count"`
	RiskLevel string `json:"risk_level,omitempty"`
	RiskTitle string `json:"risk_title,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	Status interface{} `json:"status"`
	TenantId string `json:"tenant_id,omitempty"`
	TigerbeetleTransferId string `json:"tigerbeetle_transfer_id,omitempty"`
	TranId string `json:"tran_id,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	Type *PaymentType `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	ValId string `json:"val_id,omitempty"`
	ValidatedAt time.Time `json:"validated_at,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
