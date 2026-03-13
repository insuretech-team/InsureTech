package models


// NotificationTemplateUpdateRequest represents a notification_template_update_request
type NotificationTemplateUpdateRequest struct {
	TemplateId string `json:"template_id"`
	Name string `json:"name"`
	Subject string `json:"subject,omitempty"`
	Body string `json:"body,omitempty"`
}
