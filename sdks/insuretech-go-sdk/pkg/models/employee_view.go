package models


// EmployeeView represents a employee_view
type EmployeeView struct {
	AssignedPlanName string `json:"assigned_plan_name,omitempty"`
	DepartmentName string `json:"department_name,omitempty"`
	Employee *Employee `json:"employee,omitempty"`
}
