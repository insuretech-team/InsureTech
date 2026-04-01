package models


// ChangePasswordRequest represents a change_password_request
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	UserId string `json:"user_id"`
}
