package models


// NotificationTemplateUpdateResponse represents a notification_template_update_response
type NotificationTemplateUpdateResponse struct {
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
