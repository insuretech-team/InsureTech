package handlers

import (
	"context"
	"net/http"
	"strconv"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DocumentHandler proxies document upload/retrieval to the storage gRPC service.
type DocumentHandler struct {
	client StorageClient
}

// StorageClient keeps the handler decoupled from concrete gRPC transport.
type StorageClient interface {
	UploadFile(ctx context.Context, in *storageservicev1.UploadFileRequest, opts ...grpc.CallOption) (*storageservicev1.UploadFileResponse, error)
	UploadFiles(ctx context.Context, in *storageservicev1.UploadFilesRequest, opts ...grpc.CallOption) (*storageservicev1.UploadFilesResponse, error)
	GetUploadURL(ctx context.Context, in *storageservicev1.GetUploadURLRequest, opts ...grpc.CallOption) (*storageservicev1.GetUploadURLResponse, error)
	FinalizeUpload(ctx context.Context, in *storageservicev1.FinalizeUploadRequest, opts ...grpc.CallOption) (*storageservicev1.FinalizeUploadResponse, error)
	GetFile(ctx context.Context, in *storageservicev1.GetFileRequest, opts ...grpc.CallOption) (*storageservicev1.GetFileResponse, error)
	UpdateFile(ctx context.Context, in *storageservicev1.UpdateFileRequest, opts ...grpc.CallOption) (*storageservicev1.UpdateFileResponse, error)
	GetDownloadURL(ctx context.Context, in *storageservicev1.GetDownloadURLRequest, opts ...grpc.CallOption) (*storageservicev1.GetDownloadURLResponse, error)
	DeleteFile(ctx context.Context, in *storageservicev1.DeleteFileRequest, opts ...grpc.CallOption) (*storageservicev1.DeleteFileResponse, error)
	ListFiles(ctx context.Context, in *storageservicev1.ListFilesRequest, opts ...grpc.CallOption) (*storageservicev1.ListFilesResponse, error)
}

// NewDocumentHandler creates a DocumentHandler from a gRPC connection.
func NewDocumentHandler(conn *grpc.ClientConn) *DocumentHandler {
	return &DocumentHandler{client: storageservicev1.NewStorageServiceClient(conn)}
}

// NewDocumentHandlerWithClient creates a DocumentHandler from an abstract client.
func NewDocumentHandlerWithClient(client StorageClient) *DocumentHandler {
	return &DocumentHandler{client: client}
}

// Upload handles file upload requests.
// POST /v1/storage/files
func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req storageservicev1.UploadFileRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		return h.client.UploadFile(ctx, &req)
	})
}

// UploadBatch handles multi-file upload requests.
// POST /v1/storage/files:batch
func (h *DocumentHandler) UploadBatch(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req storageservicev1.UploadFilesRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		return h.client.UploadFiles(ctx, &req)
	})
}

// GetUploadURL creates a presigned upload URL.
// POST /v1/storage/files:upload-url
func (h *DocumentHandler) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req storageservicev1.GetUploadURLRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		return h.client.GetUploadURL(ctx, &req)
	})
}

// FinalizeUpload finalizes metadata after direct upload.
// POST /v1/storage/files:finalize
func (h *DocumentHandler) FinalizeUpload(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req storageservicev1.FinalizeUploadRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		return h.client.FinalizeUpload(ctx, &req)
	})
}

// Get retrieves file metadata by ID.
// GET /v1/storage/files/{id}
func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		tenantID := r.Header.Get("X-Tenant-ID")
		return h.client.GetFile(ctx, &storageservicev1.GetFileRequest{
			TenantId: tenantID,
			FileId:   id,
		})
	})
}

// Update partially updates file metadata.
// PATCH /v1/storage/files/{id}
func (h *DocumentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req storageservicev1.UpdateFileRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.FileId == "" {
			req.FileId = id
		}
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		return h.client.UpdateFile(ctx, &req)
	})
}

// GetDownloadURL returns a presigned download URL.
// POST /v1/storage/files/{id}:download-url
func (h *DocumentHandler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req storageservicev1.GetDownloadURLRequest
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, &req); err != nil {
				return nil, err
			}
		}
		req.FileId = id
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		if q := r.URL.Query().Get("expires_in_minutes"); q != "" && req.ExpiresInMinutes == 0 {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.ExpiresInMinutes = int32(n)
			}
		}
		return h.client.GetDownloadURL(ctx, &req)
	})
}

// Delete removes file metadata and object.
// DELETE /v1/storage/files/{id}
func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteFile(ctx, &storageservicev1.DeleteFileRequest{
			TenantId: r.Header.Get("X-Tenant-ID"),
			FileId:   id,
		})
	})
}

// List returns filtered storage files.
// GET /v1/storage/files
func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &storageservicev1.ListFilesRequest{
			TenantId:      r.Header.Get("X-Tenant-ID"),
			ReferenceId:   r.URL.Query().Get("reference_id"),
			ReferenceType: r.URL.Query().Get("reference_type"),
			UploadedBy:    r.URL.Query().Get("uploaded_by"),
		}
		if req.UploadedBy == "" {
			req.UploadedBy = r.URL.Query().Get("user_id")
		}
		if q := r.URL.Query().Get("file_type"); q != "" {
			if n, err := strconv.Atoi(q); err == nil {
				req.FileType = storageentityv1.FileType(n)
			}
		}
		var page, pageSize int32
		if q := r.URL.Query().Get("page"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				page = int32(n)
			}
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				pageSize = int32(n)
			}
		}
		if page > 0 || pageSize > 0 {
			if page == 0 {
				page = 1
			}
			req.Page = &commonv1.PaginationRequest{
				Page:     page,
				PageSize: pageSize,
			}
		}
		return h.client.ListFiles(ctx, req)
	})
}
