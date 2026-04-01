package events

// Kafka topic constants for workflow domain events.
// Canonical format: insuretech.<domain>.v1.<entity>.<verb>
const (
	// Outbound — published by workflow-engine
	TopicWorkflowStarted      = "insuretech.workflow.v1.instance.started"
	TopicWorkflowTaskAssigned = "insuretech.workflow.v1.task.assigned"
	TopicWorkflowTaskCompleted = "insuretech.workflow.v1.task.completed"
	TopicWorkflowCompleted    = "insuretech.workflow.v1.instance.completed"

	// Inbound — consumed to trigger workflow actions
	// e.g. claim filed → start claim approval workflow
	TopicClaimFiled          = "insuretech.claims.v1.claim.filed"
	TopicEndorsementRequested = "insuretech.endorsement.v1.endorsement.requested"
	TopicRefundRequested     = "insuretech.refund.v1.refund.requested"
	TopicQuotationSubmitted  = "insuretech.quotation.submitted.v1"
)
