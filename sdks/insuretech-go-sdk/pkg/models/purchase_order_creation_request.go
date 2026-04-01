package models


// PurchaseOrderCreationRequest represents a purchase_order_creation_request
type PurchaseOrderCreationRequest struct {
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	DepartmentId string `json:"department_id"`
	EmployeeCount int `json:"employee_count,omitempty"`
	Notes string `json:"notes,omitempty"`
	NumberOfDependents int `json:"number_of_dependents,omitempty"`
	PlanId string `json:"plan_id"`
	RequestedBy string `json:"requested_by,omitempty"`
}
