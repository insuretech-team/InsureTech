package models


// InitiateKYCResponse represents a initiate_kyc_response
type InitiateKYCResponse struct {
	KycId string `json:"kyc_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	ProviderReference string `json:"provider_reference,omitempty"`
	SessionState string `json:"session_state,omitempty"`
	Status string `json:"status,omitempty"`
	Steps []*KYCStep `json:"steps,omitempty"`
	TotalTimeoutSeconds int `json:"total_timeout_seconds,omitempty"`
}
