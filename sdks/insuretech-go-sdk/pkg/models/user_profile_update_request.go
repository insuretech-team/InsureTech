package models

import (
	"time"
)

// UserProfileUpdateRequest represents a user_profile_update_request
type UserProfileUpdateRequest struct {
	FullName string `json:"full_name,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	District string `json:"district,omitempty"`
	Country string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	MaritalStatus string `json:"marital_status,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	EmergencyContactName string `json:"emergency_contact_name,omitempty"`
	EmergencyContactNumber string `json:"emergency_contact_number,omitempty"`
	IncomeRange string `json:"income_range,omitempty"`
	Employer string `json:"employer,omitempty"`
	City string `json:"city,omitempty"`
	Division string `json:"division,omitempty"`
	PermanentAddress string `json:"permanent_address,omitempty"`
	UserId string `json:"user_id"`
	Gender string `json:"gender,omitempty"`
	ProfilePhotoUrl string `json:"profile_photo_url,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	BloodGroup string `json:"blood_group,omitempty"`
}
