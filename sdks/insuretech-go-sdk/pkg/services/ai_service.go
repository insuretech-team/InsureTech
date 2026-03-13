package services

import (
	"context"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// AiService handles ai-related API calls
type AiService struct {
	Client Client
}

// Chat Chat with AI agent
func (s *AiService) Chat(ctx context.Context, req *models.ChatRequest) (*models.ChatResponse, error) {
	path := "/v1/ai/chat"
	var result models.ChatResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EvaluateClaim Evaluate claim
func (s *AiService) EvaluateClaim(ctx context.Context, req *models.ClaimEvaluationRequest) (*models.ClaimEvaluationResponse, error) {
	path := "/v1/ai/claims:evaluate"
	var result models.ClaimEvaluationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AssessRisk Assess risk
func (s *AiService) AssessRisk(ctx context.Context, req *models.RiskAssessmentRequest) (*models.RiskAssessmentResponse, error) {
	path := "/v1/ai/risk:assess"
	var result models.RiskAssessmentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AnalyzeDocument Analyze document
func (s *AiService) AnalyzeDocument(ctx context.Context, req *models.DocumentAnalysisRequest) (*models.DocumentAnalysisResponse, error) {
	path := "/v1/ai/documents:analyze"
	var result models.DocumentAnalysisResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DetectFraud Detect fraud
func (s *AiService) DetectFraud(ctx context.Context, req *models.DetectFraudRequest) (*models.DetectFraudResponse, error) {
	path := "/v1/ai/fraud:detect"
	var result models.DetectFraudResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

