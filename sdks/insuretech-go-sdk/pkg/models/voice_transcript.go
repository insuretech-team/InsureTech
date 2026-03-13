package models

import (
	"time"
)

// VoiceTranscript represents a voice_transcript
type VoiceTranscript struct {
	VoiceSessionId string `json:"voice_session_id"`
	Language string `json:"language"`
	Confidence float64 `json:"confidence,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	Speaker *SpeakerType `json:"speaker"`
	Text string `json:"text"`
	SequenceNumber int `json:"sequence_number"`
	Timestamp time.Time `json:"timestamp"`
}
