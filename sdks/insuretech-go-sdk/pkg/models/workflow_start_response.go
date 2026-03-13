package models


// WorkflowStartResponse represents a workflow_start_response
type WorkflowStartResponse struct {
	WorkflowInstanceId string `json:"workflow_instance_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
