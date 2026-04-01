package models


// StreamStatsRequest represents a stream_stats_request
type StreamStatsRequest struct {
	IntervalSeconds int `json:"interval_seconds,omitempty"`
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
}
