package models

// AIAgentStatus represents a ai_agent_status
type AIAgentStatus string

// AIAgentStatus values
const (
	AIAgentStatusAGENTSTATUSUNSPECIFIED AIAgentStatus = "AGENT_STATUS_UNSPECIFIED"
	AIAgentStatusAGENTSTATUSIDLE  = "AGENT_STATUS_IDLE"
	AIAgentStatusAGENTSTATUSPROCESSING  = "AGENT_STATUS_PROCESSING"
	AIAgentStatusAGENTSTATUSOFFLINE  = "AGENT_STATUS_OFFLINE"
)
