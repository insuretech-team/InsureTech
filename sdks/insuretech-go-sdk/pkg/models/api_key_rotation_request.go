package models


// APIKeyRotationRequest represents a api_key_rotation_request
type APIKeyRotationRequest struct {
	GracePeriodHours int `json:"grace_period_hours,omitempty"`
	KeyId string `json:"key_id"`
}
