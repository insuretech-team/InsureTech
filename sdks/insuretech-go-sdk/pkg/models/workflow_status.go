package models

// WorkflowStatus represents a workflow_status
type WorkflowStatus string

// WorkflowStatus values
const (
	WorkflowStatusWORKFLOWSTATUSUNSPECIFIED WorkflowStatus = "WORKFLOW_STATUS_UNSPECIFIED"
	WorkflowStatusWORKFLOWSTATUSDRAFT  = "WORKFLOW_STATUS_DRAFT"
	WorkflowStatusWORKFLOWSTATUSACTIVE  = "WORKFLOW_STATUS_ACTIVE"
	WorkflowStatusWORKFLOWSTATUSINACTIVE  = "WORKFLOW_STATUS_INACTIVE"
	WorkflowStatusWORKFLOWSTATUSARCHIVED  = "WORKFLOW_STATUS_ARCHIVED"
)
