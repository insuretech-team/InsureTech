package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/delivery"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/internalrpc"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/notifyprefs"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	notificationservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type emailSender interface {
	Send(ctx context.Context, to, subject, body string) (*delivery.EmailResponse, error)
}

type smsSender interface {
	Send(ctx context.Context, req *delivery.SMSRequest) (*delivery.SMSResponse, error)
}

type pushSender interface {
	Send(ctx context.Context, req *delivery.PushRequest) (*delivery.PushResponse, error)
}

type webhookSender interface {
	Send(ctx context.Context, req *delivery.WebhookRequest) (*delivery.WebhookResponse, error)
}

type statePublisher interface {
	PublishNotificationSent(ctx context.Context, notificationID, recipientID, channel, notificationType, correlationID string) error
	PublishNotificationDelivered(ctx context.Context, notificationID, recipientID, correlationID string) error
	PublishNotificationFailed(ctx context.Context, notificationID, recipientID, errorMessage, correlationID string) error
}

type authnPreferenceClient interface {
	UpdateNotificationPreferences(ctx context.Context, req *authnservicev1.UpdateNotificationPreferencesRequest, opts ...grpc.CallOption) (*authnservicev1.UpdateNotificationPreferencesResponse, error)
}

type Service struct {
	notificationRepo *repository.NotificationRepository
	templateRepo     *repository.TemplateRepository
	userRepo         *repository.UserRepository
	lookupRepo       *repository.LookupRepository
	pushTokenRepo    *repository.PushTokenRepository
	webhookRepo      *repository.WebhookRepository
	emailClient      emailSender
	smsClient        smsSender
	pushClient       pushSender
	webhookClient    webhookSender
	publisher        statePublisher
	authnClient      authnPreferenceClient
	cfg              *config.Config
	now              func() time.Time

	dispatchWG sync.WaitGroup
}

var _ domain.NotificationService = (*Service)(nil)
var _ events.NotificationEventService = (*Service)(nil)

func NewService(
	notificationRepo *repository.NotificationRepository,
	templateRepo *repository.TemplateRepository,
	userRepo *repository.UserRepository,
	lookupRepo *repository.LookupRepository,
	pushTokenRepo *repository.PushTokenRepository,
	webhookRepo *repository.WebhookRepository,
	emailClient emailSender,
	smsClient smsSender,
	pushClient pushSender,
	webhookClient webhookSender,
	publisher statePublisher,
	cfg *config.Config,
) *Service {
	if cfg == nil {
		cfg = &config.Config{}
	}
	return &Service{
		notificationRepo: notificationRepo,
		templateRepo:     templateRepo,
		userRepo:         userRepo,
		lookupRepo:       lookupRepo,
		pushTokenRepo:    pushTokenRepo,
		webhookRepo:      webhookRepo,
		emailClient:      emailClient,
		smsClient:        smsClient,
		pushClient:       pushClient,
		webhookClient:    webhookClient,
		publisher:        publisher,
		cfg:              cfg,
		now:              time.Now,
	}
}

func (s *Service) WithAuthNPreferenceClient(client authnPreferenceClient) *Service {
	s.authnClient = client
	return s
}

func (s *Service) SendNotification(ctx context.Context, req *notificationservicev1.SendNotificationRequest) (*notificationservicev1.SendNotificationResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if strings.TrimSpace(req.GetRecipientId()) == "" {
		return nil, errors.New("recipient_id is required")
	}

	preferences, err := s.userRepo.GetPreferences(ctx, req.GetRecipientId())
	if err != nil {
		return nil, fmt.Errorf("load recipient preferences: %w", err)
	}

	channel, err := s.resolveRequestedChannel(preferences, req.GetChannel())
	if err != nil {
		return nil, err
	}

	subject, body, err := s.resolveContent(ctx, req.GetTemplateId(), req.GetSubject(), req.GetMessage(), req.GetTemplateData(), channel, req.GetType())
	if err != nil {
		return nil, err
	}

	notification := s.buildNotificationRecord(req.GetRecipientId(), channel, req.GetType(), normalizePriority(req.GetPriority()), subject, body, cloneMap(req.GetTemplateData()), req.GetScheduleAfterSeconds())
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	return &notificationservicev1.SendNotificationResponse{
		NotificationId: notification.GetNotificationId(),
		Message:        "Notification queued successfully",
	}, nil
}

func (s *Service) SendBulkNotifications(ctx context.Context, req *notificationservicev1.SendBulkNotificationsRequest) (*notificationservicev1.SendBulkNotificationsResponse, error) {
	if req == nil || len(req.GetNotifications()) == 0 {
		return nil, errors.New("notifications are required")
	}

	ids := make([]string, 0, len(req.GetNotifications()))
	var successCount int32
	var failedCount int32

	for _, item := range req.GetNotifications() {
		resp, err := s.SendNotification(ctx, item)
		if err != nil {
			failedCount++
			appLogger.Warnf("bulk notification item failed: %v", err)
			continue
		}
		successCount++
		ids = append(ids, resp.GetNotificationId())
	}

	return &notificationservicev1.SendBulkNotificationsResponse{
		NotificationIds: ids,
		SuccessCount:    successCount,
		FailedCount:     failedCount,
	}, nil
}

func (s *Service) GetNotificationStatus(ctx context.Context, req *notificationservicev1.GetNotificationStatusRequest) (*notificationservicev1.GetNotificationStatusResponse, error) {
	if req == nil || strings.TrimSpace(req.GetNotificationId()) == "" {
		return nil, errors.New("notification_id is required")
	}
	notification, err := s.notificationRepo.GetByID(ctx, req.GetNotificationId())
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return &notificationservicev1.GetNotificationStatusResponse{Notification: notification}, nil
}

func (s *Service) GetUserNotifications(ctx context.Context, req *notificationservicev1.GetUserNotificationsRequest) (*notificationservicev1.GetUserNotificationsResponse, error) {
	if req == nil || strings.TrimSpace(req.GetUserId()) == "" {
		return nil, errors.New("user_id is required")
	}
	notifications, totalCount, unreadCount, err := s.notificationRepo.ListByRecipient(ctx, req.GetUserId(), req.GetUnreadOnly(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("list user notifications: %w", err)
	}
	return &notificationservicev1.GetUserNotificationsResponse{
		Notifications: notifications,
		TotalCount:    totalCount,
		UnreadCount:   unreadCount,
	}, nil
}

func (s *Service) MarkAsRead(ctx context.Context, req *notificationservicev1.MarkAsReadRequest) (*notificationservicev1.MarkAsReadResponse, error) {
	if req == nil || len(req.GetNotificationIds()) == 0 {
		return nil, errors.New("notification_ids are required")
	}
	if err := s.notificationRepo.MarkAsRead(ctx, req.GetNotificationIds()); err != nil {
		return nil, fmt.Errorf("mark notifications as read: %w", err)
	}
	return &notificationservicev1.MarkAsReadResponse{Message: "Notifications marked as read"}, nil
}

func (s *Service) UpdatePreferences(ctx context.Context, req *notificationservicev1.UpdatePreferencesRequest) (*notificationservicev1.UpdatePreferencesResponse, error) {
	if req == nil || strings.TrimSpace(req.GetUserId()) == "" {
		return nil, errors.New("user_id is required")
	}
	if req.GetPreferences() == nil {
		return nil, errors.New("preferences are required")
	}
	if s.authnClient == nil {
		return nil, errors.New("authn preference client is not configured")
	}

	outgoingCtx := internalrpc.OutgoingContext(ctx, "notification-service")
	if _, err := s.authnClient.UpdateNotificationPreferences(outgoingCtx, &authnservicev1.UpdateNotificationPreferencesRequest{
		UserId:                 req.GetUserId(),
		NotificationPreference: notifyprefs.Compact(req.GetPreferences()),
		PreferredLanguage:      s.preferredLanguageForUser(ctx, req.GetUserId()),
	}); err != nil {
		return nil, fmt.Errorf("update preferences via authn: %w", err)
	}
	return &notificationservicev1.UpdatePreferencesResponse{Message: "Notification preferences updated"}, nil
}

func (s *Service) CreateNotificationTemplate(ctx context.Context, req *notificationservicev1.CreateNotificationTemplateRequest) (*notificationservicev1.CreateNotificationTemplateResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if strings.TrimSpace(req.GetName()) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(req.GetBody()) == "" {
		return nil, errors.New("body is required")
	}

	templateRecord := &notificationv1.NotificationTemplate{
		TemplateId:      uuid.NewString(),
		TemplateName:    req.GetName(),
		Type:            req.GetType(),
		Channel:         notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
		SubjectTemplate: req.GetSubject(),
		BodyTemplate:    req.GetBody(),
		Language:        "en",
		IsActive:        true,
		CreatedAt:       timestamppb.New(s.now().UTC()),
		UpdatedAt:       timestamppb.New(s.now().UTC()),
	}
	if templateRecord.GetType() == notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED {
		templateRecord.Type = notificationv1.NotificationType_NOTIFICATION_TYPE_MARKETING
	}

	if err := s.templateRepo.Create(ctx, templateRecord); err != nil {
		return nil, fmt.Errorf("create notification template: %w", err)
	}

	return &notificationservicev1.CreateNotificationTemplateResponse{
		TemplateId: templateRecord.GetTemplateId(),
		Message:    "Notification template created",
	}, nil
}

func (s *Service) UpdateNotificationTemplate(ctx context.Context, req *notificationservicev1.UpdateNotificationTemplateRequest) (*notificationservicev1.UpdateNotificationTemplateResponse, error) {
	if req == nil || strings.TrimSpace(req.GetTemplateId()) == "" {
		return nil, errors.New("template_id is required")
	}
	if strings.TrimSpace(req.GetName()) == "" && req.GetSubject() == "" && req.GetBody() == "" {
		return nil, errors.New("at least one updatable field is required")
	}
	if err := s.templateRepo.Update(ctx, req.GetTemplateId(), req.GetName(), req.GetSubject(), req.GetBody()); err != nil {
		return nil, fmt.Errorf("update notification template: %w", err)
	}
	return &notificationservicev1.UpdateNotificationTemplateResponse{Message: "Notification template updated"}, nil
}

func (s *Service) DeactivateNotificationTemplate(ctx context.Context, req *notificationservicev1.DeactivateNotificationTemplateRequest) (*notificationservicev1.DeactivateNotificationTemplateResponse, error) {
	if req == nil || strings.TrimSpace(req.GetTemplateId()) == "" {
		return nil, errors.New("template_id is required")
	}
	if err := s.templateRepo.Deactivate(ctx, req.GetTemplateId()); err != nil {
		return nil, fmt.Errorf("deactivate notification template: %w", err)
	}
	return &notificationservicev1.DeactivateNotificationTemplateResponse{Message: "Notification template deactivated"}, nil
}

func (s *Service) QueueDomainNotification(ctx context.Context, input *events.DomainNotificationInput) error {
	if input == nil {
		return errors.New("notification input is required")
	}

	recipientIDs, err := s.resolveRecipientIDs(ctx, input)
	if err != nil {
		return err
	}
	if len(recipientIDs) == 0 {
		return nil
	}

	var queued int
	for _, recipientID := range recipientIDs {
		preferences, err := s.userRepo.GetPreferences(ctx, recipientID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				appLogger.Warnf("notification recipient %s not found for source topic %s", recipientID, input.Source.Topic)
				continue
			}
			return fmt.Errorf("load notification preferences for %s: %w", recipientID, err)
		}

		channels := s.channelsForDomainInput(preferences, input)
		if len(channels) == 0 {
			appLogger.Infof("notification skipped for recipient %s due to preference filtering", recipientID)
			continue
		}

		templateData := s.enrichTemplateData(input)
		for _, channel := range channels {
			notification := s.buildNotificationRecord(recipientID, channel, normalizeType(input.Type), normalizePriority(input.Priority), input.Subject, input.Message, templateData, 0)
			if err := s.notificationRepo.Create(ctx, notification); err != nil {
				return fmt.Errorf("create domain notification for %s: %w", recipientID, err)
			}
			queued++
		}
	}

	if queued == 0 {
		appLogger.Infof("no notification records queued for topic %s", input.Source.Topic)
	}
	return nil
}

func (s *Service) StartDispatcher(ctx context.Context) {
	s.dispatchWG.Add(1)
	go func() {
		defer s.dispatchWG.Done()

		interval := 15 * time.Second
		if s.cfg != nil && s.cfg.Delivery.DispatchInterval > 0 {
			interval = s.cfg.Delivery.DispatchInterval
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			if err := s.RunDispatchCycle(ctx); err != nil && ctx.Err() == nil {
				appLogger.Errorf("notification dispatch cycle failed: %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

func (s *Service) WaitForDispatcher() {
	s.dispatchWG.Wait()
}

func (s *Service) RunDispatchCycle(ctx context.Context) error {
	batchSize := 50
	if s.cfg != nil && s.cfg.Delivery.BatchSize > 0 {
		batchSize = s.cfg.Delivery.BatchSize
	}

	dueNotifications, err := s.notificationRepo.ListDue(ctx, s.now().UTC(), batchSize)
	if err != nil {
		return fmt.Errorf("list due notifications: %w", err)
	}

	var firstErr error
	for _, notification := range dueNotifications {
		if err := s.dispatchNotification(ctx, notification); err != nil && ctx.Err() == nil {
			appLogger.Errorf("dispatch notification %s failed: %v", notification.GetNotificationId(), err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	if err := s.runWebhookDispatchCycle(ctx); err != nil && ctx.Err() == nil {
		appLogger.Errorf("webhook dispatch cycle failed: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *Service) buildNotificationRecord(
	recipientID string,
	channel notificationv1.NotificationChannel,
	notificationType notificationv1.NotificationType,
	priority notificationv1.NotificationPriority,
	subject string,
	message string,
	templateData map[string]string,
	delaySeconds int64,
) *notificationv1.Notification {
	now := s.now().UTC()
	notification := &notificationv1.Notification{
		NotificationId: uuid.NewString(),
		RecipientId:    recipientID,
		Type:           notificationType,
		Channel:        channel,
		Subject:        subject,
		Message:        message,
		TemplateData:   cloneMap(templateData),
		Priority:       priority,
		Status:         notificationv1.NotificationStatus_NOTIFICATION_STATUS_QUEUED,
		CreatedAt:      timestamppb.New(now),
	}
	if delaySeconds > 0 {
		notification.ScheduledAt = timestamppb.New(now.Add(time.Duration(delaySeconds) * time.Second))
	}
	return notification
}

func (s *Service) resolveContent(
	ctx context.Context,
	templateID string,
	subject string,
	message string,
	templateData map[string]string,
	channel notificationv1.NotificationChannel,
	notificationType notificationv1.NotificationType,
) (string, string, error) {
	if strings.TrimSpace(templateID) == "" {
		if strings.TrimSpace(message) == "" {
			return "", "", errors.New("message is required when template_id is not provided")
		}
		return subject, message, nil
	}

	templateRecord, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return "", "", fmt.Errorf("get notification template: %w", err)
	}
	if !templateRecord.GetIsActive() {
		return "", "", errors.New("notification template is not active")
	}
	if templateRecord.GetChannel() != notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED &&
		channel != notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED &&
		templateRecord.GetChannel() != channel {
		appLogger.Warnf("notification template %s channel %s reused for %s", templateID, templateRecord.GetChannel().String(), channel.String())
	}
	if templateRecord.GetType() != notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED &&
		notificationType != notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED &&
		templateRecord.GetType() != notificationType {
		appLogger.Warnf("notification template %s type %s reused for %s", templateID, templateRecord.GetType().String(), notificationType.String())
	}

	renderedSubject, err := renderTextTemplate(templateRecord.GetSubjectTemplate(), templateData)
	if err != nil {
		return "", "", fmt.Errorf("render template subject: %w", err)
	}
	renderedBody, err := renderTextTemplate(templateRecord.GetBodyTemplate(), templateData)
	if err != nil {
		return "", "", fmt.Errorf("render template body: %w", err)
	}
	return renderedSubject, renderedBody, nil
}

func (s *Service) resolveRequestedChannel(preferences *notificationv1.NotificationPreference, requested notificationv1.NotificationChannel) (notificationv1.NotificationChannel, error) {
	if requested != notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED {
		if !s.channelEnabled(preferences, requested) {
			return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED, fmt.Errorf("channel %s is disabled by user preferences", requested.String())
		}
		return requested, nil
	}

	for _, channel := range defaultChannelOrder() {
		if s.channelEnabled(preferences, channel) {
			return channel, nil
		}
	}
	return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED, errors.New("no enabled notification channels found for recipient")
}

func (s *Service) channelsForDomainInput(preferences *notificationv1.NotificationPreference, input *events.DomainNotificationInput) []notificationv1.NotificationChannel {
	if preferences == nil {
		return nil
	}

	if input.Type == notificationv1.NotificationType_NOTIFICATION_TYPE_MARKETING && !preferences.GetMarketingOptIn() {
		return nil
	}
	if input.Type != notificationv1.NotificationType_NOTIFICATION_TYPE_MARKETING && !preferences.GetTransactionalOptIn() {
		return nil
	}

	requested := input.PreferredChannels
	if len(requested) == 0 {
		requested = defaultChannelOrder()
	}

	result := make([]notificationv1.NotificationChannel, 0, len(requested))
	for _, channel := range requested {
		if channel == notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED {
			continue
		}
		if !s.channelEnabled(preferences, channel) {
			continue
		}
		if slices.Contains(result, channel) {
			continue
		}
		result = append(result, channel)
	}
	return result
}

func (s *Service) channelEnabled(preferences *notificationv1.NotificationPreference, channel notificationv1.NotificationChannel) bool {
	if preferences == nil {
		return false
	}
	for _, pref := range preferences.GetChannelPreferences() {
		if pref != nil && pref.GetChannel() == channel {
			return pref.GetEnabled()
		}
	}
	return false
}

func (s *Service) resolveRecipientIDs(ctx context.Context, input *events.DomainNotificationInput) ([]string, error) {
	if strings.TrimSpace(input.RecipientID) != "" {
		return []string{strings.TrimSpace(input.RecipientID)}, nil
	}

	switch strings.ToLower(strings.TrimSpace(input.RecipientRefKind)) {
	case "payment":
		payerID, err := s.lookupRepo.GetPaymentPayer(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve payment recipient: %w", err)
		}
		return []string{payerID}, nil
	case "order":
		order, err := s.lookupRepo.GetOrder(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve order recipient: %w", err)
		}
		return []string{order.CustomerID}, nil
	case "refund":
		if paymentID := strings.TrimSpace(input.TemplateData["payment_id"]); paymentID != "" {
			payerID, err := s.lookupRepo.GetPaymentPayer(ctx, paymentID)
			if err != nil {
				return nil, fmt.Errorf("resolve refund recipient from payment: %w", err)
			}
			return []string{payerID}, nil
		}
	case "policy":
		policy, err := s.lookupRepo.GetPolicy(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve policy recipient: %w", err)
		}
		return []string{policy.CustomerID}, nil
	case "claim":
		claim, err := s.lookupRepo.GetClaim(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve claim recipient: %w", err)
		}
		return []string{claim.CustomerID}, nil
	case "renewal_schedule":
		schedule, err := s.lookupRepo.GetRenewalSchedule(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve renewal recipient: %w", err)
		}
		policy, err := s.lookupRepo.GetPolicy(ctx, schedule.PolicyID)
		if err != nil {
			return nil, fmt.Errorf("resolve renewal policy recipient: %w", err)
		}
		return []string{policy.CustomerID}, nil
	case "grace_period":
		gracePeriod, err := s.lookupRepo.GetGracePeriod(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve grace period recipient: %w", err)
		}
		policy, err := s.lookupRepo.GetPolicy(ctx, gracePeriod.PolicyID)
		if err != nil {
			return nil, fmt.Errorf("resolve grace period policy recipient: %w", err)
		}
		return []string{policy.CustomerID}, nil
	case "organisation_admin":
		userIDs, err := s.lookupRepo.ListOrganisationAdminUserIDs(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve organisation admin recipients: %w", err)
		}
		return userIDs, nil
	case "partner":
		partner, err := s.lookupRepo.GetPartner(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve partner recipient: %w", err)
		}
		if strings.TrimSpace(partner.FocalPersonID) != "" {
			return []string{partner.FocalPersonID}, nil
		}
	case "agent":
		agent, err := s.lookupRepo.GetAgent(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve agent recipient: %w", err)
		}
		if strings.TrimSpace(agent.UserID) != "" {
			return []string{agent.UserID}, nil
		}
		if strings.TrimSpace(agent.PartnerID) != "" {
			partner, err := s.lookupRepo.GetPartner(ctx, agent.PartnerID)
			if err != nil {
				return nil, fmt.Errorf("resolve agent partner recipient: %w", err)
			}
			if strings.TrimSpace(partner.FocalPersonID) != "" {
				return []string{partner.FocalPersonID}, nil
			}
		}
	case "file":
		file, err := s.lookupRepo.GetFile(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve file recipient: %w", err)
		}
		if strings.TrimSpace(file.UploadedBy) != "" {
			return []string{file.UploadedBy}, nil
		}
		if refKind := normalizedEntityRefKind(file.ReferenceType); refKind != "" && strings.TrimSpace(file.ReferenceID) != "" {
			return s.resolveRecipientIDs(ctx, &events.DomainNotificationInput{
				RecipientRefKind: refKind,
				RecipientRefID:   file.ReferenceID,
			})
		}
	case "media":
		media, err := s.lookupRepo.GetMedia(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve media recipient: %w", err)
		}
		if strings.TrimSpace(media.UploadedBy) != "" {
			return []string{media.UploadedBy}, nil
		}
		if refKind := normalizedEntityRefKind(media.EntityType); refKind != "" && strings.TrimSpace(media.EntityID) != "" {
			return s.resolveRecipientIDs(ctx, &events.DomainNotificationInput{
				RecipientRefKind: refKind,
				RecipientRefID:   media.EntityID,
			})
		}
	case "document_generation":
		generation, err := s.lookupRepo.GetDocumentGeneration(ctx, input.RecipientRefID)
		if err != nil {
			return nil, fmt.Errorf("resolve document generation recipient: %w", err)
		}
		if refKind := normalizedEntityRefKind(generation.EntityType); refKind != "" && strings.TrimSpace(generation.EntityID) != "" {
			return s.resolveRecipientIDs(ctx, &events.DomainNotificationInput{
				RecipientRefKind: refKind,
				RecipientRefID:   generation.EntityID,
			})
		}
		if strings.TrimSpace(generation.GeneratedBy) != "" {
			return []string{generation.GeneratedBy}, nil
		}
	}

	return nil, fmt.Errorf("unable to resolve recipient for ref kind %q", input.RecipientRefKind)
}

func (s *Service) enrichTemplateData(input *events.DomainNotificationInput) map[string]string {
	data := cloneMap(input.TemplateData)
	if data == nil {
		data = map[string]string{}
	}
	addIfPresent(data, "source_topic", input.Source.Topic)
	addIfPresent(data, "source_group", string(input.Source.TopicGroup))
	addIfPresent(data, "audience", string(input.Source.Audience))
	addIfPresent(data, "correlation_id", input.Source.CorrelationID)
	addIfPresent(data, "causation_id", input.Source.CausationID)
	addIfPresent(data, "tenant_id", input.Source.TenantID)
	addIfPresent(data, "organisation_id", input.Source.OrganisationID)
	addIfPresent(data, "portal", input.Source.Portal)
	addIfPresent(data, "actor_user_id", input.Source.ActorUserID)
	addIfPresent(data, "trace_id", input.Source.TraceID)
	addIfPresent(data, "recipient_ref_kind", input.RecipientRefKind)
	addIfPresent(data, "recipient_ref_id", input.RecipientRefID)
	return data
}

func renderTextTemplate(text string, values map[string]string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", nil
	}
	tmpl, err := template.New("notification").Option("missingkey=zero").Parse(text)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	if err := tmpl.Execute(&builder, values); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func defaultChannelOrder() []notificationv1.NotificationChannel {
	return []notificationv1.NotificationChannel{
		notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
		notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
		notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
		notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH,
		notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_WHATSAPP,
	}
}

func normalizePriority(priority notificationv1.NotificationPriority) notificationv1.NotificationPriority {
	if priority == notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_UNSPECIFIED {
		return notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL
	}
	return priority
}

func normalizeType(notificationType notificationv1.NotificationType) notificationv1.NotificationType {
	return notificationType
}

func cloneMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func addIfPresent(values map[string]string, key, value string) {
	if strings.TrimSpace(value) != "" {
		values[key] = value
	}
}

func normalizedEntityRefKind(entityType string) string {
	switch strings.ToLower(strings.TrimSpace(entityType)) {
	case "order":
		return "order"
	case "policy":
		return "policy"
	case "claim":
		return "claim"
	case "payment":
		return "payment"
	case "partner":
		return "partner"
	case "agent":
		return "agent"
	case "file":
		return "file"
	case "media":
		return "media"
	default:
		return ""
	}
}

func (s *Service) preferredLanguageForUser(ctx context.Context, userID string) string {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(user.GetPreferredLanguage())
}
