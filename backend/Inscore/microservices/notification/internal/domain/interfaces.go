package domain

import (
	"context"

	notificationservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/services/v1"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req *notificationservicev1.SendNotificationRequest) (*notificationservicev1.SendNotificationResponse, error)
	SendBulkNotifications(ctx context.Context, req *notificationservicev1.SendBulkNotificationsRequest) (*notificationservicev1.SendBulkNotificationsResponse, error)
	GetNotificationStatus(ctx context.Context, req *notificationservicev1.GetNotificationStatusRequest) (*notificationservicev1.GetNotificationStatusResponse, error)
	GetUserNotifications(ctx context.Context, req *notificationservicev1.GetUserNotificationsRequest) (*notificationservicev1.GetUserNotificationsResponse, error)
	MarkAsRead(ctx context.Context, req *notificationservicev1.MarkAsReadRequest) (*notificationservicev1.MarkAsReadResponse, error)
	UpdatePreferences(ctx context.Context, req *notificationservicev1.UpdatePreferencesRequest) (*notificationservicev1.UpdatePreferencesResponse, error)
	CreateNotificationTemplate(ctx context.Context, req *notificationservicev1.CreateNotificationTemplateRequest) (*notificationservicev1.CreateNotificationTemplateResponse, error)
	UpdateNotificationTemplate(ctx context.Context, req *notificationservicev1.UpdateNotificationTemplateRequest) (*notificationservicev1.UpdateNotificationTemplateResponse, error)
	DeactivateNotificationTemplate(ctx context.Context, req *notificationservicev1.DeactivateNotificationTemplateRequest) (*notificationservicev1.DeactivateNotificationTemplateResponse, error)
}
