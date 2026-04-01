package models


// InsuredPerson represents a insured_person
type InsuredPerson struct {
	Age int `json:"age,omitempty"`
	Email string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	HealthConditions []*HealthCondition `json:"health_conditions,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	Phone string `json:"phone,omitempty"`
}
