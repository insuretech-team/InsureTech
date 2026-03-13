package models

import (
	"time"
)

// UserProfileCreationRequest represents a user_profile_creation_request
type UserProfileCreationRequest struct {
	NidNumber string `json:"nid_number,omitempty"`
	EmergencyContactName string `json:"emergency_contact_name,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	PermanentAddress string `json:"permanent_address,omitempty"`
	MaritalStatus string `json:"marital_status,omitempty"`
	IncomeRange string `json:"income_range,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	District string `json:"district,omitempty"`
	Division string `json:"division,omitempty"`
	BloodGroup string `json:"blood_group,omitempty"`
	Gender string `json:"gender,omitempty"`
	Country string `json:"country,omitempty"`
	EmergencyContactNumber string `json:"emergency_contact_number,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	Employer string `json:"employer,omitempty"`
	UserId string `json:"user_id"`
	FullName string `json:"full_name,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	City string `json:"city,omitempty"`
}
