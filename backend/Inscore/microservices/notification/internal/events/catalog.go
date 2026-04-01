package events

import (
	"sort"
	"strings"
)

type TopicDirection string
type TopicGroup string
type AudienceScope string
type SubscriptionProfile string

const (
	DirectionInbound  TopicDirection = "inbound"
	DirectionOutbound TopicDirection = "outbound"

	AudienceCustomer AudienceScope = "customer"
	AudiencePartner  AudienceScope = "partner"
	AudienceInsurer  AudienceScope = "insurer"
	AudienceInternal AudienceScope = "internal"
	AudienceExternal AudienceScope = "external"

	GroupCustomerIdentity  TopicGroup = "customer_identity"
	GroupCustomerCommerce  TopicGroup = "customer_commerce"
	GroupClaimsLifecycle   TopicGroup = "claims_lifecycle"
	GroupRenewalRetention  TopicGroup = "renewal_retention"
	GroupDocumentArtifacts TopicGroup = "document_artifacts"
	GroupPartnerOps        TopicGroup = "partner_ops"
	GroupSupportEscalation TopicGroup = "support_escalation"
	GroupMarketingCampaign TopicGroup = "marketing_campaign"
	GroupIoTRiskAlerts     TopicGroup = "iot_risk_alerts"
	GroupComplianceOps     TopicGroup = "compliance_ops"
	GroupNotificationState TopicGroup = "notification_state"
	GroupWebhookFanout     TopicGroup = "webhook_fanout"

	ProfileTransactionalCore SubscriptionProfile = "transactional_core"
	ProfileCustomerCore      SubscriptionProfile = "customer_core"
	ProfileOperations        SubscriptionProfile = "operations"
	ProfilePlatformAll       SubscriptionProfile = "platform_all"
)

type TopicDescriptor struct {
	Key             string
	Topic           string
	Group           TopicGroup
	Direction       TopicDirection
	Audience        AudienceScope
	Description     string
	CurrentContract bool
	Reserved        bool
}

type TopicGroupSpec struct {
	Name        TopicGroup
	Description string
	Topics      []TopicDescriptor
}

type SubscriptionPlan struct {
	Profile               SubscriptionProfile
	EnabledGroups         []string
	DisabledGroups        []string
	AllowTopics           []string
	DenyTopics            []string
	ExtraTopics           []string
	IncludeReservedTopics bool
}

var (
	CustomerIdentity = TopicGroupSpec{
		Name:        GroupCustomerIdentity,
		Description: "Customer identity, verification, and account-safety events that may trigger transactional notifications or ops alerts.",
		Topics: []TopicDescriptor{
			{Key: "user_registered", Topic: "authn.user.registered", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Welcome and onboarding events after registration.", CurrentContract: true},
			{Key: "email_verified", Topic: "authn.email.verified", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Email verification completed and confirmation can be delivered.", CurrentContract: true},
			{Key: "password_changed", Topic: "authn.password.changed", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Security password change confirmation.", CurrentContract: true},
			{Key: "password_reset_requested", Topic: "authn.password.reset_requested", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Security alert for password reset by mobile OTP flow.", CurrentContract: true},
			{Key: "email_password_reset_requested", Topic: "authn.email.password_reset_requested", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Security alert for password reset initiated by email flow.", CurrentContract: true},
			{Key: "kyc_completed", Topic: "insuretech.authn.v1.kyc.completed", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "KYC result and onboarding progression.", Reserved: true},
			{Key: "account_locked", Topic: "authn.account.locked", Group: GroupCustomerIdentity, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Security lockouts for account-protection alerts.", CurrentContract: true},
		},
	}
	CustomerCommerce = TopicGroupSpec{
		Name:        GroupCustomerCommerce,
		Description: "Retail and B2B commerce lifecycle from order and invoice creation through payment and policy activation.",
		Topics: []TopicDescriptor{
			{Key: "order_payment_initiated", Topic: "insuretech.orders.v1.order.payment_initiated", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Customer-facing payment pending state.", CurrentContract: true},
			{Key: "order_created", Topic: "insuretech.orders.v1.order.created", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Order acknowledged and awaiting payment/fulfillment progression.", CurrentContract: true},
			{Key: "order_payment_confirmed", Topic: "insuretech.orders.v1.order.payment_confirmed", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Order payment confirmed and issuance/fulfillment will proceed.", CurrentContract: true},
			{Key: "order_cancelled", Topic: "insuretech.orders.v1.order.cancelled", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Order cancelled and customer should be informed.", CurrentContract: true},
			{Key: "order_failed", Topic: "insuretech.orders.v1.order.failed", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Order failed and the customer may need to retry.", CurrentContract: true},
			{Key: "invoice_issued", Topic: "insuretech.billing.v1.invoice.issued", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Invoice/payment artifact is ready.", Reserved: true},
			{Key: "invoice_paid", Topic: "insuretech.billing.v1.invoice.paid", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Commercial settlement confirmed.", Reserved: true},
			{Key: "invoice_overdue", Topic: "insuretech.billing.v1.invoice.overdue", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Premium overdue and billing follow-up.", Reserved: true},
			{Key: "payment_completed", Topic: "insuretech.payment.v1.payment.completed", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Payment settled successfully.", CurrentContract: true},
			{Key: "payment_failed", Topic: "insuretech.payment.v1.payment.failed", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Payment attempt failed and may need retry messaging.", CurrentContract: true},
			{Key: "refund_processed", Topic: "insuretech.payment.v1.refund.processed", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Refund settlement after cancellation or payment reversal.", CurrentContract: true},
			{Key: "policy_issued", Topic: "insuretech.insurance.v1.policy.issued", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Coverage is active and policy pack generation can begin.", CurrentContract: true},
			{Key: "policy_cancelled", Topic: "insuretech.insurance.v1.policy.cancelled", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Policy cancellation and refund communication.", Reserved: true},
			{Key: "policy_renewed", Topic: "insuretech.insurance.v1.policy.renewed", Group: GroupCustomerCommerce, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Successful renewal and refreshed coverage.", Reserved: true},
		},
	}
	ClaimsLifecycle = TopicGroupSpec{
		Name:        GroupClaimsLifecycle,
		Description: "End-to-end claims status transitions for customer, partner, and internal stakeholders.",
		Topics: []TopicDescriptor{
			{Key: "claim_submitted", Topic: "insuretech.claims.v1.claim.submitted", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Claim submission acknowledgment.", CurrentContract: true},
			{Key: "claim_documents_requested", Topic: "insuretech.claims.v1.claim.documents_requested", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Additional documents requested for claim processing.", Reserved: true},
			{Key: "claim_under_review", Topic: "insuretech.claims.v1.claim.under_review", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Claim moved to under-review state.", Reserved: true},
			{Key: "claim_approved", Topic: "insuretech.claims.v1.claim.approved", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Claim approved communication.", CurrentContract: true},
			{Key: "claim_rejected", Topic: "insuretech.claims.v1.claim.rejected", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Claim rejection communication.", CurrentContract: true},
			{Key: "claim_settled", Topic: "insuretech.claims.v1.claim.settled", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Claim settled and payout confirmation.", CurrentContract: true},
			{Key: "claim_fraud_detected", Topic: "insuretech.claims.v1.claim.fraud_detected", Group: GroupClaimsLifecycle, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Internal-only fraud escalation signals.", CurrentContract: true},
		},
	}
	RenewalRetention = TopicGroupSpec{
		Name:        GroupRenewalRetention,
		Description: "Retention, reminder, grace-period, lapse, and reinstatement journeys.",
		Topics: []TopicDescriptor{
			{Key: "renewal_due", Topic: "insuretech.renewal.v1.renewal.due", Group: GroupRenewalRetention, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Scheduled renewal reminder source.", CurrentContract: true},
			{Key: "grace_period_started", Topic: "insuretech.renewal.v1.grace_period.started", Group: GroupRenewalRetention, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Grace-period daily reminder journey starts.", CurrentContract: true},
			{Key: "policy_lapsed", Topic: "insuretech.renewal.v1.policy.lapsed", Group: GroupRenewalRetention, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Policy lapsed after grace period.", CurrentContract: true},
			{Key: "policy_lapsed_bridge", Topic: "insuretech.policy.v1.policy.lapsed", Group: GroupRenewalRetention, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Bridge topic while policy lifecycle producers standardize.", CurrentContract: true},
			{Key: "policy_reinstated", Topic: "insuretech.insurance.v1.policy.reinstated", Group: GroupRenewalRetention, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Reinstatement success notice.", Reserved: true},
		},
	}
	DocumentArtifacts = TopicGroupSpec{
		Name:        GroupDocumentArtifacts,
		Description: "Policy packs, receipts, claim attachments, and stored artifacts that can trigger download/link notifications.",
		Topics: []TopicDescriptor{
			{Key: "document_generated", Topic: "insuretech.document.v1.document.generated", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Document generation completed.", Reserved: true},
			{Key: "document_failed", Topic: "insuretech.document.v1.document.failed", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Document generation failure requiring follow-up.", Reserved: true},
			{Key: "storage_file_finalized", Topic: "insuretech.storage.v1.file.finalized", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Uploaded file is durable and safe to reference.", Reserved: true},
			{Key: "receipt_generated", Topic: "insuretech.payment.v1.payment.receipt_generated", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Payment receipt available for delivery.", CurrentContract: true},
			{Key: "storage_events", Topic: "storage.events", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Aggregate storage lifecycle stream for file upload, finalize, metadata, and deletion events.", CurrentContract: true},
			{Key: "docgen_events", Topic: "docgen-events", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Aggregate document generation stream from the docgen service.", CurrentContract: true},
			{Key: "media_events", Topic: "media-events", Group: GroupDocumentArtifacts, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Aggregate media processing stream for uploads and processing lifecycle.", CurrentContract: true},
		},
	}
	PartnerOps = TopicGroupSpec{
		Name:        GroupPartnerOps,
		Description: "Partner, agent, employer, and insurer coordination events in a multi-insurer distribution network.",
		Topics: []TopicDescriptor{
			{Key: "organisation_approved", Topic: "b2b.organisation.approved", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "Organisation approved and B2B admins should be informed.", CurrentContract: true},
			{Key: "organisation_suspended", Topic: "b2b.organisation.suspended", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "Organisation suspended and B2B admins should be informed.", CurrentContract: true},
			{Key: "org_member_added", Topic: "b2b.org_member.added", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "A user was added to an organisation.", CurrentContract: true},
			{Key: "org_member_removed", Topic: "b2b.org_member.removed", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "A user was removed from an organisation.", CurrentContract: true},
			{Key: "b2b_admin_assigned", Topic: "b2b.admin.assigned", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "A B2B admin role was assigned to a user.", CurrentContract: true},
			{Key: "partner_events", Topic: "partner-events", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "Aggregate partner distribution stream for partner, agent, and commission events.", CurrentContract: true},
			{Key: "purchase_order_approved", Topic: "insuretech.b2b.v1.purchase_order.approved", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "Corporate purchase order progression.", Reserved: true},
			{Key: "partner_suspended", Topic: "insuretech.partner.v1.partner.suspended", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudiencePartner, Description: "Partner suspension with grace handover messaging.", Reserved: true},
			{Key: "insurer_quote_queued", Topic: "insuretech.insurer.v1.quote.queued", Group: GroupPartnerOps, Direction: DirectionInbound, Audience: AudienceInsurer, Description: "Insurer API fallback queue updates.", Reserved: true},
		},
	}
	SupportEscalation = TopicGroupSpec{
		Name:        GroupSupportEscalation,
		Description: "Customer support, task management, and SLA/escalation workflows.",
		Topics: []TopicDescriptor{
			{Key: "authz_events", Topic: "authz.events", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Aggregate authz control-plane stream for role changes and portal policy updates.", CurrentContract: true},
			{Key: "workflow_task_assigned", Topic: "insuretech.workflow.v1.task.assigned", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Internal workflow task assignment notifications.", CurrentContract: true},
			{Key: "workflow_completed", Topic: "insuretech.workflow.v1.instance.completed", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Workflow completion updates that can be fanned out to stakeholders.", CurrentContract: true},
			{Key: "ticket_created", Topic: "insuretech.support.v1.ticket.created", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Ticket acknowledgment and intake confirmation.", Reserved: true},
			{Key: "ticket_updated", Topic: "insuretech.support.v1.ticket.updated", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Support update notifications.", Reserved: true},
			{Key: "task_overdue", Topic: "insuretech.tasks.v1.task.overdue", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Internal SLA breach and overdue work items.", Reserved: true},
			{Key: "escalation_triggered", Topic: "insuretech.support.v1.escalation.triggered", Group: GroupSupportEscalation, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Tiered escalation notifications.", Reserved: true},
		},
	}
	MarketingCampaign = TopicGroupSpec{
		Name:        GroupMarketingCampaign,
		Description: "Partner-branded campaign journeys, segmentation, and promotional fanout.",
		Topics: []TopicDescriptor{
			{Key: "campaign_dispatch_requested", Topic: "insuretech.marketing.v1.campaign.dispatch_requested", Group: GroupMarketingCampaign, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Approved marketing campaign dispatch.", Reserved: true},
			{Key: "audience_segment_ready", Topic: "insuretech.marketing.v1.audience.segment_ready", Group: GroupMarketingCampaign, Direction: DirectionInbound, Audience: AudiencePartner, Description: "Prepared audience segments for partner messaging.", Reserved: true},
		},
	}
	IoTRiskAlerts = TopicGroupSpec{
		Name:        GroupIoTRiskAlerts,
		Description: "Telematics, device, and real-time risk-alert journeys for motor, health, and home products.",
		Topics: []TopicDescriptor{
			{Key: "iot_alert_threshold_breached", Topic: "insuretech.iot.v1.alert.threshold_breached", Group: GroupIoTRiskAlerts, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Immediate device/risk threshold alerts.", Reserved: true},
			{Key: "iot_risk_score_updated", Topic: "insuretech.iot.v1.risk.score_updated", Group: GroupIoTRiskAlerts, Direction: DirectionInbound, Audience: AudienceCustomer, Description: "Usage-based insurance risk updates.", Reserved: true},
		},
	}
	ComplianceOps = TopicGroupSpec{
		Name:        GroupComplianceOps,
		Description: "Fraud, compliance, AML, regulator, and internal-control notifications.",
		Topics: []TopicDescriptor{
			{Key: "aml_flagged", Topic: "insuretech.compliance.v1.aml.flagged", Group: GroupComplianceOps, Direction: DirectionInbound, Audience: AudienceInternal, Description: "AML review and regulator-facing follow-up.", Reserved: true},
			{Key: "fraud_alert_opened", Topic: "insuretech.fraud.v1.alert.opened", Group: GroupComplianceOps, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Fraud alert routing for investigators.", Reserved: true},
			{Key: "incident_escalated", Topic: "insuretech.ops.v1.incident.escalated", Group: GroupComplianceOps, Direction: DirectionInbound, Audience: AudienceInternal, Description: "Operational incident escalation events.", Reserved: true},
		},
	}
	NotificationState = TopicGroupSpec{
		Name:        GroupNotificationState,
		Description: "Internal notification lifecycle events emitted by this service and consumed by downstream analytics, webhook fanout, or audit processors.",
		Topics: []TopicDescriptor{
			{Key: "notification_requested", Topic: "insuretech.notifications.v1.notification.requested", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "A normalized notification request has been accepted by the service."},
			{Key: "notification_queued", Topic: "insuretech.notifications.v1.notification.queued", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "Notification persisted and queued for dispatch."},
			{Key: "notification_sent", Topic: "insuretech.notifications.v1.notification.sent", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "Provider accepted or in-app delivery completed."},
			{Key: "notification_delivered", Topic: "insuretech.notifications.v1.notification.delivered", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "Terminal delivery confirmation when available."},
			{Key: "notification_failed", Topic: "insuretech.notifications.v1.notification.failed", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "Notification exhausted retries or encountered terminal failure."},
			{Key: "notification_read", Topic: "insuretech.notifications.v1.notification.read", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "Customer read or acknowledged an in-app notification."},
			{Key: "preferences_updated", Topic: "insuretech.notifications.v1.preferences.updated", Group: GroupNotificationState, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "User notification preference change event."},
		},
	}
	WebhookFanout = TopicGroupSpec{
		Name:        GroupWebhookFanout,
		Description: "External-system notification fanout and delivery lifecycle for partners, insurers, and ecosystem integrations.",
		Topics: []TopicDescriptor{
			{Key: "webhook_delivery_requested", Topic: "insuretech.integration.v1.webhook.delivery_requested", Group: GroupWebhookFanout, Direction: DirectionOutbound, Audience: AudienceExternal, Description: "A downstream webhook should be attempted for an external subscriber."},
			{Key: "webhook_delivery_failed", Topic: "insuretech.integration.v1.webhook.delivery_failed", Group: GroupWebhookFanout, Direction: DirectionOutbound, Audience: AudienceInternal, Description: "Webhook retry or manual intervention needed."},
		},
	}
)

func AllTopicGroups() []TopicGroupSpec {
	return []TopicGroupSpec{
		CustomerIdentity,
		CustomerCommerce,
		ClaimsLifecycle,
		RenewalRetention,
		DocumentArtifacts,
		PartnerOps,
		SupportEscalation,
		MarketingCampaign,
		IoTRiskAlerts,
		ComplianceOps,
		NotificationState,
		WebhookFanout,
	}
}

func DescriptorForTopic(topic string) (TopicDescriptor, bool) {
	for _, group := range AllTopicGroups() {
		for _, descriptor := range group.Topics {
			if descriptor.Topic == topic {
				return descriptor, true
			}
		}
	}
	return TopicDescriptor{}, false
}

func DescriptorForKey(key string) (TopicDescriptor, bool) {
	for _, group := range AllTopicGroups() {
		for _, descriptor := range group.Topics {
			if descriptor.Key == key {
				return descriptor, true
			}
		}
	}
	return TopicDescriptor{}, false
}

func ConsumerTopicsForPlan(plan SubscriptionPlan) []string {
	descriptors := ResolveSubscriptionPlan(plan)
	topics := make([]string, 0, len(descriptors))
	for _, descriptor := range descriptors {
		topics = append(topics, descriptor.Topic)
	}
	sort.Strings(topics)
	return topics
}

func ResolveSubscriptionPlan(plan SubscriptionPlan) []TopicDescriptor {
	if plan.Profile == "" {
		plan.Profile = ProfileCustomerCore
	}

	enabled := make(map[TopicGroup]struct{})
	for _, group := range defaultGroupsForProfile(plan.Profile) {
		enabled[group] = struct{}{}
	}
	for _, group := range plan.EnabledGroups {
		if parsed := parseGroup(group); parsed != "" {
			enabled[parsed] = struct{}{}
		}
	}
	for _, group := range plan.DisabledGroups {
		delete(enabled, parseGroup(group))
	}

	allowlist := make(map[string]struct{}, len(plan.AllowTopics))
	for _, topic := range plan.AllowTopics {
		if trimmed := strings.TrimSpace(topic); trimmed != "" {
			allowlist[trimmed] = struct{}{}
		}
	}
	denylist := make(map[string]struct{}, len(plan.DenyTopics))
	for _, topic := range plan.DenyTopics {
		if trimmed := strings.TrimSpace(topic); trimmed != "" {
			denylist[trimmed] = struct{}{}
		}
	}

	descriptorsByTopic := make(map[string]TopicDescriptor)
	for _, group := range AllTopicGroups() {
		if _, ok := enabled[group.Name]; !ok {
			continue
		}
		for _, descriptor := range group.Topics {
			if descriptor.Direction != DirectionInbound {
				continue
			}
			if descriptor.Reserved && !plan.IncludeReservedTopics {
				continue
			}
			if !descriptor.CurrentContract && !plan.IncludeReservedTopics {
				continue
			}
			if _, denied := denylist[descriptor.Topic]; denied {
				continue
			}
			descriptorsByTopic[descriptor.Topic] = descriptor
		}
	}

	for topic := range allowlist {
		if descriptor, ok := DescriptorForTopic(topic); ok {
			descriptorsByTopic[topic] = descriptor
			continue
		}
		descriptorsByTopic[topic] = TopicDescriptor{
			Key:             sanitizeTopicKey(topic),
			Topic:           topic,
			Group:           GroupNotificationState,
			Direction:       DirectionInbound,
			Audience:        AudienceInternal,
			Description:     "Explicitly allowlisted custom notification source topic.",
			CurrentContract: true,
		}
	}
	for _, topic := range plan.ExtraTopics {
		trimmed := strings.TrimSpace(topic)
		if trimmed == "" {
			continue
		}
		if _, denied := denylist[trimmed]; denied {
			continue
		}
		if descriptor, ok := DescriptorForTopic(trimmed); ok {
			descriptorsByTopic[trimmed] = descriptor
			continue
		}
		descriptorsByTopic[trimmed] = TopicDescriptor{
			Key:             sanitizeTopicKey(trimmed),
			Topic:           trimmed,
			Group:           GroupNotificationState,
			Direction:       DirectionInbound,
			Audience:        AudienceInternal,
			Description:     "Custom extra topic injected through notification configuration.",
			CurrentContract: true,
		}
	}

	topics := make([]string, 0, len(descriptorsByTopic))
	for topic := range descriptorsByTopic {
		topics = append(topics, topic)
	}
	sort.Strings(topics)

	result := make([]TopicDescriptor, 0, len(topics))
	for _, topic := range topics {
		result = append(result, descriptorsByTopic[topic])
	}
	return result
}

func defaultGroupsForProfile(profile SubscriptionProfile) []TopicGroup {
	switch profile {
	case ProfileTransactionalCore:
		return []TopicGroup{
			GroupCustomerCommerce,
			GroupClaimsLifecycle,
			GroupRenewalRetention,
			GroupDocumentArtifacts,
		}
	case ProfileOperations:
		return []TopicGroup{
			GroupCustomerCommerce,
			GroupClaimsLifecycle,
			GroupPartnerOps,
			GroupSupportEscalation,
			GroupComplianceOps,
			GroupDocumentArtifacts,
		}
	case ProfilePlatformAll:
		return []TopicGroup{
			GroupCustomerIdentity,
			GroupCustomerCommerce,
			GroupClaimsLifecycle,
			GroupRenewalRetention,
			GroupDocumentArtifacts,
			GroupPartnerOps,
			GroupSupportEscalation,
			GroupMarketingCampaign,
			GroupIoTRiskAlerts,
			GroupComplianceOps,
		}
	default:
		return []TopicGroup{
			GroupCustomerIdentity,
			GroupCustomerCommerce,
			GroupClaimsLifecycle,
			GroupRenewalRetention,
			GroupDocumentArtifacts,
			GroupSupportEscalation,
		}
	}
}

func parseGroup(value string) TopicGroup {
	normalized := strings.TrimSpace(strings.ToLower(value))
	switch normalized {
	case string(GroupCustomerIdentity):
		return GroupCustomerIdentity
	case string(GroupCustomerCommerce):
		return GroupCustomerCommerce
	case string(GroupClaimsLifecycle):
		return GroupClaimsLifecycle
	case string(GroupRenewalRetention):
		return GroupRenewalRetention
	case string(GroupDocumentArtifacts):
		return GroupDocumentArtifacts
	case string(GroupPartnerOps):
		return GroupPartnerOps
	case string(GroupSupportEscalation):
		return GroupSupportEscalation
	case string(GroupMarketingCampaign):
		return GroupMarketingCampaign
	case string(GroupIoTRiskAlerts):
		return GroupIoTRiskAlerts
	case string(GroupComplianceOps):
		return GroupComplianceOps
	case string(GroupNotificationState):
		return GroupNotificationState
	case string(GroupWebhookFanout):
		return GroupWebhookFanout
	default:
		return ""
	}
}

func sanitizeTopicKey(topic string) string {
	replacer := strings.NewReplacer(".", "_", "-", "_")
	return replacer.Replace(strings.TrimSpace(topic))
}
