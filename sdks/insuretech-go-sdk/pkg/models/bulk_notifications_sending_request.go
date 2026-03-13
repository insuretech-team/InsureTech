package models


// BulkNotificationsSendingRequest represents a bulk_notifications_sending_request
type BulkNotificationsSendingRequest struct {
	Notifications []*NotificationSendingRequest `json:"notifications"`
}
