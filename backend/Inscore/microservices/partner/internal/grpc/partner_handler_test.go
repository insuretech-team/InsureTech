package grpc

import (
	"context"
	"errors"
	"testing"

	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/service"
)

type fakePartnerService struct {
	createFn                 func(context.Context, *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error)
	getFn                    func(context.Context, *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error)
	updateFn                 func(context.Context, *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error)
	listFn                   func(context.Context, *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error)
	deleteFn                 func(context.Context, *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error)
	verifyFn                 func(context.Context, *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error)
	updateStatusFn           func(context.Context, *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error)
	getCommissionFn          func(context.Context, *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error)
	updateCommissionStructFn func(context.Context, *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error)
	getCredentialsFn         func(context.Context, *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error)
	rotateAPIKeyFn           func(context.Context, *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error)
}

func (f *fakePartnerService) CreatePartner(ctx context.Context, req *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error) {
	return f.createFn(ctx, req)
}
func (f *fakePartnerService) GetPartner(ctx context.Context, req *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error) {
	return f.getFn(ctx, req)
}
func (f *fakePartnerService) UpdatePartner(ctx context.Context, req *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error) {
	return f.updateFn(ctx, req)
}
func (f *fakePartnerService) ListPartners(ctx context.Context, req *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error) {
	return f.listFn(ctx, req)
}
func (f *fakePartnerService) DeletePartner(ctx context.Context, req *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error) {
	return f.deleteFn(ctx, req)
}
func (f *fakePartnerService) VerifyPartner(ctx context.Context, req *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error) {
	return f.verifyFn(ctx, req)
}
func (f *fakePartnerService) UpdatePartnerStatus(ctx context.Context, req *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error) {
	return f.updateStatusFn(ctx, req)
}
func (f *fakePartnerService) GetPartnerCommission(ctx context.Context, req *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error) {
	return f.getCommissionFn(ctx, req)
}
func (f *fakePartnerService) UpdateCommissionStructure(ctx context.Context, req *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error) {
	return f.updateCommissionStructFn(ctx, req)
}
func (f *fakePartnerService) GetPartnerAPICredentials(ctx context.Context, req *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error) {
	return f.getCredentialsFn(ctx, req)
}
func (f *fakePartnerService) RotatePartnerAPIKey(ctx context.Context, req *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error) {
	return f.rotateAPIKeyFn(ctx, req)
}
func (f *fakePartnerService) ProcessPolicyCommissionEvent(context.Context, string, partnerv1.CommissionType) error {
	return nil
}

func TestMapError(t *testing.T) {
	assert.Equal(t, codes.InvalidArgument, status.Code(mapError(service.ErrInvalidArgument)))
	assert.Equal(t, codes.NotFound, status.Code(mapError(service.ErrNotFound)))
	assert.Equal(t, codes.AlreadyExists, status.Code(mapError(service.ErrConflict)))
	assert.Equal(t, codes.Unavailable, status.Code(mapError(service.ErrUnavailable)))
	assert.Equal(t, codes.Internal, status.Code(mapError(errors.New("boom"))))
}

func TestPartnerHandlerDelegatesSuccessResponses(t *testing.T) {
	handler := NewPartnerHandler(&fakePartnerService{
		createFn: func(context.Context, *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error) {
			return &partnerservicev1.CreatePartnerResponse{PartnerId: "partner-1"}, nil
		},
		getFn: func(context.Context, *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error) {
			return &partnerservicev1.GetPartnerResponse{Partner: &partnerv1.Partner{PartnerId: "partner-1"}}, nil
		},
		updateFn: func(context.Context, *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error) {
			return &partnerservicev1.UpdatePartnerResponse{Partner: &partnerv1.Partner{PartnerId: "partner-1"}}, nil
		},
		listFn: func(context.Context, *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error) {
			return &partnerservicev1.ListPartnersResponse{Partners: []*partnerv1.Partner{{PartnerId: "partner-1"}}, TotalCount: 1}, nil
		},
		deleteFn: func(context.Context, *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error) {
			return &partnerservicev1.DeletePartnerResponse{Message: "deleted"}, nil
		},
		verifyFn: func(context.Context, *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error) {
			return &partnerservicev1.VerifyPartnerResponse{Verified: true}, nil
		},
		updateStatusFn: func(context.Context, *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error) {
			return &partnerservicev1.UpdatePartnerStatusResponse{Partner: &partnerv1.Partner{PartnerId: "partner-1"}}, nil
		},
		getCommissionFn: func(context.Context, *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error) {
			return &partnerservicev1.GetPartnerCommissionResponse{PartnerId: "partner-1"}, nil
		},
		updateCommissionStructFn: func(context.Context, *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error) {
			return &partnerservicev1.UpdateCommissionStructureResponse{Success: true}, nil
		},
		getCredentialsFn: func(context.Context, *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error) {
			return &partnerservicev1.GetPartnerAPICredentialsResponse{ApiKey: "key-1"}, nil
		},
		rotateAPIKeyFn: func(context.Context, *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error) {
			return &partnerservicev1.RotatePartnerAPIKeyResponse{NewApiKey: "key-2"}, nil
		},
	})

	createResp, err := handler.CreatePartner(context.Background(), &partnerservicev1.CreatePartnerRequest{})
	require.NoError(t, err)
	assert.Equal(t, "partner-1", createResp.GetPartnerId())

	getResp, err := handler.GetPartner(context.Background(), &partnerservicev1.GetPartnerRequest{})
	require.NoError(t, err)
	assert.Equal(t, "partner-1", getResp.GetPartner().GetPartnerId())

	updateResp, err := handler.UpdatePartner(context.Background(), &partnerservicev1.UpdatePartnerRequest{})
	require.NoError(t, err)
	assert.Equal(t, "partner-1", updateResp.GetPartner().GetPartnerId())

	listResp, err := handler.ListPartners(context.Background(), &partnerservicev1.ListPartnersRequest{})
	require.NoError(t, err)
	assert.Len(t, listResp.GetPartners(), 1)

	deleteResp, err := handler.DeletePartner(context.Background(), &partnerservicev1.DeletePartnerRequest{})
	require.NoError(t, err)
	assert.Contains(t, deleteResp.GetMessage(), "deleted")

	verifyResp, err := handler.VerifyPartner(context.Background(), &partnerservicev1.VerifyPartnerRequest{})
	require.NoError(t, err)
	assert.True(t, verifyResp.GetVerified())

	statusResp, err := handler.UpdatePartnerStatus(context.Background(), &partnerservicev1.UpdatePartnerStatusRequest{})
	require.NoError(t, err)
	assert.Equal(t, "partner-1", statusResp.GetPartner().GetPartnerId())

	commissionResp, err := handler.GetPartnerCommission(context.Background(), &partnerservicev1.GetPartnerCommissionRequest{})
	require.NoError(t, err)
	assert.Equal(t, "partner-1", commissionResp.GetPartnerId())

	structureResp, err := handler.UpdateCommissionStructure(context.Background(), &partnerservicev1.UpdateCommissionStructureRequest{})
	require.NoError(t, err)
	assert.True(t, structureResp.GetSuccess())

	credsResp, err := handler.GetPartnerAPICredentials(context.Background(), &partnerservicev1.GetPartnerAPICredentialsRequest{})
	require.NoError(t, err)
	assert.Equal(t, "key-1", credsResp.GetApiKey())

	rotateResp, err := handler.RotatePartnerAPIKey(context.Background(), &partnerservicev1.RotatePartnerAPIKeyRequest{})
	require.NoError(t, err)
	assert.Equal(t, "key-2", rotateResp.GetNewApiKey())
}

func TestPartnerHandlerMapsServiceErrors(t *testing.T) {
	sentinelErrors := []error{
		service.ErrInvalidArgument,
		service.ErrNotFound,
		service.ErrConflict,
		service.ErrUnavailable,
		errors.New("boom"),
	}

	for _, expectedErr := range sentinelErrors {
		handler := NewPartnerHandler(&fakePartnerService{
			createFn: func(context.Context, *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error) {
				return nil, expectedErr
			},
			getFn: func(context.Context, *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error) {
				return nil, expectedErr
			},
			updateFn: func(context.Context, *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error) {
				return nil, expectedErr
			},
			listFn: func(context.Context, *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error) {
				return nil, expectedErr
			},
			deleteFn: func(context.Context, *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error) {
				return nil, expectedErr
			},
			verifyFn: func(context.Context, *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error) {
				return nil, expectedErr
			},
			updateStatusFn: func(context.Context, *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error) {
				return nil, expectedErr
			},
			getCommissionFn: func(context.Context, *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error) {
				return nil, expectedErr
			},
			updateCommissionStructFn: func(context.Context, *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error) {
				return nil, expectedErr
			},
			getCredentialsFn: func(context.Context, *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error) {
				return nil, expectedErr
			},
			rotateAPIKeyFn: func(context.Context, *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error) {
				return nil, expectedErr
			},
		})

		_, err := handler.CreatePartner(context.Background(), &partnerservicev1.CreatePartnerRequest{})
		require.Error(t, err)
		assert.Equal(t, status.Code(mapError(expectedErr)), status.Code(err))

		_, err = handler.RotatePartnerAPIKey(context.Background(), &partnerservicev1.RotatePartnerAPIKeyRequest{})
		require.Error(t, err)
		assert.Equal(t, status.Code(mapError(expectedErr)), status.Code(err))
	}
}

func TestServerHelpers(t *testing.T) {
	cfg := DefaultServerConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, "50058", cfg.Port)

	server, err := NewServer(&Config{Port: "50058"}, &fakePartnerService{
		createFn: func(context.Context, *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error) {
			return &partnerservicev1.CreatePartnerResponse{}, nil
		},
		getFn: func(context.Context, *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error) {
			return &partnerservicev1.GetPartnerResponse{}, nil
		},
		updateFn: func(context.Context, *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error) {
			return &partnerservicev1.UpdatePartnerResponse{}, nil
		},
		listFn: func(context.Context, *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error) {
			return &partnerservicev1.ListPartnersResponse{}, nil
		},
		deleteFn: func(context.Context, *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error) {
			return &partnerservicev1.DeletePartnerResponse{}, nil
		},
		verifyFn: func(context.Context, *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error) {
			return &partnerservicev1.VerifyPartnerResponse{}, nil
		},
		updateStatusFn: func(context.Context, *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error) {
			return &partnerservicev1.UpdatePartnerStatusResponse{}, nil
		},
		getCommissionFn: func(context.Context, *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error) {
			return &partnerservicev1.GetPartnerCommissionResponse{}, nil
		},
		updateCommissionStructFn: func(context.Context, *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error) {
			return &partnerservicev1.UpdateCommissionStructureResponse{}, nil
		},
		getCredentialsFn: func(context.Context, *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error) {
			return &partnerservicev1.GetPartnerAPICredentialsResponse{ApiKey: "key-1"}, nil
		},
		rotateAPIKeyFn: func(context.Context, *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error) {
			return &partnerservicev1.RotatePartnerAPIKeyResponse{NewApiKey: "key-2"}, nil
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, server)

	err = server.HealthCheck(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database connection is nil")

	server.Stop()
}
