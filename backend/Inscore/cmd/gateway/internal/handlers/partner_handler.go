package handlers

import (
	"context"
	"net/http"
	"strconv"

	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// PartnerHandler proxies partner APIs to the partner gRPC service.
type PartnerHandler struct {
	client PartnerClient
}

// PartnerClient abstracts partner gRPC client methods used by gateway.
type PartnerClient interface {
	CreatePartner(ctx context.Context, in *partnerservicev1.CreatePartnerRequest, opts ...grpc.CallOption) (*partnerservicev1.CreatePartnerResponse, error)
	GetPartner(ctx context.Context, in *partnerservicev1.GetPartnerRequest, opts ...grpc.CallOption) (*partnerservicev1.GetPartnerResponse, error)
	UpdatePartner(ctx context.Context, in *partnerservicev1.UpdatePartnerRequest, opts ...grpc.CallOption) (*partnerservicev1.UpdatePartnerResponse, error)
	ListPartners(ctx context.Context, in *partnerservicev1.ListPartnersRequest, opts ...grpc.CallOption) (*partnerservicev1.ListPartnersResponse, error)
	DeletePartner(ctx context.Context, in *partnerservicev1.DeletePartnerRequest, opts ...grpc.CallOption) (*partnerservicev1.DeletePartnerResponse, error)
	VerifyPartner(ctx context.Context, in *partnerservicev1.VerifyPartnerRequest, opts ...grpc.CallOption) (*partnerservicev1.VerifyPartnerResponse, error)
	UpdatePartnerStatus(ctx context.Context, in *partnerservicev1.UpdatePartnerStatusRequest, opts ...grpc.CallOption) (*partnerservicev1.UpdatePartnerStatusResponse, error)
	GetPartnerCommission(ctx context.Context, in *partnerservicev1.GetPartnerCommissionRequest, opts ...grpc.CallOption) (*partnerservicev1.GetPartnerCommissionResponse, error)
	UpdateCommissionStructure(ctx context.Context, in *partnerservicev1.UpdateCommissionStructureRequest, opts ...grpc.CallOption) (*partnerservicev1.UpdateCommissionStructureResponse, error)
	GetPartnerAPICredentials(ctx context.Context, in *partnerservicev1.GetPartnerAPICredentialsRequest, opts ...grpc.CallOption) (*partnerservicev1.GetPartnerAPICredentialsResponse, error)
	RotatePartnerAPIKey(ctx context.Context, in *partnerservicev1.RotatePartnerAPIKeyRequest, opts ...grpc.CallOption) (*partnerservicev1.RotatePartnerAPIKeyResponse, error)
}

func NewPartnerHandler(conn *grpc.ClientConn) *PartnerHandler {
	return &PartnerHandler{client: partnerservicev1.NewPartnerServiceClient(conn)}
}

func (h *PartnerHandler) Create(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req partnerservicev1.CreatePartnerRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreatePartner(ctx, &req)
	})
}

func (h *PartnerHandler) Get(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPartner(ctx, &partnerservicev1.GetPartnerRequest{PartnerId: partnerID})
	})
}

func (h *PartnerHandler) Update(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req partnerservicev1.UpdatePartnerRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PartnerId = partnerID
		return h.client.UpdatePartner(ctx, &req)
	})
}

func (h *PartnerHandler) List(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &partnerservicev1.ListPartnersRequest{
			PageToken: r.URL.Query().Get("page_token"),
			Filter:    r.URL.Query().Get("filter"),
			OrderBy:   r.URL.Query().Get("order_by"),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListPartners(ctx, req)
	})
}

func (h *PartnerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeletePartner(ctx, &partnerservicev1.DeletePartnerRequest{PartnerId: partnerID})
	})
}

func (h *PartnerHandler) Verify(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req partnerservicev1.VerifyPartnerRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PartnerId = partnerID
		return h.client.VerifyPartner(ctx, &req)
	})
}

func (h *PartnerHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req partnerservicev1.UpdatePartnerStatusRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PartnerId = partnerID
		return h.client.UpdatePartnerStatus(ctx, &req)
	})
}

func (h *PartnerHandler) GetCommission(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &partnerservicev1.GetPartnerCommissionRequest{PartnerId: partnerID}
		if q := r.URL.Query().Get("start_date"); q != "" {
			if ts, err := parseHTTPTime(q); err == nil {
				req.StartDate = ts
			}
		}
		if q := r.URL.Query().Get("end_date"); q != "" {
			if ts, err := parseHTTPTime(q); err == nil {
				req.EndDate = ts
			}
		}
		return h.client.GetPartnerCommission(ctx, req)
	})
}

func (h *PartnerHandler) UpdateCommission(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req partnerservicev1.UpdateCommissionStructureRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PartnerId = partnerID
		return h.client.UpdateCommissionStructure(ctx, &req)
	})
}

func (h *PartnerHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPartnerAPICredentials(ctx, &partnerservicev1.GetPartnerAPICredentialsRequest{PartnerId: partnerID})
	})
}

func (h *PartnerHandler) RotateAPIKey(w http.ResponseWriter, r *http.Request) {
	partnerID := r.PathValue("partner_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req partnerservicev1.RotatePartnerAPIKeyRequest
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, &req); err != nil {
				return nil, err
			}
		}
		req.PartnerId = partnerID
		return h.client.RotatePartnerAPIKey(ctx, &req)
	})
}
