package models


// UnderwritingUnderwritingDecisionRetrievalResponse represents a underwriting_underwriting_decision_retrieval_response
type UnderwritingUnderwritingDecisionRetrievalResponse struct {
	Decision *UnderwritingDecision `json:"decision,omitempty"`
	Error *Error `json:"error,omitempty"`
}
