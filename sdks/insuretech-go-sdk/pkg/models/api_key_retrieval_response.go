package models


// ApiKeyRetrievalResponse represents a api_key_retrieval_response
type ApiKeyRetrievalResponse struct {
	ApiKey *ApiKey `json:"api_key,omitempty"`
	Error *Error `json:"error,omitempty"`
}
