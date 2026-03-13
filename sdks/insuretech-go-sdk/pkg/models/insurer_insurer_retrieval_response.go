package models


// InsurerInsurerRetrievalResponse represents a insurer_insurer_retrieval_response
type InsurerInsurerRetrievalResponse struct {
	Config *InsurerConfig `json:"config,omitempty"`
	Error *Error `json:"error,omitempty"`
	Insurer *Insurer `json:"insurer,omitempty"`
}
