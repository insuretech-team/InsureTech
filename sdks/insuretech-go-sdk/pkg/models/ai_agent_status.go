package models

// AiAgentStatus represents a ai_agent_status
type AiAgentStatus string

// AiAgentStatus values
const (
	AiAgentStatusAGENTSTATUSUNSPECIFIED AiAgentStatus = "AGENT_STATUS_UNSPECIFIED"
	AiAgentStatusAGENTSTATUSIDLE  = "AGENT_STATUS_IDLE"
	AiAgentStatusAGENTSTATUSPROCESSING  = "AGENT_STATUS_PROCESSING"
	AiAgentStatusAGENTSTATUSOFFLINE  = "AGENT_STATUS_OFFLINE"
)
