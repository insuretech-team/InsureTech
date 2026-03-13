package models

// BeneficiaryGender represents a beneficiary_gender
type BeneficiaryGender string

// BeneficiaryGender values
const (
	BeneficiaryGenderGENDERUNSPECIFIED BeneficiaryGender = "GENDER_UNSPECIFIED"
	BeneficiaryGenderGENDERMALE  = "GENDER_MALE"
	BeneficiaryGenderGENDERFEMALE  = "GENDER_FEMALE"
	BeneficiaryGenderGENDEROTHER  = "GENDER_OTHER"
)
