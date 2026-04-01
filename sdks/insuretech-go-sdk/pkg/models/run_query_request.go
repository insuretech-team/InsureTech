package models


// RunQueryRequest represents a run_query_request
type RunQueryRequest struct {
	Limit int `json:"limit,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Query string `json:"query"`
}
