package models


// NotificationTemplateCreationRequest represents a notification_template_creation_request
type NotificationTemplateCreationRequest struct {
	Name string `json:"name"`
	Type *NotificationType `json:"type"`
	Subject string `json:"subject,omitempty"`
	Body string `json:"body,omitempty"`
	Variables []string `json:"variables,omitempty"`
}
