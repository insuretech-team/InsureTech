package models


// KYCStep represents a kyc_step
type KYCStep struct {
	ChallengeType string `json:"challenge_type,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Instruction string `json:"instruction,omitempty"`
	InstructionKey string `json:"instruction_key,omitempty"`
	State string `json:"state,omitempty"`
	StepNumber int `json:"step_number,omitempty"`
	TimeoutSeconds int `json:"timeout_seconds,omitempty"`
}
