package models


// RoomConfig represents a room_config
type RoomConfig struct {
	RequireToken bool `json:"require_token,omitempty"`
	EnableRecording bool `json:"enable_recording,omitempty"`
	EnableTranscription bool `json:"enable_transcription,omitempty"`
	VideoConfig *VideoConfig `json:"video_config,omitempty"`
	AudioConfig *AudioConfig `json:"audio_config,omitempty"`
	SessionTimeoutSeconds int `json:"session_timeout_seconds,omitempty"`
	MaxParticipants int `json:"max_participants,omitempty"`
}
