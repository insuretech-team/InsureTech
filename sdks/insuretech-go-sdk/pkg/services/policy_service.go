package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// PolicyService handles policy-related API calls
type PolicyService struct {
	Client Client
}

// RenewPolicy Renew policy
func (s *PolicyService) RenewPolicy(ctx context.Context, policyId string, req *models.PolicyPolicyRenewalRequest) (*models.PolicyPolicyRenewalResponse, error) {
	path := "/v1/policies/{policy_id}:renew"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyPolicyRenewalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserPolicies List user policies
func (s *PolicyService) ListUserPolicies(ctx context.Context, customerId string) (*models.UserPoliciesListingResponse, error) {
	path := "/v1/users/{customer_id}/policies"
	path = strings.ReplaceAll(path, "{customer_id}", customerId)
	var result models.UserPoliciesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IssuePolicy Issue policy (explicit transition from pending to active after payment)
func (s *PolicyService) IssuePolicy(ctx context.Context, policyId string, req *models.PolicyIssuanceRequest) (*models.PolicyIssuanceResponse, error) {
	path := "/v1/policies/{policy_id}:issue"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyIssuanceResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelPolicy Cancel policy
func (s *PolicyService) CancelPolicy(ctx context.Context, policyId string, req *models.PolicyCancellationRequest) (*models.PolicyCancellationResponse, error) {
	path := "/v1/policies/{policy_id}:cancel"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyCancellationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePolicy Create policy
func (s *PolicyService) CreatePolicy(ctx context.Context, req *models.PolicyPolicyCreationRequest) (*models.PolicyPolicyCreationResponse, error) {
	path := "/v1/policies"
	var result models.PolicyPolicyCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPolicy Get policy details
func (s *PolicyService) GetPolicy(ctx context.Context, policyId string) (*models.PolicyPolicyRetrievalResponse, error) {
	path := "/v1/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyPolicyRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePolicy Update policy
func (s *PolicyService) UpdatePolicy(ctx context.Context, policyId string, req *models.PolicyPolicyUpdateRequest) (*models.PolicyPolicyUpdateResponse, error) {
	path := "/v1/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyPolicyUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GeneratePolicyDocument Generate policy document
func (s *PolicyService) GeneratePolicyDocument(ctx context.Context, policyId string, req *models.PolicyDocumentGenerationRequest) (*models.PolicyDocumentGenerationResponse, error) {
	path := "/v1/policies/{policy_id}:generate-document"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyDocumentGenerationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

