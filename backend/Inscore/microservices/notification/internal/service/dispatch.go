package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/delivery"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/repository"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"gorm.io/gorm"
)

const (
	webhookEventSent      = "NOTIFICATION.SENT"
	webhookEventDelivered = "NOTIFICATION.DELIVERED"
	webhookEventFailed    = "NOTIFICATION.FAILED"
)

func (s *Service) dispatchNotification(ctx context.Context, notification *notificationv1.Notification) error {
	if notification == nil {
		return nil
	}

	if err := s.notificationRepo.MarkSending(ctx, notification.GetNotificationId()); err != nil {
		return fmt.Errorf("mark notification sending: %w", err)
	}

	now := s.now().UTC()
	if err := s.deliver(ctx, notification, now); err != nil {
		return s.handleDispatchFailure(ctx, notification, err)
	}
	return nil
}

func (s *Service) deliver(ctx context.Context, notification *notificationv1.Notification, now time.Time) error {
	switch notification.GetChannel() {
	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP:
		if err := s.notificationRepo.MarkDelivered(ctx, notification.GetNotificationId(), now); err != nil {
			return fmt.Errorf("mark in-app notification delivered: %w", err)
		}
		s.publishSentAndDelivered(ctx, notification)
		return nil

	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL:
		user, err := s.userRepo.GetByID(ctx, notification.GetRecipientId())
		if err != nil {
			return fmt.Errorf("load email recipient: %w", err)
		}
		if strings.TrimSpace(user.GetEmail()) == "" {
			return delivery.Permanent(errors.New("email recipient is not configured"))
		}
		if _, err := s.emailClient.Send(ctx, user.GetEmail(), notification.GetSubject(), notification.GetMessage()); err != nil {
			return err
		}
		if err := s.notificationRepo.MarkSent(ctx, notification.GetNotificationId(), now); err != nil {
			return fmt.Errorf("mark email notification sent: %w", err)
		}
		if err := s.notificationRepo.MarkDelivered(ctx, notification.GetNotificationId(), now); err != nil {
			return fmt.Errorf("mark email notification delivered: %w", err)
		}
		s.publishSentAndDelivered(ctx, notification)
		return nil

	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS:
		user, err := s.userRepo.GetByID(ctx, notification.GetRecipientId())
		if err != nil {
			return fmt.Errorf("load SMS recipient: %w", err)
		}
		if strings.TrimSpace(user.GetMobileNumber()) == "" {
			return delivery.Permanent(errors.New("sms recipient is not configured"))
		}
		resp, err := s.smsClient.Send(ctx, &delivery.SMSRequest{
			MSISDN:     user.GetMobileNumber(),
			Message:    notification.GetMessage(),
			UseMasking: true,
			CSMSID:     notification.GetNotificationId(),
		})
		if err != nil {
			return err
		}
		if err := s.notificationRepo.MarkSent(ctx, notification.GetNotificationId(), now); err != nil {
			return fmt.Errorf("mark sms notification sent: %w", err)
		}
		if strings.EqualFold(resp.Status, "DELIVERED") {
			if err := s.notificationRepo.MarkDelivered(ctx, notification.GetNotificationId(), now); err != nil {
				return fmt.Errorf("mark sms notification delivered: %w", err)
			}
			s.publishSentAndDelivered(ctx, notification)
			return nil
		}
		s.publishSent(ctx, notification)
		return nil

	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH:
		if s.pushTokenRepo == nil {
			return delivery.Permanent(delivery.ErrPushNotConfigured)
		}
		pushTokens, err := s.pushTokenRepo.ListActiveByUserID(ctx, notification.GetRecipientId())
		if err != nil {
			return fmt.Errorf("load push tokens: %w", err)
		}
		targets := make([]delivery.PushTarget, 0, len(pushTokens))
		for _, token := range pushTokens {
			if token == nil {
				continue
			}
			targets = append(targets, delivery.PushTarget{
				Provider:    token.Provider,
				Platform:    token.Platform,
				DeviceToken: token.DeviceToken,
				DeviceID:    token.DeviceID,
				AppID:       token.AppID,
			})
		}
		resp, err := s.pushClient.Send(ctx, &delivery.PushRequest{
			RecipientID: notification.GetRecipientId(),
			Title:       notification.GetSubject(),
			Body:        notification.GetMessage(),
			Data:        cloneMap(notification.GetTemplateData()),
			Targets:     targets,
		})
		if resp != nil && len(resp.InvalidTokens) > 0 {
			if deactivateErr := s.pushTokenRepo.DeactivateByDeviceTokens(ctx, resp.InvalidTokens); deactivateErr != nil {
				appLogger.Warnf("failed to deactivate invalid push tokens for notification %s: %v", notification.GetNotificationId(), deactivateErr)
			}
		}
		if err != nil {
			return err
		}
		if err := s.notificationRepo.MarkSent(ctx, notification.GetNotificationId(), now); err != nil {
			return fmt.Errorf("mark push notification sent: %w", err)
		}
		if resp != nil && strings.EqualFold(resp.Status, "DELIVERED") {
			if err := s.notificationRepo.MarkDelivered(ctx, notification.GetNotificationId(), now); err != nil {
				return fmt.Errorf("mark push notification delivered: %w", err)
			}
			s.publishSentAndDelivered(ctx, notification)
			return nil
		}
		s.publishSent(ctx, notification)
		return nil

	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_WHATSAPP:
		return delivery.Permanent(errors.New("whatsapp delivery is not configured"))

	default:
		return delivery.Permanent(fmt.Errorf("unsupported notification channel: %s", notification.GetChannel().String()))
	}
}

func (s *Service) handleDispatchFailure(ctx context.Context, notification *notificationv1.Notification, deliveryErr error) error {
	retryCount := notification.GetRetryCount() + 1
	if s.shouldRetry(retryCount, deliveryErr) {
		nextAttempt := s.now().UTC().Add(s.retryBackoff(retryCount))
		if err := s.notificationRepo.ScheduleRetry(ctx, notification.GetNotificationId(), retryCount, nextAttempt, deliveryErr.Error()); err != nil {
			return fmt.Errorf("schedule notification retry after %v: %w", deliveryErr, err)
		}
		appLogger.Warnf("notification %s scheduled for retry #%d after failure: %v", notification.GetNotificationId(), retryCount, deliveryErr)
		return deliveryErr
	}

	if err := s.notificationRepo.MarkFailed(ctx, notification.GetNotificationId(), retryCount, deliveryErr.Error()); err != nil {
		return fmt.Errorf("mark notification failed after %v: %w", deliveryErr, err)
	}
	s.publishFailed(ctx, notification, deliveryErr)
	return deliveryErr
}

func (s *Service) shouldRetry(retryCount int32, deliveryErr error) bool {
	if deliveryErr == nil {
		return false
	}
	maxAttempts := 3
	if s.cfg != nil && s.cfg.NotificationRetry.MaxAttempts > 0 {
		maxAttempts = s.cfg.NotificationRetry.MaxAttempts
	}
	if int(retryCount) >= maxAttempts {
		return false
	}
	return !delivery.IsPermanent(deliveryErr)
}

func (s *Service) retryBackoff(retryCount int32) time.Duration {
	if s.cfg != nil && len(s.cfg.NotificationRetry.Backoff) > 0 {
		index := int(retryCount) - 1
		if index < 0 {
			index = 0
		}
		if index < len(s.cfg.NotificationRetry.Backoff) {
			return s.cfg.NotificationRetry.Backoff[index]
		}
		return s.cfg.NotificationRetry.Backoff[len(s.cfg.NotificationRetry.Backoff)-1]
	}
	return time.Minute
}

func (s *Service) publishSentAndDelivered(ctx context.Context, notification *notificationv1.Notification) {
	s.publishSent(ctx, notification)
	if s.publisher != nil {
		_ = s.publisher.PublishNotificationDelivered(ctx, notification.GetNotificationId(), notification.GetRecipientId(), notification.GetTemplateData()["correlation_id"])
	}
	s.enqueueWebhookLifecycleAttempts(ctx, webhookEventDelivered, notification, "")
}

func (s *Service) publishSent(ctx context.Context, notification *notificationv1.Notification) {
	if s.publisher != nil {
		_ = s.publisher.PublishNotificationSent(
			ctx,
			notification.GetNotificationId(),
			notification.GetRecipientId(),
			strings.TrimPrefix(notification.GetChannel().String(), "NOTIFICATION_CHANNEL_"),
			strings.TrimPrefix(notification.GetType().String(), "NOTIFICATION_TYPE_"),
			notification.GetTemplateData()["correlation_id"],
		)
	}
	s.enqueueWebhookLifecycleAttempts(ctx, webhookEventSent, notification, "")
}

func (s *Service) publishFailed(ctx context.Context, notification *notificationv1.Notification, deliveryErr error) {
	if s.publisher != nil {
		_ = s.publisher.PublishNotificationFailed(ctx, notification.GetNotificationId(), notification.GetRecipientId(), deliveryErr.Error(), notification.GetTemplateData()["correlation_id"])
	}
	s.enqueueWebhookLifecycleAttempts(ctx, webhookEventFailed, notification, deliveryErr.Error())
}

func (s *Service) enqueueWebhookLifecycleAttempts(ctx context.Context, lifecycleEvent string, notification *notificationv1.Notification, errorMessage string) {
	if s.webhookRepo == nil || notification == nil || s.cfg == nil || !s.cfg.Webhook.Enabled {
		return
	}

	sourceTopic := strings.TrimSpace(notification.GetTemplateData()["source_topic"])
	sourceGroup := strings.TrimSpace(notification.GetTemplateData()["source_group"])
	channel := strings.TrimPrefix(notification.GetChannel().String(), "NOTIFICATION_CHANNEL_")
	subscriptions, err := s.webhookRepo.ListMatchingSubscriptions(ctx, lifecycleEvent, sourceTopic, sourceGroup, channel)
	if err != nil {
		appLogger.Warnf("failed to list webhook subscriptions for notification %s: %v", notification.GetNotificationId(), err)
		return
	}
	if len(subscriptions) == 0 {
		return
	}

	payload, err := json.Marshal(buildWebhookPayload(lifecycleEvent, notification, errorMessage))
	if err != nil {
		appLogger.Warnf("failed to marshal webhook payload for notification %s: %v", notification.GetNotificationId(), err)
		return
	}

	for _, subscription := range subscriptions {
		if subscription == nil {
			continue
		}
		if err := s.webhookRepo.EnqueueAttempt(ctx, &repository.WebhookDeliveryAttempt{
			SubscriptionID: subscription.SubscriptionID,
			NotificationID: notification.GetNotificationId(),
			LifecycleEvent: lifecycleEvent,
			SourceTopic:    sourceTopic,
			Payload:        json.RawMessage(payload),
			Status:         "QUEUED",
		}); err != nil {
			appLogger.Warnf("failed to enqueue webhook attempt for notification %s subscription %s: %v", notification.GetNotificationId(), subscription.SubscriptionID, err)
		}
	}
}

func (s *Service) runWebhookDispatchCycle(ctx context.Context) error {
	if s.webhookRepo == nil || s.webhookClient == nil || s.cfg == nil || !s.cfg.Webhook.Enabled {
		return nil
	}
	batchSize := s.cfg.Webhook.BatchSize
	if batchSize <= 0 {
		batchSize = 50
	}

	attempts, err := s.webhookRepo.ListDueAttempts(ctx, s.now().UTC(), batchSize)
	if err != nil {
		return fmt.Errorf("list due webhook attempts: %w", err)
	}

	var firstErr error
	for _, attempt := range attempts {
		if err := s.dispatchWebhookAttempt(ctx, attempt); err != nil && ctx.Err() == nil {
			appLogger.Errorf("dispatch webhook attempt %s failed: %v", attempt.AttemptID, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (s *Service) dispatchWebhookAttempt(ctx context.Context, attempt *repository.WebhookDeliveryAttempt) error {
	if attempt == nil {
		return nil
	}
	if err := s.webhookRepo.MarkAttemptSending(ctx, attempt.AttemptID); err != nil {
		return fmt.Errorf("mark webhook attempt sending: %w", err)
	}

	subscription, err := s.loadWebhookSubscription(ctx, attempt.SubscriptionID)
	if err != nil {
		return s.handleWebhookFailure(ctx, attempt, nil, nil, err)
	}

	resp, err := s.webhookClient.Send(ctx, &delivery.WebhookRequest{
		TargetURL:    subscription.TargetURL,
		Secret:       subscription.Secret,
		EventType:    attempt.LifecycleEvent,
		Payload:      json.RawMessage(attempt.Payload),
		Timeout:      time.Duration(subscription.TimeoutSeconds) * time.Second,
		Subscription: subscription.SubscriptionID,
	})
	if err != nil {
		return s.handleWebhookFailure(ctx, attempt, subscription, resp, err)
	}
	if err := s.webhookRepo.MarkAttemptSucceeded(ctx, attempt.AttemptID, resp.StatusCode, resp.Body, s.now().UTC()); err != nil {
		return fmt.Errorf("mark webhook attempt succeeded: %w", err)
	}
	return nil
}

func (s *Service) loadWebhookSubscription(ctx context.Context, subscriptionID string) (*repository.WebhookSubscription, error) {
	subscription, err := s.webhookRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, delivery.Permanent(fmt.Errorf("webhook subscription %s not found", subscriptionID))
		}
		return nil, err
	}
	if subscription == nil || !subscription.IsActive {
		return nil, delivery.Permanent(fmt.Errorf("webhook subscription %s is inactive", subscriptionID))
	}
	return subscription, nil
}

func (s *Service) handleWebhookFailure(ctx context.Context, attempt *repository.WebhookDeliveryAttempt, subscription *repository.WebhookSubscription, resp *delivery.WebhookResponse, deliveryErr error) error {
	retryCount := attempt.RetryCount + 1
	responseStatus := 0
	responseBody := ""
	if resp != nil {
		responseStatus = resp.StatusCode
		responseBody = resp.Body
	}

	if s.shouldRetryWebhook(retryCount, subscription, deliveryErr) {
		nextAttempt := s.now().UTC().Add(s.webhookRetryBackoff(retryCount))
		if err := s.webhookRepo.ScheduleAttemptRetry(ctx, attempt.AttemptID, retryCount, nextAttempt, responseStatus, responseBody, deliveryErr.Error()); err != nil {
			return fmt.Errorf("schedule webhook retry after %v: %w", deliveryErr, err)
		}
		return deliveryErr
	}

	if err := s.webhookRepo.MarkAttemptFailed(ctx, attempt.AttemptID, retryCount, responseStatus, responseBody, deliveryErr.Error(), s.now().UTC()); err != nil {
		return fmt.Errorf("mark webhook attempt failed after %v: %w", deliveryErr, err)
	}
	return deliveryErr
}

func (s *Service) shouldRetryWebhook(retryCount int, subscription *repository.WebhookSubscription, deliveryErr error) bool {
	if deliveryErr == nil || delivery.IsPermanent(deliveryErr) {
		return false
	}
	maxAttempts := 5
	if s.cfg != nil && s.cfg.Webhook.MaxAttempts > 0 {
		maxAttempts = s.cfg.Webhook.MaxAttempts
	}
	if retryCount >= maxAttempts {
		return false
	}
	if subscription != nil && subscription.MaxAttempts > 0 && retryCount >= subscription.MaxAttempts {
		return false
	}
	return true
}

func (s *Service) webhookRetryBackoff(retryCount int) time.Duration {
	if s.cfg != nil && len(s.cfg.Webhook.Backoff) > 0 {
		index := retryCount - 1
		if index < 0 {
			index = 0
		}
		if index < len(s.cfg.Webhook.Backoff) {
			return s.cfg.Webhook.Backoff[index]
		}
		return s.cfg.Webhook.Backoff[len(s.cfg.Webhook.Backoff)-1]
	}
	return time.Minute
}

func buildWebhookPayload(lifecycleEvent string, notification *notificationv1.Notification, errorMessage string) map[string]any {
	payload := map[string]any{
		"event_type":      lifecycleEvent,
		"notification_id": notification.GetNotificationId(),
		"recipient_id":    notification.GetRecipientId(),
		"channel":         strings.TrimPrefix(notification.GetChannel().String(), "NOTIFICATION_CHANNEL_"),
		"type":            strings.TrimPrefix(notification.GetType().String(), "NOTIFICATION_TYPE_"),
		"priority":        strings.TrimPrefix(notification.GetPriority().String(), "NOTIFICATION_PRIORITY_"),
		"subject":         notification.GetSubject(),
		"message":         notification.GetMessage(),
		"template_data":   cloneMap(notification.GetTemplateData()),
		"occurred_at":     time.Now().UTC().Format(time.RFC3339),
	}
	if strings.TrimSpace(errorMessage) != "" {
		payload["error_message"] = errorMessage
	}
	return payload
}
