package models

// NotificationPriority represents a notification_priority
type NotificationPriority string

// NotificationPriority values
const (
	NotificationPriorityNOTIFICATIONPRIORITYUNSPECIFIED NotificationPriority = "NOTIFICATION_PRIORITY_UNSPECIFIED"
	NotificationPriorityNOTIFICATIONPRIORITYLOW  = "NOTIFICATION_PRIORITY_LOW"
	NotificationPriorityNOTIFICATIONPRIORITYNORMAL  = "NOTIFICATION_PRIORITY_NORMAL"
	NotificationPriorityNOTIFICATIONPRIORITYHIGH  = "NOTIFICATION_PRIORITY_HIGH"
	NotificationPriorityNOTIFICATIONPRIORITYURGENT  = "NOTIFICATION_PRIORITY_URGENT"
)
