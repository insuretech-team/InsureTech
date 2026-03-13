package models


// PurchaseOrderCreationRequest represents a purchase_order_creation_request
type PurchaseOrderCreationRequest struct {
	DepartmentId string `json:"department_id"`
	PlanId string `json:"plan_id"`
	EmployeeCount int `json:"employee_count,omitempty"`
	NumberOfDependents int `json:"number_of_dependents,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	Notes string `json:"notes,omitempty"`
}
