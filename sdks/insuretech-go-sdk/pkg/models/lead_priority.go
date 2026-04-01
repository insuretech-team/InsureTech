package models

// LeadPriority represents a lead_priority
type LeadPriority string

// LeadPriority values
const (
	LeadPriorityLEADPRIORITYUNSPECIFIED LeadPriority = "LEAD_PRIORITY_UNSPECIFIED"
	LeadPriorityLEADPRIORITYLOW  = "LEAD_PRIORITY_LOW"
	LeadPriorityLEADPRIORITYMEDIUM  = "LEAD_PRIORITY_MEDIUM"
	LeadPriorityLEADPRIORITYHIGH  = "LEAD_PRIORITY_HIGH"
	LeadPriorityLEADPRIORITYURGENT  = "LEAD_PRIORITY_URGENT"
)
