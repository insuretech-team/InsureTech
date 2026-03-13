package models

// CasePriority represents a case_priority
type CasePriority string

// CasePriority values
const (
	CasePriorityCASEPRIORITYUNSPECIFIED CasePriority = "CASE_PRIORITY_UNSPECIFIED"
	CasePriorityCASEPRIORITYLOW  = "CASE_PRIORITY_LOW"
	CasePriorityCASEPRIORITYMEDIUM  = "CASE_PRIORITY_MEDIUM"
	CasePriorityCASEPRIORITYHIGH  = "CASE_PRIORITY_HIGH"
	CasePriorityCASEPRIORITYURGENT  = "CASE_PRIORITY_URGENT"
)
