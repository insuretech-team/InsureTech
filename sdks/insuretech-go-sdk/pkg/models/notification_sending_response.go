package models


// NotificationSendingResponse represents a notification_sending_response
type NotificationSendingResponse struct {
	NotificationId string `json:"notification_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
