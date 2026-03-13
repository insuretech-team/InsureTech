package models


// EmployeeCreationRequest represents a employee_creation_request
type EmployeeCreationRequest struct {
	Gender *EmployeeGender `json:"gender,omitempty"`
	Name string `json:"name"`
	DepartmentId string `json:"department_id"`
	BusinessId string `json:"business_id"`
	NumberOfDependent int `json:"number_of_dependent,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	DateOfJoining string `json:"date_of_joining,omitempty"`
	EmployeeId string `json:"employee_id"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	AssignedPlanId string `json:"assigned_plan_id"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	Email string `json:"email"`
	MobileNumber string `json:"mobile_number,omitempty"`
}
