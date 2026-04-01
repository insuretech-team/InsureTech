package models


// DLRStatusUpdateRequest represents a dlrstatus_update_request
type DLRStatusUpdateRequest struct {
	ErrorCode string `json:"error_code,omitempty"`
	ProviderMessageId string `json:"provider_message_id"`
	Status string `json:"status,omitempty"`
}
