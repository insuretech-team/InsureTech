package models


// MyEmployeeProfileRetrievalResponse represents a my_employee_profile_retrieval_response
type MyEmployeeProfileRetrievalResponse struct {
	Employee *EmployeeView `json:"employee,omitempty"`
}
