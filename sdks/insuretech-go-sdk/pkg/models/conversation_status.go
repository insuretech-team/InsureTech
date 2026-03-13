package models

// ConversationStatus represents a conversation_status
type ConversationStatus string

// ConversationStatus values
const (
	ConversationStatusCONVERSATIONSTATUSUNSPECIFIED ConversationStatus = "CONVERSATION_STATUS_UNSPECIFIED"
	ConversationStatusCONVERSATIONSTATUSACTIVE  = "CONVERSATION_STATUS_ACTIVE"
	ConversationStatusCONVERSATIONSTATUSRESOLVED  = "CONVERSATION_STATUS_RESOLVED"
	ConversationStatusCONVERSATIONSTATUSESCALATED  = "CONVERSATION_STATUS_ESCALATED"
	ConversationStatusCONVERSATIONSTATUSCLOSED  = "CONVERSATION_STATUS_CLOSED"
)
