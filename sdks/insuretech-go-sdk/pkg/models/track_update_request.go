package models


// TrackUpdateRequest represents a track_update_request
type TrackUpdateRequest struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Settings *TrackSettings `json:"settings,omitempty"`
	TrackId string `json:"track_id"`
}
