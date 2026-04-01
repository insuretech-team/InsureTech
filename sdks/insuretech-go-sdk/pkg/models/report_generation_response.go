package models


// ReportGenerationResponse represents a report_generation_response
type ReportGenerationResponse struct {
	FileName string `json:"file_name,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	ReportUrl string `json:"report_url,omitempty"`
}
