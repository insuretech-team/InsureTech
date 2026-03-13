package models


// ProcessingJobRetrievalResponse represents a processing_job_retrieval_response
type ProcessingJobRetrievalResponse struct {
	Job *ingJobProcessing `json:"job,omitempty"`
	Error *Error `json:"error,omitempty"`
}
