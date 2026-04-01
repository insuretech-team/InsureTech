package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// NotificationService handles notification-related API calls
type NotificationService struct {
	Client Client
}

// CreateNotificationTemplate Create notification template
func (s *NotificationService) CreateNotificationTemplate(ctx context.Context, req *models.NotificationTemplateCreationRequest) error {
	path := "/v1/notification-templates"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdateNotificationTemplate Update notification template
func (s *NotificationService) UpdateNotificationTemplate(ctx context.Context, templateId string, req *models.NotificationTemplateUpdateRequest) error {
	path := "/v1/notification-templates/{template_id}"
	path = strings.ReplaceAll(path, "{template_id}", templateId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeactivateNotificationTemplate Deactivate notification template
func (s *NotificationService) DeactivateNotificationTemplate(ctx context.Context, templateId string, req *models.NotificationTemplateDeactivationRequest) error {
	path := "/v1/notification-templates/{template_id}:deactivate"
	path = strings.ReplaceAll(path, "{template_id}", templateId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SendNotification Send notification
func (s *NotificationService) SendNotification(ctx context.Context, req *models.NotificationSendingRequest) error {
	path := "/v1/notifications"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetNotificationStatus Get notification status
func (s *NotificationService) GetNotificationStatus(ctx context.Context, notificationId string) error {
	path := "/v1/notifications/{notification_id}"
	path = strings.ReplaceAll(path, "{notification_id}", notificationId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// MarkAsRead Mark as read
func (s *NotificationService) MarkAsRead(ctx context.Context, req *models.MarkAsReadRequest) error {
	path := "/v1/notifications:mark-as-read"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SendBulkNotifications Send bulk notifications
func (s *NotificationService) SendBulkNotifications(ctx context.Context, req *models.BulkNotificationsSendingRequest) error {
	path := "/v1/notifications:send-bulk"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdatePreferences Update notification preferences
func (s *NotificationService) UpdatePreferences(ctx context.Context, userId string, req *models.PreferencesUpdateRequest) error {
	path := "/v1/users/{user_id}/notification-preferences"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "PUT", path, req, nil)
}

// GetUserNotifications Get user notifications
func (s *NotificationService) GetUserNotifications(ctx context.Context, userId string) error {
	path := "/v1/users/{user_id}/notifications"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

