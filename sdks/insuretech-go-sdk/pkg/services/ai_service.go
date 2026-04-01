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
func (s *AiService) Chat(ctx context.Context, req *models.ChatRequest) error {
	path := "/v1/ai/chat"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// EvaluateClaim Evaluate claim
func (s *AiService) EvaluateClaim(ctx context.Context, req *models.ClaimEvaluationRequest) error {
	path := "/v1/ai/claims:evaluate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// AnalyzeDocument Analyze document
func (s *AiService) AnalyzeDocument(ctx context.Context, req *models.DocumentAnalysisRequest) error {
	path := "/v1/ai/documents:analyze"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DetectFraud Detect fraud
func (s *AiService) DetectFraud(ctx context.Context, req *models.DetectFraudRequest) error {
	path := "/v1/ai/fraud:detect"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// AssessRisk Assess risk
func (s *AiService) AssessRisk(ctx context.Context, req *models.RiskAssessmentRequest) error {
	path := "/v1/ai/risk:assess"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

