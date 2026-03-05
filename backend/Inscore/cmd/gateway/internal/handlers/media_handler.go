package handlers

import (
	"context"
	"net/http"
	"strconv"

	mediaservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MediaHandler proxies media APIs to the media gRPC service.
type MediaHandler struct {
	client MediaClient
}

// MediaClient keeps the handler decoupled from concrete gRPC transport.
type MediaClient interface {
	UploadMedia(ctx context.Context, in *mediaservicev1.UploadMediaRequest, opts ...grpc.CallOption) (*mediaservicev1.UploadMediaResponse, error)
	GetMedia(ctx context.Context, in *mediaservicev1.GetMediaRequest, opts ...grpc.CallOption) (*mediaservicev1.GetMediaResponse, error)
	ListMedia(ctx context.Context, in *mediaservicev1.ListMediaRequest, opts ...grpc.CallOption) (*mediaservicev1.ListMediaResponse, error)
	DownloadMedia(ctx context.Context, in *mediaservicev1.DownloadMediaRequest, opts ...grpc.CallOption) (*mediaservicev1.DownloadMediaResponse, error)
	DownloadOptimized(ctx context.Context, in *mediaservicev1.DownloadOptimizedRequest, opts ...grpc.CallOption) (*mediaservicev1.DownloadOptimizedResponse, error)
	DownloadThumbnail(ctx context.Context, in *mediaservicev1.DownloadThumbnailRequest, opts ...grpc.CallOption) (*mediaservicev1.DownloadThumbnailResponse, error)
	DeleteMedia(ctx context.Context, in *mediaservicev1.DeleteMediaRequest, opts ...grpc.CallOption) (*mediaservicev1.DeleteMediaResponse, error)
	ValidateMedia(ctx context.Context, in *mediaservicev1.ValidateMediaRequest, opts ...grpc.CallOption) (*mediaservicev1.ValidateMediaResponse, error)
	RequestProcessing(ctx context.Context, in *mediaservicev1.RequestProcessingRequest, opts ...grpc.CallOption) (*mediaservicev1.RequestProcessingResponse, error)
	GetProcessingJob(ctx context.Context, in *mediaservicev1.GetProcessingJobRequest, opts ...grpc.CallOption) (*mediaservicev1.GetProcessingJobResponse, error)
	ListProcessingJobs(ctx context.Context, in *mediaservicev1.ListProcessingJobsRequest, opts ...grpc.CallOption) (*mediaservicev1.ListProcessingJobsResponse, error)
}

// NewMediaHandler creates a MediaHandler from a gRPC connection.
func NewMediaHandler(conn *grpc.ClientConn) *MediaHandler {
	return &MediaHandler{client: mediaservicev1.NewMediaServiceClient(conn)}
}

// Upload handles POST /v1/media.
func (h *MediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req mediaservicev1.UploadMediaRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.TenantId == "" {
			req.TenantId = r.Header.Get("X-Tenant-ID")
		}
		return h.client.UploadMedia(ctx, &req)
	})
}

// Get handles GET /v1/media/{media_id}.
func (h *MediaHandler) Get(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetMedia(ctx, &mediaservicev1.GetMediaRequest{MediaId: mediaID})
	})
}

// List handles GET /v1/entities/{entity_type}/{entity_id}/media.
func (h *MediaHandler) List(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &mediaservicev1.ListMediaRequest{
			EntityType:       r.PathValue("entity_type"),
			EntityId:         r.PathValue("entity_id"),
			MediaType:        r.URL.Query().Get("media_type"),
			ValidationStatus: r.URL.Query().Get("validation_status"),
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
		return h.client.ListMedia(ctx, req)
	})
}

// Download handles GET /v1/media/{media_id}/download.
func (h *MediaHandler) Download(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DownloadMedia(ctx, &mediaservicev1.DownloadMediaRequest{MediaId: mediaID})
	})
}

// DownloadOptimized handles GET /v1/media/{media_id}/optimized.
func (h *MediaHandler) DownloadOptimized(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DownloadOptimized(ctx, &mediaservicev1.DownloadOptimizedRequest{MediaId: mediaID})
	})
}

// DownloadThumbnail handles GET /v1/media/{media_id}/thumbnail.
func (h *MediaHandler) DownloadThumbnail(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DownloadThumbnail(ctx, &mediaservicev1.DownloadThumbnailRequest{MediaId: mediaID})
	})
}

// Delete handles DELETE /v1/media/{media_id}.
func (h *MediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteMedia(ctx, &mediaservicev1.DeleteMediaRequest{MediaId: mediaID})
	})
}

// Validate handles POST /v1/media/{media_id}:validate.
func (h *MediaHandler) Validate(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req mediaservicev1.ValidateMediaRequest
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, &req); err != nil {
				return nil, err
			}
		}
		req.MediaId = mediaID
		return h.client.ValidateMedia(ctx, &req)
	})
}

// RequestProcessing handles POST /v1/media/{media_id}/process.
func (h *MediaHandler) RequestProcessing(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("media_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req mediaservicev1.RequestProcessingRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.MediaId = mediaID
		return h.client.RequestProcessing(ctx, &req)
	})
}

// GetProcessingJob handles GET /v1/processing-jobs/{job_id}.
func (h *MediaHandler) GetProcessingJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("job_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetProcessingJob(ctx, &mediaservicev1.GetProcessingJobRequest{JobId: jobID})
	})
}

// ListProcessingJobs handles GET /v1/processing-jobs.
func (h *MediaHandler) ListProcessingJobs(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &mediaservicev1.ListProcessingJobsRequest{
			MediaId:        r.URL.Query().Get("media_id"),
			ProcessingType: r.URL.Query().Get("processing_type"),
			Status:         r.URL.Query().Get("status"),
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
		return h.client.ListProcessingJobs(ctx, req)
	})
}
