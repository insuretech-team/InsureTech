package models


// Track represents a track
type Track struct {
	Label string `json:"label,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Muted bool `json:"muted"`
	PeerId string `json:"peer_id"`
	Settings interface{} `json:"settings"`
	State interface{} `json:"state"`
	TrackId string `json:"track_id"`
	Type *TrackType `json:"type"`
}
