package models


// RenewalRenewalScheduleRetrievalResponse represents a renewal_renewal_schedule_retrieval_response
type RenewalRenewalScheduleRetrievalResponse struct {
	Reminders []*RenewalReminder `json:"reminders,omitempty"`
	RenewalSchedule *RenewalSchedule `json:"renewal_schedule,omitempty"`
}
