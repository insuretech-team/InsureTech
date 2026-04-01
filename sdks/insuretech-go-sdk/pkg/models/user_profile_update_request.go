package models

import (
	"time"
)

// UserProfileUpdateRequest represents a user_profile_update_request
type UserProfileUpdateRequest struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	BloodGroup string `json:"blood_group,omitempty"`
	City string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	District string `json:"district,omitempty"`
	Division string `json:"division,omitempty"`
	EmergencyContactName string `json:"emergency_contact_name,omitempty"`
	EmergencyContactNumber string `json:"emergency_contact_number,omitempty"`
	Employer string `json:"employer,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	IncomeRange string `json:"income_range,omitempty"`
	MaritalStatus string `json:"marital_status,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	PermanentAddress string `json:"permanent_address,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	ProfilePhotoUrl string `json:"profile_photo_url,omitempty"`
	UserId string `json:"user_id"`
}
