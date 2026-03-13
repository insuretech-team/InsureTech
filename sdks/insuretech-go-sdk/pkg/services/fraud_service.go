package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// FraudService handles fraud-related API calls
type FraudService struct {
	Client Client
}

// DeactivateFraudRule Deactivate fraud rule
func (s *FraudService) DeactivateFraudRule(ctx context.Context, ruleId string, req *models.FraudRuleDeactivationRequest) (*models.FraudRuleDeactivationResponse, error) {
	path := "/v1/fraud-rules/{rule_id}:deactivate"
	path = strings.ReplaceAll(path, "{rule_id}", ruleId)
	var result models.FraudRuleDeactivationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateFraudCase Create fraud case
func (s *FraudService) CreateFraudCase(ctx context.Context, req *models.FraudFraudCaseCreationRequest) (*models.FraudFraudCaseCreationResponse, error) {
	path := "/v1/fraud-cases"
	var result models.FraudFraudCaseCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFraudRule Update fraud rule
func (s *FraudService) UpdateFraudRule(ctx context.Context, ruleId string, req *models.FraudFraudRuleUpdateRequest) (*models.FraudFraudRuleUpdateResponse, error) {
	path := "/v1/fraud-rules/{rule_id}"
	path = strings.ReplaceAll(path, "{rule_id}", ruleId)
	var result models.FraudFraudRuleUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFraudCase Get fraud case
func (s *FraudService) GetFraudCase(ctx context.Context, fraudCaseId string) (*models.FraudFraudCaseRetrievalResponse, error) {
	path := "/v1/fraud-cases/{fraud_case_id}"
	path = strings.ReplaceAll(path, "{fraud_case_id}", fraudCaseId)
	var result models.FraudFraudCaseRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFraudCase Update fraud case
func (s *FraudService) UpdateFraudCase(ctx context.Context, fraudCaseId string, req *models.FraudFraudCaseUpdateRequest) (*models.FraudFraudCaseUpdateResponse, error) {
	path := "/v1/fraud-cases/{fraud_case_id}"
	path = strings.ReplaceAll(path, "{fraud_case_id}", fraudCaseId)
	var result models.FraudFraudCaseUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListFraudRules List fraud rules
func (s *FraudService) ListFraudRules(ctx context.Context) (*models.FraudFraudRulesListingResponse, error) {
	path := "/v1/fraud-rules"
	var result models.FraudFraudRulesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateFraudRule Create fraud rule
func (s *FraudService) CreateFraudRule(ctx context.Context, req *models.FraudFraudRuleCreationRequest) (*models.FraudFraudRuleCreationResponse, error) {
	path := "/v1/fraud-rules"
	var result models.FraudFraudRuleCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ActivateFraudRule Activate fraud rule
func (s *FraudService) ActivateFraudRule(ctx context.Context, ruleId string, req *models.FraudRuleActivationRequest) (*models.FraudRuleActivationResponse, error) {
	path := "/v1/fraud-rules/{rule_id}:activate"
	path = strings.ReplaceAll(path, "{rule_id}", ruleId)
	var result models.FraudRuleActivationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CheckFraud Check for fraud
func (s *FraudService) CheckFraud(ctx context.Context, req *models.CheckFraudRequest) (*models.CheckFraudResponse, error) {
	path := "/v1/fraud-checks"
	var result models.CheckFraudResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFraudAlert Get fraud alert
func (s *FraudService) GetFraudAlert(ctx context.Context, fraudAlertId string) (*models.FraudFraudAlertRetrievalResponse, error) {
	path := "/v1/fraud-alerts/{fraud_alert_id}"
	path = strings.ReplaceAll(path, "{fraud_alert_id}", fraudAlertId)
	var result models.FraudFraudAlertRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListFraudAlerts List fraud alerts
func (s *FraudService) ListFraudAlerts(ctx context.Context) (*models.FraudAlertsListingResponse, error) {
	path := "/v1/fraud-alerts"
	var result models.FraudAlertsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

