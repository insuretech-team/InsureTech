package models

// FieldBehavior represents a field_behavior
type FieldBehavior string

// FieldBehavior values
const (
	FieldBehaviorFIELDBEHAVIORUNSPECIFIED FieldBehavior = "FIELD_BEHAVIOR_UNSPECIFIED"
	FieldBehaviorOPTIONAL  = "OPTIONAL"
	FieldBehaviorREQUIRED  = "REQUIRED"
	FieldBehaviorOUTPUTONLY  = "OUTPUT_ONLY"
	FieldBehaviorINPUTONLY  = "INPUT_ONLY"
	FieldBehaviorIMMUTABLE  = "IMMUTABLE"
	FieldBehaviorUNORDEREDLIST  = "UNORDERED_LIST"
	FieldBehaviorNONEMPTYDEFAULT  = "NON_EMPTY_DEFAULT"
	FieldBehaviorIDENTIFIER  = "IDENTIFIER"
)
