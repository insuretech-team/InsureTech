package models


// SetTemporaryPasswordRequest represents a set_temporary_password_request
type SetTemporaryPasswordRequest struct {
	AssignedBy string `json:"assigned_by,omitempty"`
	RequirePasswordChange bool `json:"require_password_change,omitempty"`
	TemporaryPassword string `json:"temporary_password,omitempty"`
	UserId string `json:"user_id"`
}
