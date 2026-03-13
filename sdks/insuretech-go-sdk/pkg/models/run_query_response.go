package models


// RunQueryResponse represents a run_query_response
type RunQueryResponse struct {
	Rows []*Row `json:"rows,omitempty"`
	Columns []string `json:"columns,omitempty"`
	RowCount int `json:"row_count,omitempty"`
	ExecutionTimeMs float64 `json:"execution_time_ms,omitempty"`
	Error *Error `json:"error,omitempty"`
}
