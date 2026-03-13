package models

// PartnerAgentStatus represents a partner_agent_status
type PartnerAgentStatus string

// PartnerAgentStatus values
const (
	PartnerAgentStatusAGENTSTATUSUNSPECIFIED PartnerAgentStatus = "AGENT_STATUS_UNSPECIFIED"
	PartnerAgentStatusAGENTSTATUSACTIVE  = "AGENT_STATUS_ACTIVE"
	PartnerAgentStatusAGENTSTATUSINACTIVE  = "AGENT_STATUS_INACTIVE"
	PartnerAgentStatusAGENTSTATUSSUSPENDED  = "AGENT_STATUS_SUSPENDED"
)
