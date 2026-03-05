package handlers

import (
	"context"
	"net/http"
	"strconv"

	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DocGenHandler proxies docgen APIs to the document gRPC service.
type DocGenHandler struct {
	client DocGenClient
}

// DocGenClient keeps the handler decoupled from concrete gRPC transport.
type DocGenClient interface {
	GenerateDocument(ctx context.Context, in *documentservicev1.GenerateDocumentRequest, opts ...grpc.CallOption) (*documentservicev1.GenerateDocumentResponse, error)
	GetDocument(ctx context.Context, in *documentservicev1.GetDocumentRequest, opts ...grpc.CallOption) (*documentservicev1.GetDocumentResponse, error)
	ListDocuments(ctx context.Context, in *documentservicev1.ListDocumentsRequest, opts ...grpc.CallOption) (*documentservicev1.ListDocumentsResponse, error)
	DownloadDocument(ctx context.Context, in *documentservicev1.DownloadDocumentRequest, opts ...grpc.CallOption) (*documentservicev1.DownloadDocumentResponse, error)
	DeleteDocument(ctx context.Context, in *documentservicev1.DeleteDocumentRequest, opts ...grpc.CallOption) (*documentservicev1.DeleteDocumentResponse, error)
	CreateDocumentTemplate(ctx context.Context, in *documentservicev1.CreateDocumentTemplateRequest, opts ...grpc.CallOption) (*documentservicev1.CreateDocumentTemplateResponse, error)
	GetDocumentTemplate(ctx context.Context, in *documentservicev1.GetDocumentTemplateRequest, opts ...grpc.CallOption) (*documentservicev1.GetDocumentTemplateResponse, error)
	ListDocumentTemplates(ctx context.Context, in *documentservicev1.ListDocumentTemplatesRequest, opts ...grpc.CallOption) (*documentservicev1.ListDocumentTemplatesResponse, error)
	UpdateDocumentTemplate(ctx context.Context, in *documentservicev1.UpdateDocumentTemplateRequest, opts ...grpc.CallOption) (*documentservicev1.UpdateDocumentTemplateResponse, error)
	DeactivateDocumentTemplate(ctx context.Context, in *documentservicev1.DeactivateDocumentTemplateRequest, opts ...grpc.CallOption) (*documentservicev1.DeactivateDocumentTemplateResponse, error)
	DeleteDocumentTemplate(ctx context.Context, in *documentservicev1.DeleteDocumentTemplateRequest, opts ...grpc.CallOption) (*documentservicev1.DeleteDocumentTemplateResponse, error)
}

// NewDocGenHandler creates a DocGenHandler from a gRPC connection.
func NewDocGenHandler(conn *grpc.ClientConn) *DocGenHandler {
	return &DocGenHandler{client: documentservicev1.NewDocumentServiceClient(conn)}
}

// Generate handles POST /v1/documents.
func (h *DocGenHandler) Generate(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req documentservicev1.GenerateDocumentRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.GenerateDocument(ctx, &req)
	})
}

// GetDocument handles GET /v1/documents/{document_id}.
func (h *DocGenHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("document_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetDocument(ctx, &documentservicev1.GetDocumentRequest{DocumentId: documentID})
	})
}

// ListDocuments handles GET /v1/entities/{entity_type}/{entity_id}/documents.
func (h *DocGenHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &documentservicev1.ListDocumentsRequest{
			EntityType: r.PathValue("entity_type"),
			EntityId:   r.PathValue("entity_id"),
			Status:     r.URL.Query().Get("status"),
		}
		if q := r.URL.Query().Get("page"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.Page = int32(n)
			}
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListDocuments(ctx, req)
	})
}

// Download handles GET /v1/documents/{document_id}/download.
func (h *DocGenHandler) Download(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("document_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DownloadDocument(ctx, &documentservicev1.DownloadDocumentRequest{DocumentId: documentID})
	})
}

// DeleteDocument handles DELETE /v1/documents/{document_id}.
func (h *DocGenHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("document_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteDocument(ctx, &documentservicev1.DeleteDocumentRequest{DocumentId: documentID})
	})
}

// CreateTemplate handles POST /v1/document-templates.
func (h *DocGenHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req documentservicev1.CreateDocumentTemplateRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateDocumentTemplate(ctx, &req)
	})
}

// GetTemplate handles GET /v1/document-templates/{template_id}.
func (h *DocGenHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetDocumentTemplate(ctx, &documentservicev1.GetDocumentTemplateRequest{TemplateId: templateID})
	})
}

// ListTemplates handles GET /v1/document-templates.
func (h *DocGenHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &documentservicev1.ListDocumentTemplatesRequest{Type: r.URL.Query().Get("type")}
		if q := r.URL.Query().Get("active_only"); q == "true" || q == "1" {
			req.ActiveOnly = true
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		req.PageToken = r.URL.Query().Get("page_token")
		return h.client.ListDocumentTemplates(ctx, req)
	})
}

// UpdateTemplate handles PATCH /v1/document-templates/{template_id}.
func (h *DocGenHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req documentservicev1.UpdateDocumentTemplateRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.TemplateId = templateID
		return h.client.UpdateDocumentTemplate(ctx, &req)
	})
}

// DeactivateTemplate handles POST /v1/document-templates/{template_id}:deactivate.
func (h *DocGenHandler) DeactivateTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req documentservicev1.DeactivateDocumentTemplateRequest
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, &req); err != nil {
				return nil, err
			}
		}
		req.TemplateId = templateID
		return h.client.DeactivateDocumentTemplate(ctx, &req)
	})
}

// DeleteTemplate handles DELETE /v1/document-templates/{template_id}.
func (h *DocGenHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("template_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteDocumentTemplate(ctx, &documentservicev1.DeleteDocumentTemplateRequest{TemplateId: templateID})
	})
}
