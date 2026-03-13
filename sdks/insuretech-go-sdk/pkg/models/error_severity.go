package models

// ErrorSeverity represents a error_severity
type ErrorSeverity string

// ErrorSeverity values
const (
	ErrorSeverityERRORSEVERITYUNSPECIFIED ErrorSeverity = "ERROR_SEVERITY_UNSPECIFIED"
	ErrorSeverityERRORSEVERITYINFO  = "ERROR_SEVERITY_INFO"
	ErrorSeverityERRORSEVERITYWARNING  = "ERROR_SEVERITY_WARNING"
	ErrorSeverityERRORSEVERITYERROR  = "ERROR_SEVERITY_ERROR"
	ErrorSeverityERRORSEVERITYCRITICAL  = "ERROR_SEVERITY_CRITICAL"
)
