package models


// ComplianceLogCreationResponse represents a compliance_log_creation_response
type ComplianceLogCreationResponse struct {
	ComplianceLogId string `json:"compliance_log_id,omitempty"`
	Error *Error `json:"error,omitempty"`
}
