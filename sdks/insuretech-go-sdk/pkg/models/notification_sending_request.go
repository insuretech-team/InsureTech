package models


// NotificationSendingRequest represents a notification_sending_request
type NotificationSendingRequest struct {
	Channel *NotificationChannel `json:"channel,omitempty"`
	Message string `json:"message,omitempty"`
	Priority *NotificationPriority `json:"priority,omitempty"`
	RecipientId string `json:"recipient_id"`
	ScheduleAfterSeconds string `json:"schedule_after_seconds,omitempty"`
	Subject string `json:"subject,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	TemplateId string `json:"template_id"`
	Type *NotificationType `json:"type"`
}
