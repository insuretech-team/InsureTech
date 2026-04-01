package models

import (
	"time"
)

// VoiceTranscript represents a voice_transcript
type VoiceTranscript struct {
	AuditInfo interface{} `json:"audit_info"`
	Confidence float64 `json:"confidence,omitempty"`
	Id string `json:"id"`
	Language string `json:"language"`
	SequenceNumber int `json:"sequence_number"`
	Speaker *SpeakerType `json:"speaker"`
	Text string `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	VoiceSessionId string `json:"voice_session_id"`
}
