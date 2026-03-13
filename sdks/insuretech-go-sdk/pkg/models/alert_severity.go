package models

// AlertSeverity represents a alert_severity
type AlertSeverity string

// AlertSeverity values
const (
	AlertSeverityALERTSEVERITYUNSPECIFIED AlertSeverity = "ALERT_SEVERITY_UNSPECIFIED"
	AlertSeverityALERTSEVERITYINFO  = "ALERT_SEVERITY_INFO"
	AlertSeverityALERTSEVERITYWARNING  = "ALERT_SEVERITY_WARNING"
	AlertSeverityALERTSEVERITYERROR  = "ALERT_SEVERITY_ERROR"
	AlertSeverityALERTSEVERITYCRITICAL  = "ALERT_SEVERITY_CRITICAL"
)
