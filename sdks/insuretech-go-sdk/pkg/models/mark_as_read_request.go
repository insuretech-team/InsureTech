package models


// MarkAsReadRequest represents a mark_as_read_request
type MarkAsReadRequest struct {
	NotificationIds []string `json:"notification_ids"`
}
