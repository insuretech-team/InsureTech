package models


// NotificationTemplateCreationRequest represents a notification_template_creation_request
type NotificationTemplateCreationRequest struct {
	Body string `json:"body,omitempty"`
	Name string `json:"name"`
	Subject string `json:"subject,omitempty"`
	Type *NotificationType `json:"type"`
	Variables []string `json:"variables,omitempty"`
}
