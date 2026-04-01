package models


// CSRFValidationRequest represents a csrfvalidation_request
type CSRFValidationRequest struct {
	CsrfToken string `json:"csrf_token,omitempty"`
	SessionId string `json:"session_id"`
}
