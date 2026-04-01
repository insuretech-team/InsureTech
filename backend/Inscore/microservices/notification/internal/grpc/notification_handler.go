package grpc

import (
	"context"

	notificationservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/services/v1"
)

type NotificationServiceHandler struct {
	notificationservicev1.UnimplementedNotificationServiceServer
	service NotificationServiceIface
}

func NewNotificationServiceHandler(service NotificationServiceIface) *NotificationServiceHandler {
	return &NotificationServiceHandler{service: service}
}

func (h *NotificationServiceHandler) SendNotification(ctx context.Context, req *notificationservicev1.SendNotificationRequest) (*notificationservicev1.SendNotificationResponse, error) {
	resp, err := h.service.SendNotification(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) SendBulkNotifications(ctx context.Context, req *notificationservicev1.SendBulkNotificationsRequest) (*notificationservicev1.SendBulkNotificationsResponse, error) {
	resp, err := h.service.SendBulkNotifications(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) GetNotificationStatus(ctx context.Context, req *notificationservicev1.GetNotificationStatusRequest) (*notificationservicev1.GetNotificationStatusResponse, error) {
	resp, err := h.service.GetNotificationStatus(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) GetUserNotifications(ctx context.Context, req *notificationservicev1.GetUserNotificationsRequest) (*notificationservicev1.GetUserNotificationsResponse, error) {
	resp, err := h.service.GetUserNotifications(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) MarkAsRead(ctx context.Context, req *notificationservicev1.MarkAsReadRequest) (*notificationservicev1.MarkAsReadResponse, error) {
	resp, err := h.service.MarkAsRead(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) UpdatePreferences(ctx context.Context, req *notificationservicev1.UpdatePreferencesRequest) (*notificationservicev1.UpdatePreferencesResponse, error) {
	resp, err := h.service.UpdatePreferences(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) CreateNotificationTemplate(ctx context.Context, req *notificationservicev1.CreateNotificationTemplateRequest) (*notificationservicev1.CreateNotificationTemplateResponse, error) {
	resp, err := h.service.CreateNotificationTemplate(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) UpdateNotificationTemplate(ctx context.Context, req *notificationservicev1.UpdateNotificationTemplateRequest) (*notificationservicev1.UpdateNotificationTemplateResponse, error) {
	resp, err := h.service.UpdateNotificationTemplate(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *NotificationServiceHandler) DeactivateNotificationTemplate(ctx context.Context, req *notificationservicev1.DeactivateNotificationTemplateRequest) (*notificationservicev1.DeactivateNotificationTemplateResponse, error) {
	resp, err := h.service.DeactivateNotificationTemplate(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
