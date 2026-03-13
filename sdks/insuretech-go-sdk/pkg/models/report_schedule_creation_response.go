package models


// ReportScheduleCreationResponse represents a report_schedule_creation_response
type ReportScheduleCreationResponse struct {
	ReportScheduleId string `json:"report_schedule_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
