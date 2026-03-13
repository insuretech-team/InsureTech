package models

// NotificationChannel represents a notification_channel
type NotificationChannel string

// NotificationChannel values
const (
	NotificationChannelNOTIFICATIONCHANNELUNSPECIFIED NotificationChannel = "NOTIFICATION_CHANNEL_UNSPECIFIED"
	NotificationChannelNOTIFICATIONCHANNELSMS  = "NOTIFICATION_CHANNEL_SMS"
	NotificationChannelNOTIFICATIONCHANNELEMAIL  = "NOTIFICATION_CHANNEL_EMAIL"
	NotificationChannelNOTIFICATIONCHANNELPUSH  = "NOTIFICATION_CHANNEL_PUSH"
	NotificationChannelNOTIFICATIONCHANNELWHATSAPP  = "NOTIFICATION_CHANNEL_WHATSAPP"
	NotificationChannelNOTIFICATIONCHANNELINAPP  = "NOTIFICATION_CHANNEL_IN_APP"
)
