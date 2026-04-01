package models


// PolicyHealthDeclaration represents a policy_health_declaration
type PolicyHealthDeclaration struct {
	BloodGroup string `json:"blood_group,omitempty"`
	Conditions []string `json:"conditions,omitempty"`
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	IsSmoker bool `json:"is_smoker,omitempty"`
}
