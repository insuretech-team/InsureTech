package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// UnderwritingService handles underwriting-related API calls
type UnderwritingService struct {
	Client Client
}

// GetQuote Get quote details
func (s *UnderwritingService) GetQuote(ctx context.Context, quoteId string) (*models.UnderwritingQuoteRetrievalResponse, error) {
	path := "/v1/quotes/{quote_id}"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	var result models.UnderwritingQuoteRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ApproveUnderwriting Approve underwriting (manual)
func (s *UnderwritingService) ApproveUnderwriting(ctx context.Context, quoteId string, req *models.UnderwritingApprovalRequest) (*models.UnderwritingApprovalResponse, error) {
	path := "/v1/quotes/{quote_id}"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	var result models.UnderwritingApprovalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUnderwritingDecision Get underwriting decision
func (s *UnderwritingService) GetUnderwritingDecision(ctx context.Context, quoteId string) (*models.UnderwritingUnderwritingDecisionRetrievalResponse, error) {
	path := "/v1/quotes/{quote_id}/decision"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	var result models.UnderwritingUnderwritingDecisionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SubmitHealthDeclaration Submit health declaration
func (s *UnderwritingService) SubmitHealthDeclaration(ctx context.Context, quoteId string, req *models.HealthDeclarationSubmissionRequest) (*models.HealthDeclarationSubmissionResponse, error) {
	path := "/v1/quotes/{quote_id}/health-declaration"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	var result models.HealthDeclarationSubmissionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetHealthDeclaration Get health declaration
func (s *UnderwritingService) GetHealthDeclaration(ctx context.Context, quoteId string) (*models.UnderwritingHealthDeclarationRetrievalResponse, error) {
	path := "/v1/quotes/{quote_id}/health-declaration"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	var result models.UnderwritingHealthDeclarationRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListQuotes List quotes for beneficiary
func (s *UnderwritingService) ListQuotes(ctx context.Context, beneficiaryId string) (*models.UnderwritingQuotesListingResponse, error) {
	path := "/v1/beneficiaries/{beneficiary_id}/quotes"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	var result models.UnderwritingQuotesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestQuote Request premium quote
func (s *UnderwritingService) RequestQuote(ctx context.Context, req *models.RequestQuoteRequest) (*models.RequestQuoteResponse, error) {
	path := "/v1/quotes"
	var result models.RequestQuoteResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

