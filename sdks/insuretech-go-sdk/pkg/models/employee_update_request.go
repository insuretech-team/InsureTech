package models


// EmployeeUpdateRequest represents a employee_update_request
type EmployeeUpdateRequest struct {
	AssignedPlanId string `json:"assigned_plan_id"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	DateOfJoining string `json:"date_of_joining,omitempty"`
	DepartmentId string `json:"department_id"`
	Email string `json:"email"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	Gender *EmployeeGender `json:"gender,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Name string `json:"name"`
	NumberOfDependent int `json:"number_of_dependent,omitempty"`
	Status *EmployeeStatus `json:"status,omitempty"`
	UserId string `json:"user_id"`
}
