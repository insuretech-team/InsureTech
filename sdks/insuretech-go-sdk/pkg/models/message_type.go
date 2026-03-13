package models

// MessageType represents a message_type
type MessageType string

// MessageType values
const (
	MessageTypeMESSAGETYPEUNSPECIFIED MessageType = "MESSAGE_TYPE_UNSPECIFIED"
	MessageTypeMESSAGETYPEREPLY  = "MESSAGE_TYPE_REPLY"
	MessageTypeMESSAGETYPENOTE  = "MESSAGE_TYPE_NOTE"
	MessageTypeMESSAGETYPESTATUSCHANGE  = "MESSAGE_TYPE_STATUS_CHANGE"
	MessageTypeMESSAGETYPEASSIGNMENT  = "MESSAGE_TYPE_ASSIGNMENT"
)
