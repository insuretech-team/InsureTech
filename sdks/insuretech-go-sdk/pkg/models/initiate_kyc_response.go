package models


// InitiateKYCResponse represents a initiate_kyc_response
type InitiateKYCResponse struct {
	KycId string `json:"kyc_id,omitempty"`
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
