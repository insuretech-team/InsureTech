package models


// ReportGenerationResponse represents a report_generation_response
type ReportGenerationResponse struct {
	ReportUrl string `json:"report_url,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	Error *Error `json:"error,omitempty"`
}
