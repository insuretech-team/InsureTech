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

// ListQuotes List quotes for beneficiary
func (s *UnderwritingService) ListQuotes(ctx context.Context, beneficiaryId string) error {
	path := "/v1/beneficiaries/{beneficiary_id}/quotes"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RequestQuote Request premium quote
func (s *UnderwritingService) RequestQuote(ctx context.Context, req *models.RequestQuoteRequest) error {
	path := "/v1/quotes"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetQuote Get quote details
func (s *UnderwritingService) GetQuote(ctx context.Context, quoteId string) error {
	path := "/v1/quotes/{quote_id}"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetUnderwritingDecision Get underwriting decision
func (s *UnderwritingService) GetUnderwritingDecision(ctx context.Context, quoteId string) error {
	path := "/v1/quotes/{quote_id}/decision"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetHealthDeclaration Get health declaration
func (s *UnderwritingService) GetHealthDeclaration(ctx context.Context, quoteId string) error {
	path := "/v1/quotes/{quote_id}/health-declaration"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// SubmitHealthDeclaration Submit health declaration
func (s *UnderwritingService) SubmitHealthDeclaration(ctx context.Context, quoteId string, req *models.HealthDeclarationSubmissionRequest) error {
	path := "/v1/quotes/{quote_id}/health-declaration"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ApproveUnderwriting Approve underwriting (manual)
func (s *UnderwritingService) ApproveUnderwriting(ctx context.Context, quoteId string, req *models.UnderwritingApprovalRequest) error {
	path := "/v1/quotes/{quote_id}:approve"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ConvertQuoteToPolicy Convert quote to policy
func (s *UnderwritingService) ConvertQuoteToPolicy(ctx context.Context, quoteId string, req *models.UnderwritingConvertQuoteToPolicyRequest) error {
	path := "/v1/quotes/{quote_id}:convert"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RejectUnderwriting Reject underwriting (manual)
func (s *UnderwritingService) RejectUnderwriting(ctx context.Context, quoteId string, req *models.UnderwritingRejectionRequest) error {
	path := "/v1/quotes/{quote_id}:reject"
	path = strings.ReplaceAll(path, "{quote_id}", quoteId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

