package models

// SessionType represents a session_type
type SessionType string

// SessionType values
const (
	SessionTypeSESSIONTYPEUNSPECIFIED SessionType = "SESSION_TYPE_UNSPECIFIED"
	SessionTypeSESSIONTYPESERVERSIDE  = "SESSION_TYPE_SERVER_SIDE"
	SessionTypeSESSIONTYPEJWT  = "SESSION_TYPE_JWT"
)
