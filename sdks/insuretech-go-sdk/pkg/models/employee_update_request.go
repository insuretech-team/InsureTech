package models


// EmployeeUpdateRequest represents a employee_update_request
type EmployeeUpdateRequest struct {
	Status *EmployeeStatus `json:"status,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	Name string `json:"name"`
	DepartmentId string `json:"department_id"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Gender *EmployeeGender `json:"gender,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	NumberOfDependent int `json:"number_of_dependent,omitempty"`
	Email string `json:"email"`
	MobileNumber string `json:"mobile_number,omitempty"`
	DateOfJoining string `json:"date_of_joining,omitempty"`
	AssignedPlanId string `json:"assigned_plan_id"`
}
