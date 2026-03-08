package events

// Kafka topic constants for order domain events.
// Canonical format: insuretech.<domain>.v1.<entity>.<verb>
const (
	// Outbound — published by orders-service
	TopicOrderCreated              = "insuretech.orders.v1.order.created"
	TopicOrderPaymentInitiated     = "insuretech.orders.v1.order.payment_initiated"
	TopicOrderPaymentConfirmed     = "insuretech.orders.v1.order.payment_confirmed"
	TopicOrderCancelled            = "insuretech.orders.v1.order.cancelled"
	TopicOrderFailed               = "insuretech.orders.v1.order.failed"
	TopicOrderFulfillmentCompleted = "insuretech.orders.v1.order.fulfillment_completed"

	// Inbound — payment events consumed by orders-service
	TopicPaymentCompleted       = "insuretech.payment.v1.payment.completed"
	TopicPaymentFailed          = "insuretech.payment.v1.payment.failed"
	TopicPaymentVerified        = "insuretech.payment.v1.payment.verified"         // manual review approved → re-trigger order
	TopicManualReviewRequested  = "insuretech.payment.v1.payment.manual_review_requested" // proof submitted → hold order
	TopicManualPaymentReviewed  = "insuretech.payment.v1.payment.manual_review_completed" // reviewed (approve/reject)

	// Inbound — insurance events consumed by orders-service
	TopicPolicyIssued = "insuretech.insurance.v1.policy.issued"

	// Inbound — B2B events consumed by orders-service
	// The B2B service publishes purchase_order.status_changed for all status transitions.
	// The consumer filters to APPROVED transitions only.
	TopicB2BPurchaseOrderApproved = "insuretech.b2b.v1.purchase_order.status_changed"
)
