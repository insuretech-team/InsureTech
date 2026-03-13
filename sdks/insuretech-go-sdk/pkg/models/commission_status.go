package models

// CommissionStatus represents a commission_status
type CommissionStatus string

// CommissionStatus values
const (
	CommissionStatusCOMMISSIONSTATUSUNSPECIFIED CommissionStatus = "COMMISSION_STATUS_UNSPECIFIED"
	CommissionStatusCOMMISSIONSTATUSPENDING  = "COMMISSION_STATUS_PENDING"
	CommissionStatusCOMMISSIONSTATUSPROCESSING  = "COMMISSION_STATUS_PROCESSING"
	CommissionStatusCOMMISSIONSTATUSPAID  = "COMMISSION_STATUS_PAID"
	CommissionStatusCOMMISSIONSTATUSCANCELLED  = "COMMISSION_STATUS_CANCELLED"
)
