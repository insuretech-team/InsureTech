package models


// NotificationTemplateDeactivationResponse represents a notification_template_deactivation_response
type NotificationTemplateDeactivationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
