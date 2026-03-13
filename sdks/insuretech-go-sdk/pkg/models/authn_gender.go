package models

// AuthnGender represents a authn_gender
type AuthnGender string

// AuthnGender values
const (
	AuthnGenderGENDERUNSPECIFIED AuthnGender = "GENDER_UNSPECIFIED"
	AuthnGenderGENDERMALE  = "GENDER_MALE"
	AuthnGenderGENDERFEMALE  = "GENDER_FEMALE"
	AuthnGenderGENDEROTHER  = "GENDER_OTHER"
)
