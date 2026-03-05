package handlers

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
)

// KYCHandler proxies KYC verification requests to the KYC gRPC service.
type KYCHandler struct {
	client kycservicev1.KYCServiceClient
}

// NewKYCHandler creates a KYCHandler from a gRPC connection to the KYC service.
func NewKYCHandler(conn *grpc.ClientConn) *KYCHandler {
	return &KYCHandler{client: kycservicev1.NewKYCServiceClient(conn)}
}

// Create initiates a new KYC verification.
// POST /v1/kyc/verifications
func (h *KYCHandler) Create(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req kycservicev1.StartKYCVerificationRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.StartKYCVerification(ctx, &req)
	})
}

// Get retrieves a KYC verification by ID.
// GET /v1/kyc/verifications/{id}
func (h *KYCHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetKYCVerification(ctx, &kycservicev1.GetKYCVerificationRequest{
			KycVerificationId: id,
		})
	})
}

// List retrieves pending KYC verifications.
// GET /v1/kyc/verifications
func (h *KYCHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListPendingVerifications(ctx, &kycservicev1.ListPendingVerificationsRequest{
			PageToken: q.Get("page_token"),
		})
	})
}
