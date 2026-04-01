package models


// PrimaryContact represents a primary_contact
type PrimaryContact struct {
	Department string `json:"department,omitempty"`
	Email string `json:"email,omitempty"`
	Name string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}
