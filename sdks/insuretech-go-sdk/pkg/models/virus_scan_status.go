package models

// VirusScanStatus represents a virus_scan_status
type VirusScanStatus string

// VirusScanStatus values
const (
	VirusScanStatusVIRUSSCANSTATUSUNSPECIFIED VirusScanStatus = "VIRUS_SCAN_STATUS_UNSPECIFIED"
	VirusScanStatusVIRUSSCANSTATUSPENDING  = "VIRUS_SCAN_STATUS_PENDING"
	VirusScanStatusVIRUSSCANSTATUSCLEAN  = "VIRUS_SCAN_STATUS_CLEAN"
	VirusScanStatusVIRUSSCANSTATUSINFECTED  = "VIRUS_SCAN_STATUS_INFECTED"
)
