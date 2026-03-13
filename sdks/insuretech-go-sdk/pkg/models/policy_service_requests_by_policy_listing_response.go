package models


// PolicyServiceRequestsByPolicyListingResponse represents a policy_service_requests_by_policy_listing_response
type PolicyServiceRequestsByPolicyListingResponse struct {
	Requests []*PolicyServiceRequest `json:"requests,omitempty"`
}
