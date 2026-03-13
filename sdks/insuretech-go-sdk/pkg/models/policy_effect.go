package models

// PolicyEffect represents a policy_effect
type PolicyEffect string

// PolicyEffect values
const (
	PolicyEffectPOLICYEFFECTUNSPECIFIED PolicyEffect = "POLICY_EFFECT_UNSPECIFIED"
	PolicyEffectPOLICYEFFECTALLOW  = "POLICY_EFFECT_ALLOW"
	PolicyEffectPOLICYEFFECTDENY  = "POLICY_EFFECT_DENY"
)
