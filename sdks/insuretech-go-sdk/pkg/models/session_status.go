package models

// SessionStatus represents a session_status
type SessionStatus string

// SessionStatus values
const (
	SessionStatusSESSIONSTATUSUNSPECIFIED SessionStatus = "SESSION_STATUS_UNSPECIFIED"
	SessionStatusSESSIONSTATUSACTIVE  = "SESSION_STATUS_ACTIVE"
	SessionStatusSESSIONSTATUSCOMPLETED  = "SESSION_STATUS_COMPLETED"
	SessionStatusSESSIONSTATUSFAILED  = "SESSION_STATUS_FAILED"
	SessionStatusSESSIONSTATUSTIMEOUT  = "SESSION_STATUS_TIMEOUT"
)
