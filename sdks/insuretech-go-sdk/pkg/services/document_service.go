package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// DocumentService handles document-related API calls
type DocumentService struct {
	Client Client
}

// ListDocumentTemplates List templates
func (s *DocumentService) ListDocumentTemplates(ctx context.Context) error {
	path := "/v1/document-templates"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateDocumentTemplate Create template
func (s *DocumentService) CreateDocumentTemplate(ctx context.Context, req *models.DocumentTemplateCreationRequest) error {
	path := "/v1/document-templates"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetDocumentTemplate Get template
func (s *DocumentService) GetDocumentTemplate(ctx context.Context, templateId string) error {
	path := "/v1/document-templates/{template_id}"
	path = strings.ReplaceAll(path, "{template_id}", templateId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateDocumentTemplate Update template
func (s *DocumentService) UpdateDocumentTemplate(ctx context.Context, templateId string, req *models.DocumentTemplateUpdateRequest) error {
	path := "/v1/document-templates/{template_id}"
	path = strings.ReplaceAll(path, "{template_id}", templateId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteDocumentTemplate Delete template
func (s *DocumentService) DeleteDocumentTemplate(ctx context.Context, templateId string) error {
	path := "/v1/document-templates/{template_id}"
	path = strings.ReplaceAll(path, "{template_id}", templateId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// DeactivateDocumentTemplate Deactivate template
func (s *DocumentService) DeactivateDocumentTemplate(ctx context.Context, templateId string, req *models.DocumentTemplateDeactivationRequest) error {
	path := "/v1/document-templates/{template_id}:deactivate"
	path = strings.ReplaceAll(path, "{template_id}", templateId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetDocument Get document
func (s *DocumentService) GetDocument(ctx context.Context, documentId string) error {
	path := "/v1/documents/{document_id}"
	path = strings.ReplaceAll(path, "{document_id}", documentId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// DeleteDocument Delete document
func (s *DocumentService) DeleteDocument(ctx context.Context, documentId string) error {
	path := "/v1/documents/{document_id}"
	path = strings.ReplaceAll(path, "{document_id}", documentId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// DownloadDocument Download document
func (s *DocumentService) DownloadDocument(ctx context.Context, documentId string) error {
	path := "/v1/documents/{document_id}/download"
	path = strings.ReplaceAll(path, "{document_id}", documentId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GenerateDocument Generate document
func (s *DocumentService) GenerateDocument(ctx context.Context, req *models.DocumentGenerationRequest) error {
	path := "/v1/documents:generate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListDocuments List documents for entity
func (s *DocumentService) ListDocuments(ctx context.Context, entityType string, entityId string) error {
	path := "/v1/entities/{entity_type}/{entity_id}/documents"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

