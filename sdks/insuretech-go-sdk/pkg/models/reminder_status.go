package models

// ReminderStatus represents a reminder_status
type ReminderStatus string

// ReminderStatus values
const (
	ReminderStatusREMINDERSTATUSUNSPECIFIED ReminderStatus = "REMINDER_STATUS_UNSPECIFIED"
	ReminderStatusREMINDERSTATUSPENDING  = "REMINDER_STATUS_PENDING"
	ReminderStatusREMINDERSTATUSSENT  = "REMINDER_STATUS_SENT"
	ReminderStatusREMINDERSTATUSFAILED  = "REMINDER_STATUS_FAILED"
	ReminderStatusREMINDERSTATUSCANCELLED  = "REMINDER_STATUS_CANCELLED"
)
