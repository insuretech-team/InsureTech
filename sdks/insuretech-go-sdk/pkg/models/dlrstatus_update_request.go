package models


// DLRStatusUpdateRequest represents a dlrstatus_update_request
type DLRStatusUpdateRequest struct {
	ProviderMessageId string `json:"provider_message_id"`
	Status string `json:"status,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}
