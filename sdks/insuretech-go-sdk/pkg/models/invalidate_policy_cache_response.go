package models


// InvalidatePolicyCacheResponse represents a invalidate_policy_cache_response
type InvalidatePolicyCacheResponse struct {
	Invalidated bool `json:"invalidated,omitempty"`
	Error *Error `json:"error,omitempty"`
}
