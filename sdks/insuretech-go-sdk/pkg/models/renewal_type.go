package models

// RenewalType represents a renewal_type
type RenewalType string

// RenewalType values
const (
	RenewalTypeRENEWALTYPEUNSPECIFIED RenewalType = "RENEWAL_TYPE_UNSPECIFIED"
	RenewalTypeRENEWALTYPEAUTOMATIC  = "RENEWAL_TYPE_AUTOMATIC"
	RenewalTypeRENEWALTYPEMANUAL  = "RENEWAL_TYPE_MANUAL"
)
