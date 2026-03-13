package models


// RequestEndorsementResponse represents a request_endorsement_response
type RequestEndorsementResponse struct {
	EndorsementId string `json:"endorsement_id,omitempty"`
	EndorsementNumber string `json:"endorsement_number,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
