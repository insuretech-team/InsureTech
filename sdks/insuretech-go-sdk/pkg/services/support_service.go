package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// SupportService handles support-related API calls
type SupportService struct {
	Client Client
}

// GetTicket Get ticket
func (s *SupportService) GetTicket(ctx context.Context, ticketId string) (*models.TicketRetrievalResponse, error) {
	path := "/v1/tickets/{ticket_id}"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	var result models.TicketRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AssignTicket Assign ticket
func (s *SupportService) AssignTicket(ctx context.Context, ticketId string, req *models.TicketAssignmentRequest) (*models.TicketAssignmentResponse, error) {
	path := "/v1/tickets/{ticket_id}"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	var result models.TicketAssignmentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateKnowledgeBaseArticle Update Knowledge Base Article
func (s *SupportService) UpdateKnowledgeBaseArticle(ctx context.Context, articleId string, req *models.KnowledgeBaseArticleUpdateRequest) (*models.KnowledgeBaseArticleUpdateResponse, error) {
	path := "/v1/knowledge-base/{article_id}"
	path = strings.ReplaceAll(path, "{article_id}", articleId)
	var result models.KnowledgeBaseArticleUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteKnowledgeBaseArticle Delete Knowledge Base Article
func (s *SupportService) DeleteKnowledgeBaseArticle(ctx context.Context, articleId string) error {
	path := "/v1/knowledge-base/{article_id}"
	path = strings.ReplaceAll(path, "{article_id}", articleId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// CreateTicket Create ticket
func (s *SupportService) CreateTicket(ctx context.Context, req *models.TicketCreationRequest) (*models.TicketCreationResponse, error) {
	path := "/v1/tickets"
	var result models.TicketCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListTickets List tickets
func (s *SupportService) ListTickets(ctx context.Context) (*models.TicketsListingResponse, error) {
	path := "/v1/tickets"
	var result models.TicketsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchKnowledgeBase Search knowledge base
func (s *SupportService) SearchKnowledgeBase(ctx context.Context) (*models.KnowledgeBaseSearchResponse, error) {
	path := "/v1/knowledge-base/search"
	var result models.KnowledgeBaseSearchResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateKnowledgeBaseArticle Create Knowledge Base Article
func (s *SupportService) CreateKnowledgeBaseArticle(ctx context.Context, req *models.KnowledgeBaseArticleCreationRequest) (*models.KnowledgeBaseArticleCreationResponse, error) {
	path := "/v1/knowledge-base"
	var result models.KnowledgeBaseArticleCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListFAQs List FAQs
func (s *SupportService) ListFAQs(ctx context.Context) (*models.FAQsListingResponse, error) {
	path := "/v1/faqs"
	var result models.FAQsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateFAQ Create FAQ
func (s *SupportService) CreateFAQ(ctx context.Context, req *models.FAQCreationRequest) (*models.FAQCreationResponse, error) {
	path := "/v1/faqs"
	var result models.FAQCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddTicketMessage Add ticket message
func (s *SupportService) AddTicketMessage(ctx context.Context, ticketId string, req *models.AddTicketMessageRequest) (*models.AddTicketMessageResponse, error) {
	path := "/v1/tickets/{ticket_id}/messages"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	var result models.AddTicketMessageResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKnowledgeBaseArticle Get knowledge base article
func (s *SupportService) GetKnowledgeBaseArticle(ctx context.Context, slug string) (*models.KnowledgeBaseArticleRetrievalResponse, error) {
	path := "/v1/knowledge-base/{slug}"
	path = strings.ReplaceAll(path, "{slug}", slug)
	var result models.KnowledgeBaseArticleRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFAQ Update FAQ
func (s *SupportService) UpdateFAQ(ctx context.Context, faqId string, req *models.FAQUpdateRequest) (*models.FAQUpdateResponse, error) {
	path := "/v1/faqs/{faq_id}"
	path = strings.ReplaceAll(path, "{faq_id}", faqId)
	var result models.FAQUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFAQ Delete FAQ
func (s *SupportService) DeleteFAQ(ctx context.Context, faqId string) error {
	path := "/v1/faqs/{faq_id}"
	path = strings.ReplaceAll(path, "{faq_id}", faqId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// UpdateTicketStatus Update ticket status
func (s *SupportService) UpdateTicketStatus(ctx context.Context, ticketId string, req *models.TicketStatusUpdateRequest) (*models.TicketStatusUpdateResponse, error) {
	path := "/v1/tickets/{ticket_id}/status"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	var result models.TicketStatusUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

