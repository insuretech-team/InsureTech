package models

// CommissionType represents a commission_type
type CommissionType string

// CommissionType values
const (
	CommissionTypeCOMMISSIONTYPEUNSPECIFIED CommissionType = "COMMISSION_TYPE_UNSPECIFIED"
	CommissionTypeCOMMISSIONTYPEACQUISITION  = "COMMISSION_TYPE_ACQUISITION"
	CommissionTypeCOMMISSIONTYPERENEWAL  = "COMMISSION_TYPE_RENEWAL"
	CommissionTypeCOMMISSIONTYPECLAIMSASSISTANCE  = "COMMISSION_TYPE_CLAIMS_ASSISTANCE"
)
