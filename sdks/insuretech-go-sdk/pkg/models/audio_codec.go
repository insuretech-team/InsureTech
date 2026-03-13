package models

// AudioCodec represents a audio_codec
type AudioCodec string

// AudioCodec values
const (
	AudioCodecAUDIOCODECUNSPECIFIED AudioCodec = "AUDIO_CODEC_UNSPECIFIED"
	AudioCodecAUDIOCODECOPUS  = "AUDIO_CODEC_OPUS"
	AudioCodecAUDIOCODECPCMU  = "AUDIO_CODEC_PCMU"
	AudioCodecAUDIOCODECPCMA  = "AUDIO_CODEC_PCMA"
	AudioCodecAUDIOCODECISAC  = "AUDIO_CODEC_ISAC"
)
