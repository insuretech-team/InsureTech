package models


// ProcessingJobsListingRequest represents a processing_jobs_listing_request
type ProcessingJobsListingRequest struct {
	MediaId string `json:"media_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	ProcessingType string `json:"processing_type,omitempty"`
	Status string `json:"status,omitempty"`
}
