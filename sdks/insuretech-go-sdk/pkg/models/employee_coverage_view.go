package models


// EmployeeCoverageView represents a employee_coverage_view
type EmployeeCoverageView struct {
	AssignedPlanId string `json:"assigned_plan_id,omitempty"`
	AssignedPlanName string `json:"assigned_plan_name,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	EmployeeId string `json:"employee_id,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	NumberOfDependent int `json:"number_of_dependent,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	OrganisationName string `json:"organisation_name,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
}
