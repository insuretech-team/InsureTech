package models


// EmployeeRetrievalResponse represents a employee_retrieval_response
type EmployeeRetrievalResponse struct {
	Employee *EmployeeView `json:"employee,omitempty"`
	Error *Error `json:"error,omitempty"`
}
