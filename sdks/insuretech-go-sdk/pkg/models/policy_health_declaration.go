package models


// PolicyHealthDeclaration represents a policy_health_declaration
type PolicyHealthDeclaration struct {
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	Conditions []string `json:"conditions,omitempty"`
	IsSmoker bool `json:"is_smoker,omitempty"`
	BloodGroup string `json:"blood_group,omitempty"`
}
