package models

// DecisionMethod represents a decision_method
type DecisionMethod string

// DecisionMethod values
const (
	DecisionMethodDECISIONMETHODUNSPECIFIED DecisionMethod = "DECISION_METHOD_UNSPECIFIED"
	DecisionMethodDECISIONMETHODAUTOMATIC  = "DECISION_METHOD_AUTOMATIC"
	DecisionMethodDECISIONMETHODMANUAL  = "DECISION_METHOD_MANUAL"
	DecisionMethodDECISIONMETHODHYBRID  = "DECISION_METHOD_HYBRID"
)
