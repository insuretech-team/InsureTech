package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	authneventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/events/v1"
	authzeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/events/v1"
	b2beventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/events/v1"
	claimeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/claims/events/v1"
	mediaeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/events/v1"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	orderseventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/events/v1"
	partnereventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/events/v1"
	paymenteventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/events/v1"
	policyeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/events/v1"
	renewaleventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/renewal/events/v1"
	storageeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/events/v1"
	workfloweventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type NotificationEventService interface {
	QueueDomainNotification(ctx context.Context, input *DomainNotificationInput) error
}

type EventContext struct {
	Topic          string
	TopicGroup     TopicGroup
	Audience       AudienceScope
	CorrelationID  string
	CausationID    string
	TenantID       string
	OrganisationID string
	Portal         string
	ActorUserID    string
	TraceID        string
}

type DomainNotificationInput struct {
	RecipientID       string
	RecipientRefKind  string
	RecipientRefID    string
	Type              notificationv1.NotificationType
	Priority          notificationv1.NotificationPriority
	Subject           string
	Message           string
	TemplateData      map[string]string
	Source            EventContext
	PreferredChannels []notificationv1.NotificationChannel
}

type Consumer struct {
	service NotificationEventService
}

func NewConsumer(service NotificationEventService) *Consumer {
	return &Consumer{service: service}
}

func (c *Consumer) HandleMessage(ctx context.Context, topic string, payload []byte) error {
	if c == nil || c.service == nil || len(payload) == 0 {
		return nil
	}

	descriptor, _ := DescriptorForTopic(topic)

	switch topic {
	case mustTopic("user_registered"):
		var evt authneventsv1.UserRegisteredEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode user registered event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:     "Welcome to Labaid InsureTech",
			Message:     "Your account has been created successfully.",
			TemplateData: map[string]string{
				"user_id":       evt.GetUserId(),
				"mobile_number": evt.GetMobileNumber(),
				"email":         evt.GetEmail(),
				"portal":        evt.GetPortal(),
				"tenant_id":     evt.GetTenantId(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
				TenantID:   evt.GetTenantId(),
				Portal:     evt.GetPortal(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
			},
		})
	case mustTopic("email_verified"):
		var evt authneventsv1.EmailVerifiedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode email verified event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:     "Email verified",
			Message:     "Your email address has been verified successfully.",
			TemplateData: map[string]string{
				"user_id": evt.GetUserId(),
				"email":   evt.GetEmail(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("password_changed"):
		var evt authneventsv1.PasswordChangedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode password changed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:     "Password changed",
			Message:     "Your account password has been changed.",
			TemplateData: map[string]string{
				"user_id":    evt.GetUserId(),
				"ip_address": evt.GetIpAddress(),
				"changed_by": evt.GetChangedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("password_reset_requested"), mustTopic("email_password_reset_requested"):
		var evt authneventsv1.PasswordResetRequestedEvent
		if topic == mustTopic("email_password_reset_requested") {
			var emailEvt authneventsv1.PasswordResetByEmailRequestedEvent
			if err := protojson.Unmarshal(payload, &emailEvt); err != nil {
				return fmt.Errorf("decode email password reset requested event: %w", err)
			}
			return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
				RecipientID: emailEvt.GetUserId(),
				Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
				Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
				Subject:     "Password reset requested",
				Message:     "A password reset request was initiated for your account.",
				TemplateData: map[string]string{
					"user_id":      emailEvt.GetUserId(),
					"email_masked": emailEvt.GetEmailMasked(),
					"otp_id":       emailEvt.GetOtpId(),
					"ip_address":   emailEvt.GetIpAddress(),
				},
				Source: EventContext{
					Topic:      descriptor.Topic,
					TopicGroup: descriptor.Group,
					Audience:   descriptor.Audience,
				},
				PreferredChannels: []notificationv1.NotificationChannel{
					notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
					notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
					notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				},
			})
		}
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode password reset requested event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:     "Password reset requested",
			Message:     "A password reset request was initiated for your account.",
			TemplateData: map[string]string{
				"user_id":       evt.GetUserId(),
				"mobile_number": evt.GetMobileNumber(),
				"ip_address":    evt.GetIpAddress(),
				"device_type":   evt.GetDeviceType(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("account_locked"):
		var evt authneventsv1.AccountLockedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode account locked event: %w", err)
		}
		templateData := map[string]string{
			"user_id": evt.GetUserId(),
			"reason":  evt.GetReason(),
		}
		if evt.GetLockedUntil() != nil {
			templateData["locked_until"] = evt.GetLockedUntil().AsTime().Format(time.RFC3339)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:  evt.GetUserId(),
			Type:         notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:     notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT,
			Subject:      "Account locked",
			Message:      "Your account has been locked for security reasons.",
			TemplateData: templateData,
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("order_created"):
		var evt orderseventsv1.OrderCreatedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode order created event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "order",
			RecipientRefID:   evt.GetOrderId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Order created",
			Message:          "Your order has been created and is awaiting payment confirmation.",
			TemplateData: map[string]string{
				"order_id":      evt.GetOrderId(),
				"order_number":  evt.GetOrderNumber(),
				"quotation_id":  evt.GetQuotationId(),
				"total_payable": formatMoney(evt.GetTotalPayable()),
			},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
			},
		})
	case mustTopic("order_payment_confirmed"):
		var evt orderseventsv1.OrderPaymentConfirmedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode order payment confirmed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "order",
			RecipientRefID:   evt.GetOrderId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Order payment confirmed",
			Message:          "Your payment has been confirmed and policy issuance is now in progress.",
			TemplateData: map[string]string{
				"order_id":      evt.GetOrderId(),
				"payment_id":    evt.GetPaymentId(),
				"quotation_id":  evt.GetQuotationId(),
				"invoice_id":    evt.GetInvoiceId(),
				"total_payable": formatMoney(evt.GetTotalPayable()),
			},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("order_cancelled"):
		var evt orderseventsv1.OrderCancelledEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode order cancelled event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "order",
			RecipientRefID:   evt.GetOrderId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Order cancelled",
			Message:          "Your order has been cancelled.",
			TemplateData: map[string]string{
				"order_id": evt.GetOrderId(),
				"reason":   evt.GetReason(),
			},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("order_failed"):
		var evt orderseventsv1.OrderFailedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode order failed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "order",
			RecipientRefID:   evt.GetOrderId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Order failed",
			Message:          "We could not complete your order. Please review the failure details and try again.",
			TemplateData: map[string]string{
				"order_id":      evt.GetOrderId(),
				"payment_id":    evt.GetPaymentId(),
				"error_code":    evt.GetErrorCode(),
				"error_message": evt.GetErrorMessage(),
			},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("workflow_task_assigned"):
		var evt workfloweventsv1.WorkflowTaskAssignedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode workflow task assigned event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetAssignedTo(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:     "Workflow task assigned",
			Message:     "A workflow task has been assigned to you.",
			TemplateData: map[string]string{
				"task_id":              evt.GetTaskId(),
				"workflow_instance_id": evt.GetWorkflowInstanceId(),
				"step_name":            evt.GetStepName(),
				"assigned_to":          evt.GetAssignedTo(),
			},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
			},
		})
	case mustTopic("workflow_completed"):
		var evt workfloweventsv1.WorkflowCompletedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode workflow completed event: %w", err)
		}
		refKind := workflowEntityRecipientRef(evt.GetEntityType())
		if refKind == "" {
			return nil
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: refKind,
			RecipientRefID:   evt.GetEntityId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Workflow update",
			Message:          "A workflow related to your request has been updated.",
			TemplateData: map[string]string{
				"workflow_instance_id": evt.GetWorkflowInstanceId(),
				"entity_type":          evt.GetEntityType(),
				"entity_id":            evt.GetEntityId(),
				"status":               evt.GetStatus(),
			},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
			},
		})
	case mustTopic("organisation_approved"):
		var evt b2beventsv1.OrganisationApprovedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode organisation approved event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "organisation_admin",
			RecipientRefID:   evt.GetOrganisationId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Organisation approved",
			Message:          "Your organisation has been approved and is ready for B2B operations.",
			TemplateData: map[string]string{
				"organisation_id": evt.GetOrganisationId(),
				"approved_by":     evt.GetApprovedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("organisation_suspended"):
		var evt b2beventsv1.OrganisationSuspendedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode organisation suspended event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "organisation_admin",
			RecipientRefID:   evt.GetOrganisationId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT,
			Subject:          "Organisation suspended",
			Message:          "Your organisation has been suspended.",
			TemplateData: map[string]string{
				"organisation_id": evt.GetOrganisationId(),
				"reason":          evt.GetReason(),
				"suspended_by":    evt.GetSuspendedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("org_member_added"):
		var evt b2beventsv1.OrgMemberAddedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode org member added event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:     "You were added to an organisation",
			Message:     "You have been added to an organisation in the B2B portal.",
			TemplateData: map[string]string{
				"organisation_id": evt.GetOrganisationId(),
				"member_id":       evt.GetMemberId(),
				"role":            evt.GetRole().String(),
				"added_by":        evt.GetAddedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("org_member_removed"):
		var evt b2beventsv1.OrgMemberRemovedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode org member removed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:     "Organisation membership removed",
			Message:     "Your organisation membership has been removed.",
			TemplateData: map[string]string{
				"organisation_id": evt.GetOrganisationId(),
				"member_id":       evt.GetMemberId(),
				"removed_by":      evt.GetRemovedBy(),
				"reason":          evt.GetReason(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("b2b_admin_assigned"):
		var evt b2beventsv1.B2BAdminAssignedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode b2b admin assigned event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: evt.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:     "B2B admin role assigned",
			Message:     "You have been assigned as a B2B admin.",
			TemplateData: map[string]string{
				"organisation_id": evt.GetOrganisationId(),
				"user_id":         evt.GetUserId(),
				"assigned_by":     evt.GetAssignedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("authz_events"):
		return c.handleAuthzEvents(ctx, descriptor, payload)
	case mustTopic("storage_events"):
		return c.handleStorageEvents(ctx, descriptor, payload)
	case mustTopic("partner_events"):
		return c.handlePartnerEvents(ctx, descriptor, payload)
	case mustTopic("media_events"):
		return c.handleMediaEvents(ctx, descriptor, payload)
	case mustTopic("docgen_events"):
		return c.handleDocgenEvents(ctx, descriptor, payload)
	case mustTopic("payment_completed"):
		var evt paymenteventsv1.PaymentCompletedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode payment completed event: %w", err)
		}
		amount := ""
		if evt.Amount != nil {
			amount = fmt.Sprintf("%s %0.2f", evt.Amount.GetCurrency(), float64(evt.Amount.GetAmount())/100)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetActorUserId(),
			RecipientRefKind: "payment",
			RecipientRefID:   evt.GetPaymentId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Payment confirmation",
			Message:          "Your payment has been completed successfully.",
			TemplateData:     map[string]string{"payment_id": evt.GetPaymentId(), "transaction_id": evt.GetTransactionId(), "amount": amount, "receipt_number": evt.GetReceiptNumber()},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("payment_failed"):
		var evt paymenteventsv1.PaymentFailedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode payment failed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      firstNonEmpty(evt.GetActorUserId(), evt.GetPayerId()),
			RecipientRefKind: "payment",
			RecipientRefID:   evt.GetPaymentId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Payment failed",
			Message:          "We could not complete your payment.",
			TemplateData:     map[string]string{"payment_id": evt.GetPaymentId(), "reason": evt.GetErrorMessage()},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("refund_processed"):
		var evt paymenteventsv1.RefundProcessedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode refund processed event: %w", err)
		}
		amount := ""
		if evt.Amount != nil {
			amount = fmt.Sprintf("%s %0.2f", evt.Amount.GetCurrency(), float64(evt.Amount.GetAmount())/100)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      firstNonEmpty(evt.GetRecipientId(), evt.GetActorUserId()),
			RecipientRefKind: "refund",
			RecipientRefID:   evt.GetRefundId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Refund processed",
			Message:          "Your refund has been processed.",
			TemplateData:     map[string]string{"refund_id": evt.GetRefundId(), "payment_id": evt.GetOriginalPaymentId(), "amount": amount},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
				Portal:         evt.GetPortal(),
				ActorUserID:    evt.GetActorUserId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("receipt_generated"):
		var evt paymenteventsv1.ReceiptGeneratedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode receipt generated event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "payment",
			RecipientRefID:   evt.GetPaymentId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Payment receipt available",
			Message:          "Your payment receipt is ready.",
			TemplateData:     map[string]string{"payment_id": evt.GetPaymentId(), "receipt_number": evt.GetReceiptNumber(), "receipt_file_id": evt.GetReceiptFileId()},
			Source: EventContext{
				Topic:          descriptor.Topic,
				TopicGroup:     descriptor.Group,
				Audience:       descriptor.Audience,
				CorrelationID:  evt.GetCorrelationId(),
				CausationID:    evt.GetCausationId(),
				TenantID:       evt.GetTenantId(),
				OrganisationID: evt.GetOrganisationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("policy_issued"):
		var evt policyeventsv1.PolicyIssuedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode policy issued event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "policy",
			RecipientRefID:   evt.GetPolicyId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_POLICY_ISSUED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Policy issued",
			Message:          "Your insurance policy has been issued successfully.",
			TemplateData:     map[string]string{"policy_id": evt.GetPolicyId(), "policy_number": evt.GetPolicyNumber()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("claim_submitted"):
		var evt claimeventsv1.ClaimSubmittedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode claim submitted event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "claim",
			RecipientRefID:   evt.GetClaimId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_CLAIM_SUBMITTED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Claim submitted",
			Message:          "We received your claim submission.",
			TemplateData:     map[string]string{"claim_id": evt.GetClaimId(), "claim_number": evt.GetClaimNumber(), "policy_id": evt.GetPolicyId()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
			},
		})
	case mustTopic("claim_approved"):
		var evt claimeventsv1.ClaimApprovedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode claim approved event: %w", err)
		}
		amount := ""
		if evt.ApprovedAmount != nil {
			amount = fmt.Sprintf("%s %0.2f", evt.ApprovedAmount.GetCurrency(), float64(evt.ApprovedAmount.GetAmount())/100)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "claim",
			RecipientRefID:   evt.GetClaimId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_CLAIM_APPROVED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Claim approved",
			Message:          "Your claim has been approved.",
			TemplateData:     map[string]string{"claim_id": evt.GetClaimId(), "claim_number": evt.GetClaimNumber(), "approved_amount": amount},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("claim_rejected"):
		var evt claimeventsv1.ClaimRejectedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode claim rejected event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "claim",
			RecipientRefID:   evt.GetClaimId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_CLAIM_REJECTED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Claim rejected",
			Message:          "Your claim could not be approved.",
			TemplateData:     map[string]string{"claim_id": evt.GetClaimId(), "claim_number": evt.GetClaimNumber(), "reason": evt.GetReason()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("claim_settled"):
		var evt claimeventsv1.ClaimSettledEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode claim settled event: %w", err)
		}
		amount := ""
		if evt.SettledAmount != nil {
			amount = fmt.Sprintf("%s %0.2f", evt.SettledAmount.GetCurrency(), float64(evt.SettledAmount.GetAmount())/100)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "claim",
			RecipientRefID:   evt.GetClaimId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_CLAIM_APPROVED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Claim settled",
			Message:          "Your claim settlement has been completed.",
			TemplateData:     map[string]string{"claim_id": evt.GetClaimId(), "claim_number": evt.GetClaimNumber(), "settled_amount": amount, "payment_reference": evt.GetPaymentReference()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("renewal_due"):
		var evt renewaleventsv1.RenewalDueEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode renewal due event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "renewal_schedule",
			RecipientRefID:   evt.GetRenewalScheduleId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_RENEWAL_REMINDER,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Renewal due soon",
			Message:          "Your policy renewal is due soon.",
			TemplateData:     map[string]string{"policy_id": evt.GetPolicyId(), "renewal_schedule_id": evt.GetRenewalScheduleId(), "renewal_due_date": evt.GetRenewalDueDate()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("grace_period_started"):
		var evt renewaleventsv1.GracePeriodStartedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode grace period started event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "grace_period",
			RecipientRefID:   evt.GetGracePeriodId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_GRACE_PERIOD,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT,
			Subject:          "Grace period started",
			Message:          "Your policy is now in grace period.",
			TemplateData:     map[string]string{"policy_id": evt.GetPolicyId(), "grace_period_id": evt.GetGracePeriodId(), "end_date": evt.GetEndDate()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("policy_lapsed"):
		var evt renewaleventsv1.PolicyLapsedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode renewal policy lapsed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "policy",
			RecipientRefID:   evt.GetPolicyId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_POLICY_LAPSED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT,
			Subject:          "Policy lapsed",
			Message:          "Your policy has lapsed.",
			TemplateData:     map[string]string{"policy_id": evt.GetPolicyId(), "reason": evt.GetReason(), "grace_period_id": evt.GetGracePeriodId()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	case mustTopic("policy_lapsed_bridge"):
		var evt policyeventsv1.PolicyLapsedEvent
		if err := protojson.Unmarshal(payload, &evt); err != nil {
			return fmt.Errorf("decode policy lapsed event: %w", err)
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      evt.GetCustomerId(),
			RecipientRefKind: "policy",
			RecipientRefID:   evt.GetPolicyId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_POLICY_LAPSED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT,
			Subject:          "Policy lapsed",
			Message:          "Your policy has lapsed.",
			TemplateData:     map[string]string{"policy_id": evt.GetPolicyId(), "policy_number": evt.GetPolicyNumber(), "reason": evt.GetReason()},
			Source: EventContext{
				Topic:         descriptor.Topic,
				TopicGroup:    descriptor.Group,
				Audience:      descriptor.Audience,
				CorrelationID: evt.GetCorrelationId(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	default:
		return nil
	}
}

func mustTopic(key string) string {
	descriptor, _ := DescriptorForKey(key)
	return descriptor.Topic
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func workflowEntityRecipientRef(entityType string) string {
	switch strings.ToLower(strings.TrimSpace(entityType)) {
	case "claim":
		return "claim"
	case "policy":
		return "policy"
	case "payment":
		return "payment"
	case "order":
		return "order"
	default:
		return ""
	}
}

func formatMoney(amount interface {
	GetCurrency() string
	GetAmount() int64
}) string {
	if amount == nil {
		return ""
	}
	return fmt.Sprintf("%s %0.2f", amount.GetCurrency(), float64(amount.GetAmount())/100)
}

func (c *Consumer) handleAuthzEvents(ctx context.Context, descriptor TopicDescriptor, payload []byte) error {
	var roleAssigned authzeventsv1.RoleAssignedEvent
	if err := protojson.Unmarshal(payload, &roleAssigned); err == nil && roleAssigned.GetUserId() != "" && roleAssigned.GetRoleId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: roleAssigned.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:     "Role assigned",
			Message:     "A new role has been assigned to your account.",
			TemplateData: map[string]string{
				"user_id":     roleAssigned.GetUserId(),
				"role_id":     roleAssigned.GetRoleId(),
				"role_name":   roleAssigned.GetRoleName(),
				"domain":      roleAssigned.GetDomain(),
				"assigned_by": roleAssigned.GetAssignedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var roleRemoved authzeventsv1.RoleRemovedEvent
	if err := protojson.Unmarshal(payload, &roleRemoved); err == nil && roleRemoved.GetUserId() != "" && roleRemoved.GetRoleId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: roleRemoved.GetUserId(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:     "Role removed",
			Message:     "A role has been removed from your account.",
			TemplateData: map[string]string{
				"user_id":    roleRemoved.GetUserId(),
				"role_id":    roleRemoved.GetRoleId(),
				"role_name":  roleRemoved.GetRoleName(),
				"domain":     roleRemoved.GetDomain(),
				"removed_by": roleRemoved.GetRemovedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var portalUpdated authzeventsv1.PortalConfigUpdatedEvent
	if err := protojson.Unmarshal(payload, &portalUpdated); err == nil && portalUpdated.GetUpdatedBy() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: portalUpdated.GetUpdatedBy(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:     "Portal configuration updated",
			Message:     "Portal access configuration has been updated.",
			TemplateData: map[string]string{
				"portal":         portalUpdated.GetPortal().String(),
				"updated_by":     portalUpdated.GetUpdatedBy(),
				"changed_fields": fmt.Sprintf("%d", len(portalUpdated.GetChangedFields())),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	return nil
}

func (c *Consumer) handleStorageEvents(ctx context.Context, descriptor TopicDescriptor, payload []byte) error {
	var uploaded storageeventsv1.FileUploadedEvent
	if err := protojson.Unmarshal(payload, &uploaded); err == nil && uploaded.GetFileId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      uploaded.GetUploadedBy(),
			RecipientRefKind: "file",
			RecipientRefID:   uploaded.GetFileId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "File uploaded",
			Message:          "Your file upload has been received successfully.",
			TemplateData: map[string]string{
				"file_id":        uploaded.GetFileId(),
				"filename":       uploaded.GetFilename(),
				"reference_id":   uploaded.GetReferenceId(),
				"reference_type": uploaded.GetReferenceType(),
				"content_type":   uploaded.GetContentType(),
			},
			Source: EventContext{
				Topic:       descriptor.Topic,
				TopicGroup:  descriptor.Group,
				Audience:    descriptor.Audience,
				TenantID:    uploaded.GetTenantId(),
				ActorUserID: uploaded.GetUploadedBy(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var finalized storageeventsv1.FileUploadFinalizedEvent
	if err := protojson.Unmarshal(payload, &finalized); err == nil && finalized.GetFileId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      finalized.GetFinalizedBy(),
			RecipientRefKind: "file",
			RecipientRefID:   finalized.GetFileId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "File ready",
			Message:          "Your uploaded file has been finalized and is ready to use.",
			TemplateData: map[string]string{
				"file_id":      finalized.GetFileId(),
				"content_type": finalized.GetContentType(),
			},
			Source: EventContext{
				Topic:       descriptor.Topic,
				TopicGroup:  descriptor.Group,
				Audience:    descriptor.Audience,
				TenantID:    finalized.GetTenantId(),
				ActorUserID: finalized.GetFinalizedBy(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var deleted storageeventsv1.FileDeletedEvent
	if err := protojson.Unmarshal(payload, &deleted); err == nil && deleted.GetFileId() != "" && deleted.GetDeletedBy() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID: deleted.GetDeletedBy(),
			Type:        notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:     "File deleted",
			Message:     "A stored file has been deleted.",
			TemplateData: map[string]string{
				"file_id":     deleted.GetFileId(),
				"storage_key": deleted.GetStorageKey(),
			},
			Source: EventContext{
				Topic:       descriptor.Topic,
				TopicGroup:  descriptor.Group,
				Audience:    descriptor.Audience,
				TenantID:    deleted.GetTenantId(),
				ActorUserID: deleted.GetDeletedBy(),
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	return nil
}

func (c *Consumer) handlePartnerEvents(ctx context.Context, descriptor TopicDescriptor, payload []byte) error {
	var onboarded partnereventsv1.PartnerOnboardedEvent
	if err := protojson.Unmarshal(payload, &onboarded); err == nil && onboarded.GetPartnerId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientID:      onboarded.GetFocalPersonId(),
			RecipientRefKind: "partner",
			RecipientRefID:   onboarded.GetPartnerId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Partner onboarding completed",
			Message:          "Your partner organisation has been onboarded.",
			TemplateData: map[string]string{
				"partner_id":        onboarded.GetPartnerId(),
				"organization_name": onboarded.GetOrganizationName(),
				"partner_type":      onboarded.GetPartnerType(),
				"focal_person_id":   onboarded.GetFocalPersonId(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var verified partnereventsv1.PartnerVerifiedEvent
	if err := protojson.Unmarshal(payload, &verified); err == nil && verified.GetPartnerId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "partner",
			RecipientRefID:   verified.GetPartnerId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Partner verified",
			Message:          "Your partner profile has been verified.",
			TemplateData: map[string]string{
				"partner_id":  verified.GetPartnerId(),
				"verified_by": verified.GetVerifiedBy(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var agentRegistered partnereventsv1.AgentRegisteredEvent
	if err := protojson.Unmarshal(payload, &agentRegistered); err == nil && agentRegistered.GetAgentId() != "" {
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "agent",
			RecipientRefID:   agentRegistered.GetAgentId(),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Agent registration completed",
			Message:          "An agent profile has been registered successfully.",
			TemplateData: map[string]string{
				"agent_id":   agentRegistered.GetAgentId(),
				"partner_id": agentRegistered.GetPartnerId(),
				"agent_name": agentRegistered.GetAgentName(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	var commission partnereventsv1.CommissionCalculatedEvent
	if err := protojson.Unmarshal(payload, &commission); err == nil && commission.GetCommissionId() != "" {
		refKind := "partner"
		refID := commission.GetPartnerId()
		if commission.GetAgentId() != "" {
			refKind = "agent"
			refID = commission.GetAgentId()
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: refKind,
			RecipientRefID:   refID,
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Commission calculated",
			Message:          "A commission has been calculated for your distribution activity.",
			TemplateData: map[string]string{
				"commission_id":     commission.GetCommissionId(),
				"partner_id":        commission.GetPartnerId(),
				"agent_id":          commission.GetAgentId(),
				"policy_id":         commission.GetPolicyId(),
				"commission_amount": formatMoney(commission.GetCommissionAmount()),
				"commission_type":   commission.GetCommissionType(),
			},
			Source: EventContext{
				Topic:      descriptor.Topic,
				TopicGroup: descriptor.Group,
				Audience:   descriptor.Audience,
			},
			PreferredChannels: []notificationv1.NotificationChannel{
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
			},
		})
	}

	return nil
}

func (c *Consumer) handleMediaEvents(ctx context.Context, descriptor TopicDescriptor, payload []byte) error {
	var generic mediaAggregateEvent
	if err := json.Unmarshal(payload, &generic); err != nil || generic.EventType == "" {
		var uploaded mediaeventsv1.MediaFileUploadedEvent
		if err := protojson.Unmarshal(payload, &uploaded); err == nil && uploaded.GetMediaId() != "" {
			return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
				RecipientID:      uploaded.GetUploadedBy(),
				RecipientRefKind: "media",
				RecipientRefID:   uploaded.GetMediaId(),
				Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
				Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
				Subject:          "Media uploaded",
				Message:          "Your media upload has been received.",
				TemplateData: map[string]string{
					"media_id":    uploaded.GetMediaId(),
					"file_id":     uploaded.GetFileId(),
					"entity_type": uploaded.GetEntityType(),
					"entity_id":   uploaded.GetEntityId(),
					"media_type":  uploaded.GetMediaType(),
					"mime_type":   uploaded.GetMimeType(),
				},
				Source: EventContext{
					Topic:         descriptor.Topic,
					TopicGroup:    descriptor.Group,
					Audience:      descriptor.Audience,
					TenantID:      uploaded.GetTenantId(),
					ActorUserID:   uploaded.GetUploadedBy(),
					CorrelationID: uploaded.GetCorrelationId(),
				},
				PreferredChannels: []notificationv1.NotificationChannel{
					notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				},
			})
		}
		return nil
	}

	switch generic.EventType {
	case "media.file.uploaded":
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "media",
			RecipientRefID:   generic.MediaID,
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Media uploaded",
			Message:          "Your media upload has been received.",
			TemplateData: map[string]string{
				"media_id":   generic.MediaID,
				"event_type": generic.EventType,
			},
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience, TenantID: generic.TenantID},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP},
		})
	case "media.processing.completed":
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind:  "media",
			RecipientRefID:    generic.MediaID,
			Type:              notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:          notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:           "Media processing completed",
			Message:           "Media processing has completed successfully.",
			TemplateData:      aggregateTemplateData(generic),
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience, TenantID: generic.TenantID},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP},
		})
	case "media.processing.failed":
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind:  "media",
			RecipientRefID:    generic.MediaID,
			Type:              notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:          notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:           "Media processing failed",
			Message:           "Media processing failed and may need attention.",
			TemplateData:      aggregateTemplateData(generic),
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience, TenantID: generic.TenantID},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL},
		})
	case "media.virus.detected":
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind:  "media",
			RecipientRefID:    generic.MediaID,
			Type:              notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:          notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT,
			Subject:           "Virus detected in media upload",
			Message:           "A virus was detected in an uploaded media file.",
			TemplateData:      aggregateTemplateData(generic),
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience, TenantID: generic.TenantID},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL},
		})
	default:
		return nil
	}
}

func (c *Consumer) handleDocgenEvents(ctx context.Context, descriptor TopicDescriptor, payload []byte) error {
	var generic docgenAggregateEvent
	if err := json.Unmarshal(payload, &generic); err != nil || generic.EventType == "" {
		return nil
	}

	switch generic.EventType {
	case "document.generated":
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: firstNonEmpty(workflowEntityRecipientRef(generic.EntityType), "document_generation"),
			RecipientRefID:   firstNonEmpty(generic.EntityID, generic.DocumentGenerationID),
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Document ready",
			Message:          "A generated document is now available.",
			TemplateData: map[string]string{
				"document_generation_id": generic.DocumentGenerationID,
				"entity_type":            generic.EntityType,
				"entity_id":              generic.EntityID,
				"file_url":               generic.FileURL,
				"correlation_id":         generic.CorrelationID,
			},
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP},
		})
	case "document.generation.failed":
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: "document_generation",
			RecipientRefID:   generic.DocumentGenerationID,
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
			Subject:          "Document generation failed",
			Message:          "A document generation request could not be completed.",
			TemplateData: map[string]string{
				"document_generation_id": generic.DocumentGenerationID,
				"error_message":          generic.ErrorMessage,
				"correlation_id":         generic.CorrelationID,
			},
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL},
		})
	case "document.generation.requested":
		refKind := workflowEntityRecipientRef(generic.EntityType)
		if refKind == "" {
			return nil
		}
		return c.service.QueueDomainNotification(ctx, &DomainNotificationInput{
			RecipientRefKind: refKind,
			RecipientRefID:   generic.EntityID,
			Type:             notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED,
			Priority:         notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Subject:          "Document generation started",
			Message:          "A document generation request has started for your record.",
			TemplateData: map[string]string{
				"document_generation_id": generic.DocumentGenerationID,
				"document_template_id":   generic.DocumentTemplateID,
				"entity_type":            generic.EntityType,
				"entity_id":              generic.EntityID,
				"correlation_id":         generic.CorrelationID,
			},
			Source:            EventContext{Topic: descriptor.Topic, TopicGroup: descriptor.Group, Audience: descriptor.Audience},
			PreferredChannels: []notificationv1.NotificationChannel{notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP},
		})
	default:
		return nil
	}
}

type mediaAggregateEvent struct {
	EventType string         `json:"event_type"`
	TenantID  string         `json:"tenant_id"`
	MediaID   string         `json:"media_id"`
	JobID     string         `json:"job_id"`
	Data      map[string]any `json:"data"`
}

type docgenAggregateEvent struct {
	EventType            string `json:"event_type"`
	TenantID             string `json:"tenant_id"`
	TemplateID           string `json:"template_id"`
	DocumentGenerationID string `json:"document_generation_id"`
	DocumentTemplateID   string `json:"document_template_id"`
	EntityType           string `json:"entity_type"`
	EntityID             string `json:"entity_id"`
	FileURL              string `json:"file_url"`
	ErrorMessage         string `json:"error_message"`
	CorrelationID        string `json:"correlation_id"`
}

func aggregateTemplateData(evt mediaAggregateEvent) map[string]string {
	values := map[string]string{
		"event_type": evt.EventType,
		"tenant_id":  evt.TenantID,
		"media_id":   evt.MediaID,
		"job_id":     evt.JobID,
	}
	for key, value := range evt.Data {
		values[key] = fmt.Sprintf("%v", value)
	}
	return values
}
