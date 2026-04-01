package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	"google.golang.org/grpc"
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

// Generate handles POST /v1/documents:generate
//
// Google API style format selection (in order of precedence):
//  1. Body field:       {"output_format": "xlsx"}
//  2. Query parameter:  ?format=xlsx  or  ?output_format=xlsx
//
// Google API style media download ($alt=media):
//   ?$alt=media  or  ?alt=media  or body field {"alt":"media"}
//   → responds with raw file bytes + correct Content-Type + Content-Disposition
//   Default (no alt=media) → JSON metadata response with file_url
func (h *DocGenHandler) Generate(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Determine if caller wants raw bytes (alt=media Google API style).
	wantsMedia := q.Get("$alt") == "media" || q.Get("alt") == "media"

	// Determine format override from query params.
	formatFromQuery := strings.ToLower(strings.TrimSpace(q.Get("format")))
	if formatFromQuery == "" {
		formatFromQuery = strings.ToLower(strings.TrimSpace(q.Get("output_format")))
	}

	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req documentservicev1.GenerateDocumentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		// Query param takes precedence over body field (easier for client-side use).
		if formatFromQuery != "" {
			req.OutputFormat = formatFromQuery
		}
		// Propagate alt=media intent to the service so it can short-circuit storage.
		if wantsMedia && req.Alt == "" {
			req.Alt = "media"
		}
		return h.client.GenerateDocument(ctx, &req)
	})
}

// GenerateAndDownload handles the combined generate+download flow:
//   POST /v1/documents:generate?$alt=media
//
// Generates the document and immediately streams the raw file bytes back
// with the correct Content-Type and Content-Disposition header.
// This is a convenience alias; the standard Generate endpoint also supports
// ?$alt=media — this separate handler exists for explicit route matching.
func (h *DocGenHandler) GenerateAndDownload(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	formatFromQuery := strings.ToLower(strings.TrimSpace(q.Get("format")))
	if formatFromQuery == "" {
		formatFromQuery = strings.ToLower(strings.TrimSpace(q.Get("output_format")))
	}

	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req documentservicev1.GenerateDocumentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if formatFromQuery != "" {
			req.OutputFormat = formatFromQuery
		}
		req.Alt = "media"

		resp, err := h.client.GenerateDocument(ctx, &req)
		if err != nil {
			return nil, err
		}
		// If the service returned raw bytes directly in the response,
		// stream them immediately instead of going through JSON serialization.
		if resp != nil && resp.FileBytes != nil && len(resp.FileBytes) > 0 {
			ct := resp.ContentType
			if ct == "" {
				ct = contentTypeForFormat(req.OutputFormat)
			}
			filename := resp.Filename
			if filename == "" {
				filename = filenameForFormat(req.OutputFormat)
			}
			w.Header().Set("Content-Type", ct)
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
			w.Header().Set("Content-Length", strconv.Itoa(len(resp.FileBytes)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resp.FileBytes)
			return nil, nil // signal to callUnary that we've already written
		}
		return resp, nil
	})
}

// Download handles GET /v1/documents/{document_id}:download
// Streams raw file bytes with Content-Disposition: attachment.
func (h *DocGenHandler) DownloadRaw(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("document_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		resp, err := h.client.DownloadDocument(ctx, &documentservicev1.DownloadDocumentRequest{DocumentId: documentID})
		if err != nil {
			return nil, err
		}
		if resp != nil && len(resp.Content) > 0 {
			ct := resp.ContentType
			if ct == "" {
				ct = "application/octet-stream"
			}
			filename := resp.Filename
			if filename == "" {
				filename = "document"
			}
			w.Header().Set("Content-Type", ct)
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
			w.Header().Set("Content-Length", strconv.Itoa(len(resp.Content)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(resp.Content)
			return nil, nil
		}
		return resp, nil
	})
}

// contentTypeForFormat returns the MIME type for a given format string.
func contentTypeForFormat(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "pdf":
		return "application/pdf"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "html":
		return "text/html; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}

// filenameForFormat returns a sensible default filename for a given format string.
func filenameForFormat(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "pdf":
		return "document.pdf"
	case "docx":
		return "document.docx"
	case "xlsx":
		return "document.xlsx"
	case "html":
		return "document.html"
	default:
		return "document"
	}
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
		if err := protoUnmarshal(body, &req); err != nil {
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
		if err := protoUnmarshal(body, &req); err != nil {
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
			if err := protoUnmarshal(body, &req); err != nil {
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
