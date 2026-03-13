package models


// InsurerRetrievalResponse represents a insurer_retrieval_response
type InsurerRetrievalResponse struct {
	Insurer *Insurer `json:"insurer,omitempty"`
	Config *InsurerConfig `json:"config,omitempty"`
	Error *Error `json:"error,omitempty"`
}
