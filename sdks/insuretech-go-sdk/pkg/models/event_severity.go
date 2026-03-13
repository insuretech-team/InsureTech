package models

// EventSeverity represents a event_severity
type EventSeverity string

// EventSeverity values
const (
	EventSeverityEVENTSEVERITYUNSPECIFIED EventSeverity = "EVENT_SEVERITY_UNSPECIFIED"
	EventSeverityEVENTSEVERITYINFO  = "EVENT_SEVERITY_INFO"
	EventSeverityEVENTSEVERITYWARNING  = "EVENT_SEVERITY_WARNING"
	EventSeverityEVENTSEVERITYERROR  = "EVENT_SEVERITY_ERROR"
	EventSeverityEVENTSEVERITYCRITICAL  = "EVENT_SEVERITY_CRITICAL"
)
