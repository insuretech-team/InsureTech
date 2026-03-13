package models


// ProcessingJobsListingResponse represents a processing_jobs_listing_response
type ProcessingJobsListingResponse struct {
	Jobs []*ingJobProcessing `json:"jobs,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
