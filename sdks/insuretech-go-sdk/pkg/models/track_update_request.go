package models


// TrackUpdateRequest represents a track_update_request
type TrackUpdateRequest struct {
	TrackId string `json:"track_id"`
	Settings *TrackSettings `json:"settings,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
