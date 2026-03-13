package models


// ParticipantAnalytics represents a participant_analytics
type ParticipantAnalytics struct {
	PeerId string `json:"peer_id,omitempty"`
	TotalTimeSeconds string `json:"total_time_seconds,omitempty"`
}
