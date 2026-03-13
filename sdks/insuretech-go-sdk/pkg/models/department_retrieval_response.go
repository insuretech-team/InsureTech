package models


// DepartmentRetrievalResponse represents a department_retrieval_response
type DepartmentRetrievalResponse struct {
	Department *Department `json:"department,omitempty"`
	Error *Error `json:"error,omitempty"`
}
