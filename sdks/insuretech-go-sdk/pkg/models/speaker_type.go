package models

// SpeakerType represents a speaker_type
type SpeakerType string

// SpeakerType values
const (
	SpeakerTypeSPEAKERTYPEUNSPECIFIED SpeakerType = "SPEAKER_TYPE_UNSPECIFIED"
	SpeakerTypeSPEAKERTYPEUSER  = "SPEAKER_TYPE_USER"
	SpeakerTypeSPEAKERTYPESYSTEM  = "SPEAKER_TYPE_SYSTEM"
	SpeakerTypeSPEAKERTYPEAGENT  = "SPEAKER_TYPE_AGENT"
)
