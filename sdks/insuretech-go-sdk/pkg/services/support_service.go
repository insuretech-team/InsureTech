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

// ListFAQs List FAQs
func (s *SupportService) ListFAQs(ctx context.Context) error {
	path := "/v1/faqs"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateFAQ Create FAQ
func (s *SupportService) CreateFAQ(ctx context.Context, req *models.FAQCreationRequest) error {
	path := "/v1/faqs"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdateFAQ Update FAQ
func (s *SupportService) UpdateFAQ(ctx context.Context, faqId string, req *models.FAQUpdateRequest) error {
	path := "/v1/faqs/{faq_id}"
	path = strings.ReplaceAll(path, "{faq_id}", faqId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteFAQ Delete FAQ
func (s *SupportService) DeleteFAQ(ctx context.Context, faqId string) error {
	path := "/v1/faqs/{faq_id}"
	path = strings.ReplaceAll(path, "{faq_id}", faqId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// CreateKnowledgeBaseArticle Create Knowledge Base Article
func (s *SupportService) CreateKnowledgeBaseArticle(ctx context.Context, req *models.KnowledgeBaseArticleCreationRequest) error {
	path := "/v1/knowledge-base"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SearchKnowledgeBase Search knowledge base
func (s *SupportService) SearchKnowledgeBase(ctx context.Context) error {
	path := "/v1/knowledge-base/search"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateKnowledgeBaseArticle Update Knowledge Base Article
func (s *SupportService) UpdateKnowledgeBaseArticle(ctx context.Context, articleId string, req *models.KnowledgeBaseArticleUpdateRequest) error {
	path := "/v1/knowledge-base/{article_id}"
	path = strings.ReplaceAll(path, "{article_id}", articleId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteKnowledgeBaseArticle Delete Knowledge Base Article
func (s *SupportService) DeleteKnowledgeBaseArticle(ctx context.Context, articleId string) error {
	path := "/v1/knowledge-base/{article_id}"
	path = strings.ReplaceAll(path, "{article_id}", articleId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetKnowledgeBaseArticle Get knowledge base article
func (s *SupportService) GetKnowledgeBaseArticle(ctx context.Context, slug string) error {
	path := "/v1/knowledge-base/{slug}"
	path = strings.ReplaceAll(path, "{slug}", slug)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListTickets List tickets
func (s *SupportService) ListTickets(ctx context.Context) error {
	path := "/v1/tickets"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateTicket Create ticket
func (s *SupportService) CreateTicket(ctx context.Context, req *models.TicketCreationRequest) error {
	path := "/v1/tickets"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetTicket Get ticket
func (s *SupportService) GetTicket(ctx context.Context, ticketId string) error {
	path := "/v1/tickets/{ticket_id}"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// AddTicketMessage Add ticket message
func (s *SupportService) AddTicketMessage(ctx context.Context, ticketId string, req *models.AddTicketMessageRequest) error {
	path := "/v1/tickets/{ticket_id}/messages"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdateTicketStatus Update ticket status
func (s *SupportService) UpdateTicketStatus(ctx context.Context, ticketId string, req *models.TicketStatusUpdateRequest) error {
	path := "/v1/tickets/{ticket_id}/status"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// AssignTicket Assign ticket
func (s *SupportService) AssignTicket(ctx context.Context, ticketId string, req *models.TicketAssignmentRequest) error {
	path := "/v1/tickets/{ticket_id}:assign"
	path = strings.ReplaceAll(path, "{ticket_id}", ticketId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

