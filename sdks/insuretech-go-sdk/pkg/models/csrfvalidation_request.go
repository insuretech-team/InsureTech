package models


// CSRFValidationRequest represents a csrfvalidation_request
type CSRFValidationRequest struct {
	SessionId string `json:"session_id"`
	CsrfToken string `json:"csrf_token,omitempty"`
}
