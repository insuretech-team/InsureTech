package models


// MyEmployeeCoverageRetrievalResponse represents a my_employee_coverage_retrieval_response
type MyEmployeeCoverageRetrievalResponse struct {
	Coverage *EmployeeCoverageView `json:"coverage,omitempty"`
}
