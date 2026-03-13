package models

// ReminderChannel represents a reminder_channel
type ReminderChannel string

// ReminderChannel values
const (
	ReminderChannelREMINDERCHANNELUNSPECIFIED ReminderChannel = "REMINDER_CHANNEL_UNSPECIFIED"
	ReminderChannelREMINDERCHANNELEMAIL  = "REMINDER_CHANNEL_EMAIL"
	ReminderChannelREMINDERCHANNELSMS  = "REMINDER_CHANNEL_SMS"
	ReminderChannelREMINDERCHANNELPUSH  = "REMINDER_CHANNEL_PUSH"
	ReminderChannelREMINDERCHANNELVOICE  = "REMINDER_CHANNEL_VOICE"
)
