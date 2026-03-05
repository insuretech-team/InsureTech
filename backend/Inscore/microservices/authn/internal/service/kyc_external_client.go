package service

import (
	"context"

	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"google.golang.org/grpc"
)

// ExternalKYCClient is the minimal downstream KYC dependency used by AuthN.
// It is intentionally narrow so unit/live tests can keep using the local path.
type ExternalKYCClient interface {
	StartKYCVerification(ctx context.Context, in *kycservicev1.StartKYCVerificationRequest, opts ...grpc.CallOption) (*kycservicev1.StartKYCVerificationResponse, error)
	UploadDocument(ctx context.Context, in *kycservicev1.UploadDocumentRequest, opts ...grpc.CallOption) (*kycservicev1.UploadDocumentResponse, error)
	VerifyKYC(ctx context.Context, in *kycservicev1.VerifyKYCRequest, opts ...grpc.CallOption) (*kycservicev1.VerifyKYCResponse, error)
}

// SetExternalKYCClient injects optional downstream KYC client.
// When nil, AuthN falls back to local repository-based KYC flow.
func (s *AuthService) SetExternalKYCClient(client ExternalKYCClient) {
	if s == nil {
		return
	}
	s.externalKYC = client
}
