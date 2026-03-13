package models


// RunQueryRequest represents a run_query_request
type RunQueryRequest struct {
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Limit int `json:"limit,omitempty"`
	Query string `json:"query"`
}
