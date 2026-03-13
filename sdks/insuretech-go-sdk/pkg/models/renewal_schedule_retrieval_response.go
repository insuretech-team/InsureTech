package models


// RenewalScheduleRetrievalResponse represents a renewal_schedule_retrieval_response
type RenewalScheduleRetrievalResponse struct {
	RenewalSchedule *RenewalSchedule `json:"renewal_schedule,omitempty"`
	Reminders []*RenewalReminder `json:"reminders,omitempty"`
	Error *Error `json:"error,omitempty"`
}
