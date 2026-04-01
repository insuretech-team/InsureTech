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

// ListFraudAlerts List fraud alerts
func (s *FraudService) ListFraudAlerts(ctx context.Context) error {
	path := "/v1/fraud-alerts"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetFraudAlert Get fraud alert
func (s *FraudService) GetFraudAlert(ctx context.Context, fraudAlertId string) error {
	path := "/v1/fraud-alerts/{fraud_alert_id}"
	path = strings.ReplaceAll(path, "{fraud_alert_id}", fraudAlertId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateFraudCase Create fraud case
func (s *FraudService) CreateFraudCase(ctx context.Context, req *models.FraudFraudCaseCreationRequest) error {
	path := "/v1/fraud-cases"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetFraudCase Get fraud case
func (s *FraudService) GetFraudCase(ctx context.Context, fraudCaseId string) error {
	path := "/v1/fraud-cases/{fraud_case_id}"
	path = strings.ReplaceAll(path, "{fraud_case_id}", fraudCaseId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateFraudCase Update fraud case
func (s *FraudService) UpdateFraudCase(ctx context.Context, fraudCaseId string, req *models.FraudFraudCaseUpdateRequest) error {
	path := "/v1/fraud-cases/{fraud_case_id}"
	path = strings.ReplaceAll(path, "{fraud_case_id}", fraudCaseId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// CheckFraud Check for fraud
func (s *FraudService) CheckFraud(ctx context.Context, req *models.CheckFraudRequest) error {
	path := "/v1/fraud-checks"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListFraudRules List fraud rules
func (s *FraudService) ListFraudRules(ctx context.Context) error {
	path := "/v1/fraud-rules"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateFraudRule Create fraud rule
func (s *FraudService) CreateFraudRule(ctx context.Context, req *models.FraudFraudRuleCreationRequest) error {
	path := "/v1/fraud-rules"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdateFraudRule Update fraud rule
func (s *FraudService) UpdateFraudRule(ctx context.Context, ruleId string, req *models.FraudFraudRuleUpdateRequest) error {
	path := "/v1/fraud-rules/{rule_id}"
	path = strings.ReplaceAll(path, "{rule_id}", ruleId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// ActivateFraudRule Activate fraud rule
func (s *FraudService) ActivateFraudRule(ctx context.Context, ruleId string, req *models.FraudRuleActivationRequest) error {
	path := "/v1/fraud-rules/{rule_id}:activate"
	path = strings.ReplaceAll(path, "{rule_id}", ruleId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DeactivateFraudRule Deactivate fraud rule
func (s *FraudService) DeactivateFraudRule(ctx context.Context, ruleId string, req *models.FraudRuleDeactivationRequest) error {
	path := "/v1/fraud-rules/{rule_id}:deactivate"
	path = strings.ReplaceAll(path, "{rule_id}", ruleId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

