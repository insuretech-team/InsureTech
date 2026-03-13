package models

import (
	"time"
)

// UserProfile represents a user_profile
type UserProfile struct {
	FullName string `json:"full_name,omitempty"`
	Gender *AuthnGender `json:"gender,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country string `json:"country,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	KycVerifiedAt time.Time `json:"kyc_verified_at,omitempty"`
	PhotographSelfieUrl string `json:"photograph_selfie_url,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	EmergencyContactNumber string `json:"emergency_contact_number,omitempty"`
	ProfilePhotoUrl string `json:"profile_photo_url,omitempty"`
	MaritalStatus string `json:"marital_status,omitempty"`
	Employer string `json:"employer,omitempty"`
	IdUploadFrontUrl string `json:"id_upload_front_url,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	City string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	PermanentAddress string `json:"permanent_address,omitempty"`
	EmergencyContactName string `json:"emergency_contact_name,omitempty"`
	IdType string `json:"id_type,omitempty"`
	IdUploadBackUrl string `json:"id_upload_back_url,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Division string `json:"division,omitempty"`
	ProofOfAddressUrl string `json:"proof_of_address_url,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	KycVerified bool `json:"kyc_verified,omitempty"`
	ConsentPrivacyAcceptance bool `json:"consent_privacy_acceptance,omitempty"`
}
