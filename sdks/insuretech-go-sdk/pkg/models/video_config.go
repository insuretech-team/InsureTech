package models


// VideoConfig represents a video_config
type VideoConfig struct {
	AllowedQualities []*VideoQuality `json:"allowed_qualities,omitempty"`
	DefaultQuality *VideoQuality `json:"default_quality,omitempty"`
	EnableSimulcast bool `json:"enable_simulcast,omitempty"`
	MaxBitrateKbps int `json:"max_bitrate_kbps,omitempty"`
	MinBitrateKbps int `json:"min_bitrate_kbps,omitempty"`
}
