package models


// PrimaryContact represents a primary_contact
type PrimaryContact struct {
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
}
