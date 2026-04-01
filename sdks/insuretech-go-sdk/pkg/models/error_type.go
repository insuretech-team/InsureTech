package models

// ErrorType represents a error_type
type ErrorType string

// ErrorType values
const (
	ErrorTypeERRORTYPEUNSPECIFIED ErrorType = "ERROR_TYPE_UNSPECIFIED"
	ErrorTypeERRORTYPEERROR  = "ERROR_TYPE_ERROR"
	ErrorTypeERRORTYPEWARNING  = "ERROR_TYPE_WARNING"
)
