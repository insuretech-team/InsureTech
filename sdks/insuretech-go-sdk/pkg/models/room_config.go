package models


// RoomConfig represents a room_config
type RoomConfig struct {
	AudioConfig *AudioConfig `json:"audio_config,omitempty"`
	EnableRecording bool `json:"enable_recording,omitempty"`
	EnableTranscription bool `json:"enable_transcription,omitempty"`
	MaxParticipants int `json:"max_participants,omitempty"`
	RequireToken bool `json:"require_token,omitempty"`
	SessionTimeoutSeconds int `json:"session_timeout_seconds,omitempty"`
	VideoConfig *VideoConfig `json:"video_config,omitempty"`
}
