package models


// UnderwritingDecisionRetrievalResponse represents a underwriting_decision_retrieval_response
type UnderwritingDecisionRetrievalResponse struct {
	Decision *UnderwritingDecision `json:"decision,omitempty"`
	Error *Error `json:"error,omitempty"`
}
