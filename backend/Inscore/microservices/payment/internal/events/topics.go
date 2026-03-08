package events

// Kafka topic constants for payment domain events.
// Canonical format: insuretech.<domain>.v1.<entity>.<verb>
const (
	// Outbound — published by payment-service
	TopicPaymentInitiated           = "insuretech.payment.v1.payment.initiated"
	TopicPaymentCompleted           = "insuretech.payment.v1.payment.completed"
	TopicPaymentFailed              = "insuretech.payment.v1.payment.failed"
	TopicRefundProcessed            = "insuretech.payment.v1.refund.processed"
	TopicPaymentVerified            = "insuretech.payment.v1.payment.verified"
	TopicManualReviewRequested      = "insuretech.payment.v1.payment.manual_review_requested"
	TopicManualPaymentReviewed      = "insuretech.payment.v1.payment.manual_review_completed"
	TopicReceiptGenerated           = "insuretech.payment.v1.payment.receipt_generated"
	TopicReconciliationMismatch     = "insuretech.payment.v1.payment.reconciliation_mismatch"
)
