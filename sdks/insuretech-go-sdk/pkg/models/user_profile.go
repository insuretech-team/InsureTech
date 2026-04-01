package models

import (
	"time"
)

// UserProfile represents a user_profile
type UserProfile struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City string `json:"city,omitempty"`
	ConsentPrivacyAcceptance bool `json:"consent_privacy_acceptance,omitempty"`
	Country string `json:"country,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	District string `json:"district,omitempty"`
	Division string `json:"division,omitempty"`
	EmergencyContactName string `json:"emergency_contact_name,omitempty"`
	EmergencyContactNumber string `json:"emergency_contact_number,omitempty"`
	Employer string `json:"employer,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Gender *AuthnGender `json:"gender,omitempty"`
	IdType string `json:"id_type,omitempty"`
	IdUploadBackUrl string `json:"id_upload_back_url,omitempty"`
	IdUploadFrontUrl string `json:"id_upload_front_url,omitempty"`
	KycVerified bool `json:"kyc_verified,omitempty"`
	KycVerifiedAt time.Time `json:"kyc_verified_at,omitempty"`
	MaritalStatus string `json:"marital_status,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	PermanentAddress string `json:"permanent_address,omitempty"`
	PhotographSelfieUrl string `json:"photograph_selfie_url,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	ProfilePhotoUrl string `json:"profile_photo_url,omitempty"`
	ProofOfAddressUrl string `json:"proof_of_address_url,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
