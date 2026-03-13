package models


// EmployeeView represents a employee_view
type EmployeeView struct {
	DepartmentName string `json:"department_name,omitempty"`
	AssignedPlanName string `json:"assigned_plan_name,omitempty"`
	Employee *Employee `json:"employee,omitempty"`
}
