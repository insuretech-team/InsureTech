package models

// MessageRole represents a message_role
type MessageRole string

// MessageRole values
const (
	MessageRoleMESSAGEROLEUNSPECIFIED MessageRole = "MESSAGE_ROLE_UNSPECIFIED"
	MessageRoleMESSAGEROLEUSER  = "MESSAGE_ROLE_USER"
	MessageRoleMESSAGEROLEASSISTANT  = "MESSAGE_ROLE_ASSISTANT"
	MessageRoleMESSAGEROLESYSTEM  = "MESSAGE_ROLE_SYSTEM"
)
