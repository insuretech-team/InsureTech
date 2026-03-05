package server

import (
	"context"
	"errors"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/service"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	"google.golang.org/grpc/metadata"
)

// DocumentHandler implements the DocumentService gRPC server.
type DocumentHandler struct {
	documentservicev1.UnimplementedDocumentServiceServer
	docService *service.DocumentService
}

func NewDocumentHandler(docService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{docService: docService}
}

func actorFromContext(ctx context.Context, fallback string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fallback
	}
	for _, k := range []string{"x-user-id", "x-actor-id", "x-sub", "user-id"} {
		vals := md.Get(k)
		if len(vals) > 0 {
			v := strings.TrimSpace(vals[0])
			if v != "" {
				return v
			}
		}
	}
	return fallback
}

func tenantFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	for _, k := range []string{"x-tenant-id", "tenant-id", "x-tenant"} {
		vals := md.Get(k)
		if len(vals) > 0 {
			v := strings.TrimSpace(vals[0])
			if v != "" {
				return v
			}
		}
	}
	return ""
}

func mapErr(err error, notFoundMessage string) *commonv1.Error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return &commonv1.Error{Code: "INVALID_ARGUMENT", Message: err.Error()}
	case errors.Is(err, service.ErrTemplateNotFound), errors.Is(err, service.ErrDocumentNotFound):
		return &commonv1.Error{Code: "NOT_FOUND", Message: notFoundMessage}
	case errors.Is(err, service.ErrUnsupportedOutput):
		return &commonv1.Error{Code: "UNIMPLEMENTED", Message: err.Error()}
	default:
		return &commonv1.Error{Code: "INTERNAL", Message: err.Error()}
	}
}

func (h *DocumentHandler) GenerateDocument(ctx context.Context, req *documentservicev1.GenerateDocumentRequest) (*documentservicev1.GenerateDocumentResponse, error) {
	tenantID := tenantFromContext(ctx)
	generatedBy := actorFromContext(ctx, tenantID)
	doc, err := h.docService.GenerateDocument(ctx, req.TemplateId, req.EntityType, req.EntityId, req.Data, req.IncludeQrCode, tenantID, generatedBy)
	if err != nil {
		return &documentservicev1.GenerateDocumentResponse{Error: mapErr(err, "template or document context not found")}, nil
	}
	return &documentservicev1.GenerateDocumentResponse{
		DocumentId: doc.Id,
		FileUrl:    doc.FileUrl,
		Message:    "Document generated successfully",
	}, nil
}

func (h *DocumentHandler) GetDocument(ctx context.Context, req *documentservicev1.GetDocumentRequest) (*documentservicev1.GetDocumentResponse, error) {
	doc, err := h.docService.GetDocument(ctx, req.DocumentId)
	if err != nil {
		return &documentservicev1.GetDocumentResponse{Error: mapErr(err, "document not found")}, nil
	}
	return &documentservicev1.GetDocumentResponse{Document: doc}, nil
}

func (h *DocumentHandler) ListDocuments(ctx context.Context, req *documentservicev1.ListDocumentsRequest) (*documentservicev1.ListDocumentsResponse, error) {
	docs, total, err := h.docService.ListDocuments(ctx, req.EntityType, req.EntityId, req.Status, int(req.Page), int(req.PageSize))
	if err != nil {
		return &documentservicev1.ListDocumentsResponse{Error: mapErr(err, "documents not found")}, nil
	}
	return &documentservicev1.ListDocumentsResponse{Documents: docs, TotalCount: int32(total)}, nil
}

func (h *DocumentHandler) DownloadDocument(ctx context.Context, req *documentservicev1.DownloadDocumentRequest) (*documentservicev1.DownloadDocumentResponse, error) {
	content, contentType, filename, err := h.docService.DownloadDocument(ctx, req.DocumentId)
	if err != nil {
		return &documentservicev1.DownloadDocumentResponse{Error: mapErr(err, "document not found")}, nil
	}
	return &documentservicev1.DownloadDocumentResponse{Content: content, ContentType: contentType, Filename: filename}, nil
}

func (h *DocumentHandler) DeleteDocument(ctx context.Context, req *documentservicev1.DeleteDocumentRequest) (*documentservicev1.DeleteDocumentResponse, error) {
	tenantID := tenantFromContext(ctx)
	if err := h.docService.DeleteDocument(ctx, req.DocumentId, tenantID); err != nil {
		return &documentservicev1.DeleteDocumentResponse{Error: mapErr(err, "document not found")}, nil
	}
	return &documentservicev1.DeleteDocumentResponse{Message: "Document deleted successfully"}, nil
}

func (h *DocumentHandler) CreateDocumentTemplate(ctx context.Context, req *documentservicev1.CreateDocumentTemplateRequest) (*documentservicev1.CreateDocumentTemplateResponse, error) {
	actor := actorFromContext(ctx, "")
	templateID, err := h.docService.CreateTemplate(ctx, req.Name, req.Type, req.Description, req.TemplateContent, req.OutputFormat, req.Variables, actor)
	if err != nil {
		return &documentservicev1.CreateDocumentTemplateResponse{Error: mapErr(err, "failed to create template")}, nil
	}
	return &documentservicev1.CreateDocumentTemplateResponse{TemplateId: templateID, Message: "Template created successfully"}, nil
}

func (h *DocumentHandler) GetDocumentTemplate(ctx context.Context, req *documentservicev1.GetDocumentTemplateRequest) (*documentservicev1.GetDocumentTemplateResponse, error) {
	tpl, err := h.docService.GetTemplate(ctx, req.TemplateId)
	if err != nil {
		return &documentservicev1.GetDocumentTemplateResponse{Error: mapErr(err, "template not found")}, nil
	}
	return &documentservicev1.GetDocumentTemplateResponse{Template: tpl}, nil
}

func (h *DocumentHandler) ListDocumentTemplates(ctx context.Context, req *documentservicev1.ListDocumentTemplatesRequest) (*documentservicev1.ListDocumentTemplatesResponse, error) {
	items, next, total, err := h.docService.ListTemplates(ctx, req.Type, req.ActiveOnly, int(req.PageSize), req.PageToken)
	if err != nil {
		return &documentservicev1.ListDocumentTemplatesResponse{Error: mapErr(err, "failed to list templates")}, nil
	}
	return &documentservicev1.ListDocumentTemplatesResponse{Templates: items, NextPageToken: next, TotalCount: int32(total)}, nil
}

func (h *DocumentHandler) UpdateDocumentTemplate(ctx context.Context, req *documentservicev1.UpdateDocumentTemplateRequest) (*documentservicev1.UpdateDocumentTemplateResponse, error) {
	if err := h.docService.UpdateTemplate(ctx, req.TemplateId, req.Template); err != nil {
		return &documentservicev1.UpdateDocumentTemplateResponse{Error: mapErr(err, "template not found")}, nil
	}
	return &documentservicev1.UpdateDocumentTemplateResponse{Message: "Template updated successfully"}, nil
}

func (h *DocumentHandler) DeactivateDocumentTemplate(ctx context.Context, req *documentservicev1.DeactivateDocumentTemplateRequest) (*documentservicev1.DeactivateDocumentTemplateResponse, error) {
	if err := h.docService.DeactivateTemplate(ctx, req.TemplateId); err != nil {
		return &documentservicev1.DeactivateDocumentTemplateResponse{Error: mapErr(err, "template not found")}, nil
	}
	return &documentservicev1.DeactivateDocumentTemplateResponse{Message: "Template deactivated successfully"}, nil
}

func (h *DocumentHandler) DeleteDocumentTemplate(ctx context.Context, req *documentservicev1.DeleteDocumentTemplateRequest) (*documentservicev1.DeleteDocumentTemplateResponse, error) {
	if err := h.docService.DeleteTemplate(ctx, req.TemplateId); err != nil {
		return &documentservicev1.DeleteDocumentTemplateResponse{Error: mapErr(err, "template not found")}, nil
	}
	return &documentservicev1.DeleteDocumentTemplateResponse{Message: "Template deleted successfully"}, nil
}
