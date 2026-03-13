package models


// NotificationTemplateDeactivationRequest represents a notification_template_deactivation_request
type NotificationTemplateDeactivationRequest struct {
	Reason string `json:"reason,omitempty"`
	TemplateId string `json:"template_id"`
}
