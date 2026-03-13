package models

import (
	"time"
)

// Employee represents a employee
type Employee struct {
	Name string `json:"name,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	DateOfJoining string `json:"date_of_joining,omitempty"`
	UserId string `json:"user_id,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	NumberOfDependent int `json:"number_of_dependent,omitempty"`
	EmployeeId string `json:"employee_id,omitempty"`
	AssignedPlanId string `json:"assigned_plan_id,omitempty"`
	Email string `json:"email,omitempty"`
	Gender *EmployeeGender `json:"gender,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	Status *EmployeeStatus `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
}
