package models

// InsuranceAgentStatus represents a insurance_agent_status
type InsuranceAgentStatus string

// InsuranceAgentStatus values
const (
	InsuranceAgentStatusAGENTSTATUSUNSPECIFIED InsuranceAgentStatus = "AGENT_STATUS_UNSPECIFIED"
	InsuranceAgentStatusAGENTSTATUSACTIVE  = "AGENT_STATUS_ACTIVE"
	InsuranceAgentStatusAGENTSTATUSINACTIVE  = "AGENT_STATUS_INACTIVE"
	InsuranceAgentStatusAGENTSTATUSSUSPENDED  = "AGENT_STATUS_SUSPENDED"
)
