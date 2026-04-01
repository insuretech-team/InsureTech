package models


// TrackSettings represents a track_settings
type TrackSettings struct {
	Bitrate int `json:"bitrate,omitempty"`
	ChannelCount int `json:"channel_count,omitempty"`
	Codec string `json:"codec,omitempty"`
	FrameRate float64 `json:"frame_rate,omitempty"`
	Height int `json:"height,omitempty"`
	SampleRate int `json:"sample_rate,omitempty"`
	Width int `json:"width,omitempty"`
}
