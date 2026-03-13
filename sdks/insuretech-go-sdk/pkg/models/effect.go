package models

// Effect represents a effect
type Effect string

// Effect values
const (
	EffectEFFECTUNSPECIFIED Effect = "EFFECT_UNSPECIFIED"
	EffectEFFECTALLOW  = "EFFECT_ALLOW"
	EffectEFFECTDENY  = "EFFECT_DENY"
)
