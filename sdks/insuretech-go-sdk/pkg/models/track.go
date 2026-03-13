package models


// Track represents a track
type Track struct {
	State interface{} `json:"state"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Settings interface{} `json:"settings"`
	TrackId string `json:"track_id"`
	PeerId string `json:"peer_id"`
	Type *TrackType `json:"type"`
	Label string `json:"label,omitempty"`
	Muted bool `json:"muted"`
}
