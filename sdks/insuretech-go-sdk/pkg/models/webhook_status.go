package models

// WebhookStatus represents a webhook_status
type WebhookStatus string

// WebhookStatus values
const (
	WebhookStatusWEBHOOKSTATUSUNSPECIFIED WebhookStatus = "WEBHOOK_STATUS_UNSPECIFIED"
	WebhookStatusWEBHOOKSTATUSPENDING  = "WEBHOOK_STATUS_PENDING"
	WebhookStatusWEBHOOKSTATUSPROCESSED  = "WEBHOOK_STATUS_PROCESSED"
	WebhookStatusWEBHOOKSTATUSFAILED  = "WEBHOOK_STATUS_FAILED"
)
