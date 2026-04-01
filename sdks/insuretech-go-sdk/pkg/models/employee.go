package models

import (
	"time"
)

// Employee represents a employee
type Employee struct {
	AssignedPlanId string `json:"assigned_plan_id,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	DateOfJoining string `json:"date_of_joining,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	Email string `json:"email,omitempty"`
	EmployeeId string `json:"employee_id,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	Gender *EmployeeGender `json:"gender,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Name string `json:"name,omitempty"`
	NumberOfDependent int `json:"number_of_dependent,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	Status *EmployeeStatus `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
