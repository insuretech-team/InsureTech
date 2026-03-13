package models

// ReportStatus represents a report_status
type ReportStatus string

// ReportStatus values
const (
	ReportStatusREPORTSTATUSUNSPECIFIED ReportStatus = "REPORT_STATUS_UNSPECIFIED"
	ReportStatusREPORTSTATUSGENERATING  = "REPORT_STATUS_GENERATING"
	ReportStatusREPORTSTATUSCOMPLETED  = "REPORT_STATUS_COMPLETED"
	ReportStatusREPORTSTATUSFAILED  = "REPORT_STATUS_FAILED"
)
