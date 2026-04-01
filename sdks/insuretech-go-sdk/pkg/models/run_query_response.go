package models


// RunQueryResponse represents a run_query_response
type RunQueryResponse struct {
	Columns []string `json:"columns,omitempty"`
	ExecutionTimeMs float64 `json:"execution_time_ms,omitempty"`
	RowCount int `json:"row_count,omitempty"`
	Rows []*Row `json:"rows,omitempty"`
}
