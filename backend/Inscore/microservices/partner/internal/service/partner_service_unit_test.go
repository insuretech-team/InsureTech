package service

import (
	"context"
	"errors"
	"testing"
	"time"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/domain"
)

type fakePartnerRepo struct {
	createFn          func(context.Context, *partnerv1.Partner) error
	getByIDFn         func(context.Context, string) (*partnerv1.Partner, error)
	updateStatusFn    func(context.Context, string, partnerv1.PartnerStatus) error
	listWithFiltersFn func(context.Context, int, int, string, string) ([]*partnerv1.Partner, int32, error)
	updateFn          func(context.Context, string, *partnerv1.Partner, []string) error
	softDeleteFn      func(context.Context, string) error
}

func (f *fakePartnerRepo) Create(ctx context.Context, partner *partnerv1.Partner) error {
	return f.createFn(ctx, partner)
}

func (f *fakePartnerRepo) GetByID(ctx context.Context, id string) (*partnerv1.Partner, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakePartnerRepo) UpdateStatus(ctx context.Context, partnerID string, status partnerv1.PartnerStatus) error {
	return f.updateStatusFn(ctx, partnerID, status)
}

func (f *fakePartnerRepo) ListWithFilters(ctx context.Context, limit, offset int, filter, orderBy string) ([]*partnerv1.Partner, int32, error) {
	return f.listWithFiltersFn(ctx, limit, offset, filter, orderBy)
}

func (f *fakePartnerRepo) Update(ctx context.Context, partnerID string, partner *partnerv1.Partner, updateMask []string) error {
	return f.updateFn(ctx, partnerID, partner, updateMask)
}

func (f *fakePartnerRepo) SoftDelete(ctx context.Context, partnerID string) error {
	return f.softDeleteFn(ctx, partnerID)
}

type fakeCommissionRepo struct {
	createFn                    func(context.Context, *partnerv1.Commission) error
	listByPartnerAndDateRangeFn func(context.Context, string, *time.Time, *time.Time, int, int) ([]*partnerv1.Commission, int32, error)
	sumByPartnerAndDateRangeFn  func(context.Context, string, *time.Time, *time.Time) (int64, error)
	existsByPolicyAndTypeFn     func(context.Context, string, partnerv1.CommissionType) (bool, error)
	resolvePolicyInputFn        func(context.Context, string) (*domain.PolicyCommissionInput, error)
}

func (f *fakeCommissionRepo) Create(ctx context.Context, comm *partnerv1.Commission) error {
	return f.createFn(ctx, comm)
}

func (f *fakeCommissionRepo) ListByPartnerAndDateRange(ctx context.Context, partnerID string, start, end *time.Time, limit, offset int) ([]*partnerv1.Commission, int32, error) {
	return f.listByPartnerAndDateRangeFn(ctx, partnerID, start, end, limit, offset)
}

func (f *fakeCommissionRepo) SumByPartnerAndDateRange(ctx context.Context, partnerID string, start, end *time.Time) (int64, error) {
	return f.sumByPartnerAndDateRangeFn(ctx, partnerID, start, end)
}

func (f *fakeCommissionRepo) ExistsByPolicyAndType(ctx context.Context, policyID string, cType partnerv1.CommissionType) (bool, error) {
	return f.existsByPolicyAndTypeFn(ctx, policyID, cType)
}

func (f *fakeCommissionRepo) ResolvePolicyCommissionInput(ctx context.Context, policyID string) (*domain.PolicyCommissionInput, error) {
	return f.resolvePolicyInputFn(ctx, policyID)
}

type fakeEventPublisher struct {
	onboardedPartners  []string
	verifiedPartnerIDs []string
	verifiedBy         []string
	commissionEvents   []*partnerv1.Commission
}

func (f *fakeEventPublisher) PublishPartnerOnboarded(_ context.Context, partner *partnerv1.Partner) error {
	f.onboardedPartners = append(f.onboardedPartners, partner.GetPartnerId())
	return nil
}

func (f *fakeEventPublisher) PublishPartnerVerified(_ context.Context, partnerID string, verifiedBy string) error {
	f.verifiedPartnerIDs = append(f.verifiedPartnerIDs, partnerID)
	f.verifiedBy = append(f.verifiedBy, verifiedBy)
	return nil
}

func (f *fakeEventPublisher) PublishAgentRegistered(context.Context, *partnerv1.Agent) error {
	return nil
}

func (f *fakeEventPublisher) PublishCommissionCalculated(_ context.Context, commission *partnerv1.Commission) error {
	f.commissionEvents = append(f.commissionEvents, commission)
	return nil
}

type fakeAuthnClient struct {
	listFn   func(context.Context, *authnservicev1.ListAPIKeysRequest, ...grpc.CallOption) (*authnservicev1.ListAPIKeysResponse, error)
	createFn func(context.Context, *authnservicev1.CreateAPIKeyRequest, ...grpc.CallOption) (*authnservicev1.CreateAPIKeyResponse, error)
	rotateFn func(context.Context, *authnservicev1.RotateAPIKeyRequest, ...grpc.CallOption) (*authnservicev1.RotateAPIKeyResponse, error)
}

func (f *fakeAuthnClient) ListAPIKeys(ctx context.Context, in *authnservicev1.ListAPIKeysRequest, opts ...grpc.CallOption) (*authnservicev1.ListAPIKeysResponse, error) {
	return f.listFn(ctx, in, opts...)
}

func (f *fakeAuthnClient) CreateAPIKey(ctx context.Context, in *authnservicev1.CreateAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.CreateAPIKeyResponse, error) {
	return f.createFn(ctx, in, opts...)
}

func (f *fakeAuthnClient) RotateAPIKey(ctx context.Context, in *authnservicev1.RotateAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.RotateAPIKeyResponse, error) {
	return f.rotateFn(ctx, in, opts...)
}

func TestPartnerCRUDServiceFlows(t *testing.T) {
	t.Run("create partner validates and publishes", func(t *testing.T) {
		pub := &fakeEventPublisher{}
		svc := NewPartnerService(
			&fakePartnerRepo{
				createFn: func(context.Context, *partnerv1.Partner) error { return nil },
			},
			nil,
			nil,
			pub,
			nil,
		)

		_, err := svc.CreatePartner(context.Background(), &partnerservicev1.CreatePartnerRequest{})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)

		resp, err := svc.CreatePartner(context.Background(), &partnerservicev1.CreatePartnerRequest{
			Partner: &partnerv1.Partner{PartnerId: "partner-1"},
		})
		require.NoError(t, err)
		assert.Equal(t, "partner-1", resp.GetPartnerId())
		assert.Equal(t, []string{"partner-1"}, pub.onboardedPartners)
		assert.Equal(t, int64(1), svc.MetricsSnapshot()["partner_created"])
	})

	t.Run("get, update, list and delete map repository behavior", func(t *testing.T) {
		var listedLimit int
		var listedOffset int
		var listedFilter string
		var listedOrder string

		svc := NewPartnerService(
			&fakePartnerRepo{
				getByIDFn: func(_ context.Context, id string) (*partnerv1.Partner, error) {
					if id == "missing" {
						return nil, gorm.ErrRecordNotFound
					}
					return &partnerv1.Partner{PartnerId: id, OrganizationName: "Org " + id}, nil
				},
				updateFn: func(_ context.Context, id string, partner *partnerv1.Partner, updateMask []string) error {
					if id == "missing" {
						return gorm.ErrRecordNotFound
					}
					assert.Equal(t, []string{"organization_name"}, updateMask)
					assert.Equal(t, "Updated", partner.GetOrganizationName())
					return nil
				},
				listWithFiltersFn: func(_ context.Context, limit, offset int, filter, orderBy string) ([]*partnerv1.Partner, int32, error) {
					listedLimit = limit
					listedOffset = offset
					listedFilter = filter
					listedOrder = orderBy
					return []*partnerv1.Partner{{PartnerId: "p-1"}, {PartnerId: "p-2"}}, 4, nil
				},
				softDeleteFn: func(_ context.Context, id string) error {
					if id == "missing" {
						return gorm.ErrRecordNotFound
					}
					return nil
				},
			},
			nil,
			nil,
			nil,
			nil,
		)

		_, err := svc.GetPartner(context.Background(), &partnerservicev1.GetPartnerRequest{})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)

		getResp, err := svc.GetPartner(context.Background(), &partnerservicev1.GetPartnerRequest{PartnerId: "partner-1"})
		require.NoError(t, err)
		assert.Equal(t, "partner-1", getResp.GetPartner().GetPartnerId())

		_, err = svc.GetPartner(context.Background(), &partnerservicev1.GetPartnerRequest{PartnerId: "missing"})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)

		_, err = svc.UpdatePartner(context.Background(), &partnerservicev1.UpdatePartnerRequest{})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)

		updateResp, err := svc.UpdatePartner(context.Background(), &partnerservicev1.UpdatePartnerRequest{
			PartnerId:  "partner-1",
			Partner:    &partnerv1.Partner{OrganizationName: "Updated"},
			UpdateMask: []string{"organization_name"},
		})
		require.NoError(t, err)
		assert.Equal(t, "partner-1", updateResp.GetPartner().GetPartnerId())

		_, err = svc.UpdatePartner(context.Background(), &partnerservicev1.UpdatePartnerRequest{
			PartnerId:  "missing",
			Partner:    &partnerv1.Partner{OrganizationName: "Updated"},
			UpdateMask: []string{"organization_name"},
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)

		listResp, err := svc.ListPartners(context.Background(), &partnerservicev1.ListPartnersRequest{
			PageSize:  0,
			PageToken: "1",
			Filter:    "status=ACTIVE",
			OrderBy:   "-created_at",
		})
		require.NoError(t, err)
		assert.Len(t, listResp.GetPartners(), 2)
		assert.Equal(t, int32(4), listResp.GetTotalCount())
		assert.Equal(t, "3", listResp.GetNextPageToken())
		assert.Equal(t, 50, listedLimit)
		assert.Equal(t, 1, listedOffset)
		assert.Equal(t, "status=ACTIVE", listedFilter)
		assert.Equal(t, "-created_at", listedOrder)

		deleteResp, err := svc.DeletePartner(context.Background(), &partnerservicev1.DeletePartnerRequest{PartnerId: "partner-1"})
		require.NoError(t, err)
		assert.Contains(t, deleteResp.GetMessage(), "deleted")
		assert.Equal(t, int64(1), svc.MetricsSnapshot()["partner_deleted"])

		_, err = svc.DeletePartner(context.Background(), &partnerservicev1.DeletePartnerRequest{PartnerId: "missing"})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestPartnerVerificationAndStatusFlows(t *testing.T) {
	pub := &fakeEventPublisher{}
	svc := NewPartnerService(
		&fakePartnerRepo{
			updateStatusFn: func(_ context.Context, partnerID string, status partnerv1.PartnerStatus) error {
				if partnerID == "missing" {
					return gorm.ErrRecordNotFound
				}
				assert.NotEqual(t, partnerv1.PartnerStatus_PARTNER_STATUS_UNSPECIFIED, status)
				return nil
			},
			getByIDFn: func(_ context.Context, id string) (*partnerv1.Partner, error) {
				if id == "missing" {
					return nil, gorm.ErrRecordNotFound
				}
				return &partnerv1.Partner{PartnerId: id, Status: partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE}, nil
			},
		},
		nil,
		nil,
		pub,
		nil,
	)

	_, err := svc.VerifyPartner(context.Background(), &partnerservicev1.VerifyPartnerRequest{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidArgument)

	verifyResp, err := svc.VerifyPartner(context.Background(), &partnerservicev1.VerifyPartnerRequest{
		PartnerId:        "partner-1",
		VerificationType: "MANUAL",
	})
	require.NoError(t, err)
	assert.True(t, verifyResp.GetVerified())
	assert.Equal(t, "system", verifyResp.GetVerifiedBy())
	assert.Equal(t, []string{"partner-1"}, pub.verifiedPartnerIDs)
	assert.Equal(t, []string{"system"}, pub.verifiedBy)

	_, err = svc.VerifyPartner(context.Background(), &partnerservicev1.VerifyPartnerRequest{
		PartnerId:        "missing",
		VerificationType: "MANUAL",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)

	_, err = svc.UpdatePartnerStatus(context.Background(), &partnerservicev1.UpdatePartnerStatusRequest{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidArgument)

	_, err = svc.UpdatePartnerStatus(context.Background(), &partnerservicev1.UpdatePartnerStatusRequest{
		PartnerId: "partner-1",
		Status:    "not-a-status",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidArgument)

	statusResp, err := svc.UpdatePartnerStatus(context.Background(), &partnerservicev1.UpdatePartnerStatusRequest{
		PartnerId: "partner-1",
		Status:    "active",
	})
	require.NoError(t, err)
	assert.Equal(t, "partner-1", statusResp.GetPartner().GetPartnerId())
	assert.Equal(t, "admin", pub.verifiedBy[len(pub.verifiedBy)-1])
}

func TestPartnerCommissionFlows(t *testing.T) {
	pub := &fakeEventPublisher{}
	var createdCommission *partnerv1.Commission
	svc := NewPartnerService(
		&fakePartnerRepo{
			getByIDFn: func(_ context.Context, id string) (*partnerv1.Partner, error) {
				if id == "missing" {
					return nil, gorm.ErrRecordNotFound
				}
				return &partnerv1.Partner{
					PartnerId:                 id,
					AcquisitionCommissionRate: 10,
					RenewalCommissionRate:     5,
					ClaimsAssistanceRate:      2,
				}, nil
			},
			updateFn: func(context.Context, string, *partnerv1.Partner, []string) error { return nil },
		},
		nil,
		&fakeCommissionRepo{
			listByPartnerAndDateRangeFn: func(_ context.Context, partnerID string, start, end *time.Time, limit, offset int) ([]*partnerv1.Commission, int32, error) {
				assert.Equal(t, 200, limit)
				return []*partnerv1.Commission{{
					PolicyId: "policy-1",
					Type:     partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION,
					CommissionAmount: &commonv1.Money{
						Amount:   1500,
						Currency: "",
					},
					CreatedAt: timestamppb.Now(),
				}}, 1, nil
			},
			sumByPartnerAndDateRangeFn: func(context.Context, string, *time.Time, *time.Time) (int64, error) {
				return 1500, nil
			},
			existsByPolicyAndTypeFn: func(_ context.Context, policyID string, cType partnerv1.CommissionType) (bool, error) {
				return policyID == "duplicate", nil
			},
			resolvePolicyInputFn: func(_ context.Context, policyID string) (*domain.PolicyCommissionInput, error) {
				switch policyID {
				case "missing":
					return nil, gorm.ErrRecordNotFound
				case "no-partner":
					return &domain.PolicyCommissionInput{PolicyID: policyID, PremiumAmount: 10000}, nil
				default:
					return &domain.PolicyCommissionInput{
						PolicyID:      policyID,
						PartnerID:     "partner-1",
						AgentID:       "agent-1",
						PremiumAmount: 10000,
						Currency:      "BDT",
					}, nil
				}
			},
			createFn: func(_ context.Context, comm *partnerv1.Commission) error {
				createdCommission = comm
				return nil
			},
		},
		pub,
		nil,
	)

	_, err := svc.GetPartnerCommission(context.Background(), &partnerservicev1.GetPartnerCommissionRequest{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidArgument)

	commissionResp, err := svc.GetPartnerCommission(context.Background(), &partnerservicev1.GetPartnerCommissionRequest{
		PartnerId: "partner-1",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1500), commissionResp.GetTotalCommission().GetAmount())
	assert.Equal(t, "BDT", commissionResp.GetCurrency())
	assert.Len(t, commissionResp.GetDetails(), 1)
	assert.Equal(t, "ACQUISITION", commissionResp.GetDetails()[0].GetCommissionType())

	updateResp, err := svc.UpdateCommissionStructure(context.Background(), &partnerservicev1.UpdateCommissionStructureRequest{
		PartnerId: "partner-1",
		CommissionRates: map[string]float64{
			"acquisition": 12,
			"renewal":     6,
		},
	})
	require.NoError(t, err)
	assert.True(t, updateResp.GetSuccess())

	err = svc.ProcessPolicyCommissionEvent(context.Background(), "duplicate", partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION)
	require.NoError(t, err)

	err = svc.ProcessPolicyCommissionEvent(context.Background(), "no-partner", partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION)
	require.NoError(t, err)

	err = svc.ProcessPolicyCommissionEvent(context.Background(), "missing", partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)

	err = svc.ProcessPolicyCommissionEvent(context.Background(), "policy-1", partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION)
	require.NoError(t, err)
	require.NotNil(t, createdCommission)
	assert.Equal(t, int64(1000), createdCommission.GetCommissionAmount().GetAmount())
	assert.Len(t, pub.commissionEvents, 1)
}

func TestPartnerAPIKeyFlows(t *testing.T) {
	expiresAt := timestamppb.New(time.Now().Add(24 * time.Hour))
	oldExpiresAt := timestamppb.New(time.Now().Add(2 * time.Hour))
	authnClient := &fakeAuthnClient{
		listFn: func(_ context.Context, in *authnservicev1.ListAPIKeysRequest, _ ...grpc.CallOption) (*authnservicev1.ListAPIKeysResponse, error) {
			switch in.PageSize {
			case 10:
				if in.OwnerId == "partner-with-key" {
					return &authnservicev1.ListAPIKeysResponse{
						Keys: []*authnservicev1.APIKeySummary{{
							KeyId:     "key-1",
							Status:    "active",
							ExpiresAt: expiresAt,
						}},
					}, nil
				}
				if in.OwnerId == "rotate-existing" {
					return &authnservicev1.ListAPIKeysResponse{
						Keys: []*authnservicev1.APIKeySummary{{
							KeyId:     "old-key",
							Status:    "active",
							ExpiresAt: oldExpiresAt,
						}},
					}, nil
				}
				if in.OwnerId == "list-fails" {
					return nil, errors.New("list failed")
				}
				return &authnservicev1.ListAPIKeysResponse{}, nil
			default:
				return &authnservicev1.ListAPIKeysResponse{
					Keys: []*authnservicev1.APIKeySummary{{
						KeyId:     "rotated-key",
						Status:    "active",
						ExpiresAt: expiresAt,
					}},
				}, nil
			}
		},
		createFn: func(_ context.Context, in *authnservicev1.CreateAPIKeyRequest, _ ...grpc.CallOption) (*authnservicev1.CreateAPIKeyResponse, error) {
			return &authnservicev1.CreateAPIKeyResponse{
				KeyId:     "created-key",
				RawKey:    "secret",
				ExpiresAt: expiresAt,
			}, nil
		},
		rotateFn: func(_ context.Context, in *authnservicev1.RotateAPIKeyRequest, _ ...grpc.CallOption) (*authnservicev1.RotateAPIKeyResponse, error) {
			assert.Equal(t, "old-key", in.GetKeyId())
			return &authnservicev1.RotateAPIKeyResponse{
				NewKeyId:        "rotated-key",
				RawKey:          "rotated-secret",
				OldKeyExpiresAt: oldExpiresAt,
			}, nil
		},
	}

	svc := NewPartnerService(
		&fakePartnerRepo{},
		nil,
		nil,
		nil,
		authnClient,
	)

	_, err := svc.GetPartnerAPICredentials(context.Background(), &partnerservicev1.GetPartnerAPICredentialsRequest{})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidArgument)

	getResp, err := svc.GetPartnerAPICredentials(context.Background(), &partnerservicev1.GetPartnerAPICredentialsRequest{
		PartnerId: "partner-with-key",
	})
	require.NoError(t, err)
	assert.Equal(t, "key-1", getResp.GetApiKey())
	assert.Empty(t, getResp.GetApiSecret())

	createResp, err := svc.GetPartnerAPICredentials(context.Background(), &partnerservicev1.GetPartnerAPICredentialsRequest{
		PartnerId: "partner-new",
	})
	require.NoError(t, err)
	assert.Equal(t, "created-key", createResp.GetApiKey())
	assert.Equal(t, "secret", createResp.GetApiSecret())

	rotateResp, err := svc.RotatePartnerAPIKey(context.Background(), &partnerservicev1.RotatePartnerAPIKeyRequest{
		PartnerId: "rotate-existing",
	})
	require.NoError(t, err)
	assert.Equal(t, "rotated-key", rotateResp.GetNewApiKey())
	assert.Equal(t, expiresAt, rotateResp.GetExpiresAt())
	assert.Equal(t, int64(1), svc.MetricsSnapshot()["api_key_rotated"])

	rotateCreateResp, err := svc.RotatePartnerAPIKey(context.Background(), &partnerservicev1.RotatePartnerAPIKeyRequest{
		PartnerId: "partner-without-key",
	})
	require.NoError(t, err)
	assert.Equal(t, "created-key", rotateCreateResp.GetNewApiKey())

	_, err = svc.GetPartnerAPICredentials(context.Background(), &partnerservicev1.GetPartnerAPICredentialsRequest{
		PartnerId: "list-fails",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lookup partner API keys")
}

func TestPartnerServiceHelpers(t *testing.T) {
	assert.Equal(t, 0, decodePartnerPageToken(""))
	assert.Equal(t, 0, decodePartnerPageToken("not-a-number"))
	assert.Equal(t, 7, decodePartnerPageToken("7"))

	status, ok := parsePartnerStatus("active")
	assert.True(t, ok)
	assert.Equal(t, partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE, status)

	status, ok = parsePartnerStatus("PARTNER_STATUS_SUSPENDED")
	assert.True(t, ok)
	assert.Equal(t, partnerv1.PartnerStatus_PARTNER_STATUS_SUSPENDED, status)

	status, ok = parsePartnerStatus("")
	assert.False(t, ok)
	assert.Equal(t, partnerv1.PartnerStatus_PARTNER_STATUS_UNSPECIFIED, status)

	startTS := timestamppb.New(time.Now().Add(-time.Hour))
	endTS := timestamppb.New(time.Now())
	start, end := toTimeRange(startTS, endTS)
	require.NotNil(t, start)
	require.NotNil(t, end)

	partner := &partnerv1.Partner{AcquisitionCommissionRate: 10, RenewalCommissionRate: 5, ClaimsAssistanceRate: 2}
	applyCommissionRates(partner, map[string]float64{
		"acquisition_rate":       12,
		"claims_assistance_rate": 3,
	})
	assert.Equal(t, float64(12), partner.GetAcquisitionCommissionRate())
	assert.Equal(t, float64(5), partner.GetRenewalCommissionRate())
	assert.Equal(t, float64(3), partner.GetClaimsAssistanceRate())

	assert.Equal(t, float64(12), commissionRateByType(partner, partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION))
	assert.Equal(t, float64(5), commissionRateByType(partner, partnerv1.CommissionType_COMMISSION_TYPE_RENEWAL))
	assert.Equal(t, float64(3), commissionRateByType(partner, partnerv1.CommissionType_COMMISSION_TYPE_CLAIMS_ASSISTANCE))
	assert.Equal(t, defaultPartnerAPIScopes(), []string{"policy:read", "policy:write", "claim:read", "claim:write", "partner:read"})
}
