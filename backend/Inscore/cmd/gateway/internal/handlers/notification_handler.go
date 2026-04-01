package handlers

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/services/v1"
)

// NotificationHandler proxies notification requests to the notification gRPC service.
// BUG-006 FIX: Replaces PoliSyncHandler (HTTP proxy) which cannot reach the
// gRPC-only notification service (no HTTP companion server on :50231).
type NotificationHandler struct {
	client notificationv1.NotificationServiceClient
}

// NewNotificationHandler creates a NotificationHandler from a gRPC connection.
func NewNotificationHandler(conn *grpc.ClientConn) *NotificationHandler {
	return &NotificationHandler{client: notificationv1.NewNotificationServiceClient(conn)}
}

// Send handles POST /v1/notifications
func (h *NotificationHandler) Send(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req notificationv1.SendNotificationRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SendNotification(ctx, &req)
	})
}

// GetStatus handles GET /v1/notifications/{notification_id}
func (h *NotificationHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	notifID := r.PathValue("notification_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetNotificationStatus(ctx, &notificationv1.GetNotificationStatusRequest{
			NotificationId: notifID,
		})
	})
}

// GetUserNotifications handles GET /v1/users/{user_id}/notifications
func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetUserNotifications(ctx, &notificationv1.GetUserNotificationsRequest{
			UserId: userID,
		})
	})
}

// MarkAsRead handles POST /v1/notifications/mark-as-read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req notificationv1.MarkAsReadRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.MarkAsRead(ctx, &req)
	})
}

// SendBulk handles POST /v1/notifications/send-bulk
func (h *NotificationHandler) SendBulk(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req notificationv1.SendBulkNotificationsRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SendBulkNotifications(ctx, &req)
	})
}

// CreateTemplate handles POST /v1/notification-templates
func (h *NotificationHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req notificationv1.CreateNotificationTemplateRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateNotificationTemplate(ctx, &req)
	})
}

// UpdateTemplate handles PATCH /v1/notification-templates/{template_id}
func (h *NotificationHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req notificationv1.UpdateNotificationTemplateRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.UpdateNotificationTemplate(ctx, &req)
	})
}

// DeactivateTemplate handles POST /v1/notification-templates/{template_id}/deactivate
func (h *NotificationHandler) DeactivateTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeactivateNotificationTemplate(ctx, &notificationv1.DeactivateNotificationTemplateRequest{
			TemplateId: templateID,
		})
	})
}
