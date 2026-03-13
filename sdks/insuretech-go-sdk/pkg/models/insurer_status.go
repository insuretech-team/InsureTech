package models

// InsurerStatus represents a insurer_status
type InsurerStatus string

// InsurerStatus values
const (
	InsurerStatusINSURERSTATUSUNSPECIFIED InsurerStatus = "INSURER_STATUS_UNSPECIFIED"
	InsurerStatusINSURERSTATUSACTIVE  = "INSURER_STATUS_ACTIVE"
	InsurerStatusINSURERSTATUSSUSPENDED  = "INSURER_STATUS_SUSPENDED"
	InsurerStatusINSURERSTATUSINACTIVE  = "INSURER_STATUS_INACTIVE"
)
