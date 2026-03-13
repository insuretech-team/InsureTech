package models

// ExecutionStatus represents a execution_status
type ExecutionStatus string

// ExecutionStatus values
const (
	ExecutionStatusEXECUTIONSTATUSUNSPECIFIED ExecutionStatus = "EXECUTION_STATUS_UNSPECIFIED"
	ExecutionStatusEXECUTIONSTATUSPENDING  = "EXECUTION_STATUS_PENDING"
	ExecutionStatusEXECUTIONSTATUSRUNNING  = "EXECUTION_STATUS_RUNNING"
	ExecutionStatusEXECUTIONSTATUSCOMPLETED  = "EXECUTION_STATUS_COMPLETED"
	ExecutionStatusEXECUTIONSTATUSFAILED  = "EXECUTION_STATUS_FAILED"
)
