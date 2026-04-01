package models


// NotificationTemplateUpdateRequest represents a notification_template_update_request
type NotificationTemplateUpdateRequest struct {
	Body string `json:"body,omitempty"`
	Name string `json:"name"`
	Subject string `json:"subject,omitempty"`
	TemplateId string `json:"template_id"`
}
