package models

import (
	"time"
)

// Department represents a department
type Department struct {
	BusinessId string `json:"business_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EmployeeNo int `json:"employee_no,omitempty"`
	Name string `json:"name,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
