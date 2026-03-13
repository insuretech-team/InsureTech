package models


// APIKeyRotationRequest represents a api_key_rotation_request
type APIKeyRotationRequest struct {
	KeyId string `json:"key_id"`
	GracePeriodHours int `json:"grace_period_hours,omitempty"`
}
