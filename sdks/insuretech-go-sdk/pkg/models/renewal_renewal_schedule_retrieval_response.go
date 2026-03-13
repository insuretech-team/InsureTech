package models


// RenewalRenewalScheduleRetrievalResponse represents a renewal_renewal_schedule_retrieval_response
type RenewalRenewalScheduleRetrievalResponse struct {
	RenewalSchedule *RenewalSchedule `json:"renewal_schedule,omitempty"`
	Reminders []*RenewalReminder `json:"reminders,omitempty"`
	Error *Error `json:"error,omitempty"`
}
