package models


// TrackSettings represents a track_settings
type TrackSettings struct {
	Bitrate int `json:"bitrate,omitempty"`
	Codec string `json:"codec,omitempty"`
	Width int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
	FrameRate float64 `json:"frame_rate,omitempty"`
	SampleRate int `json:"sample_rate,omitempty"`
	ChannelCount int `json:"channel_count,omitempty"`
}
