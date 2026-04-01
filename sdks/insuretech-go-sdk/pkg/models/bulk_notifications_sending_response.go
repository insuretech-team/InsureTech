package models


// BulkNotificationsSendingResponse represents a bulk_notifications_sending_response
type BulkNotificationsSendingResponse struct {
	FailedCount int `json:"failed_count,omitempty"`
	NotificationIds []string `json:"notification_ids,omitempty"`
	SuccessCount int `json:"success_count,omitempty"`
}
