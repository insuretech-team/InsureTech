package domain

import (
	"context"
	"time"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"google.golang.org/grpc"
)

// PolicyCommissionInput is the minimum policy data required for commission calculation.
type PolicyCommissionInput struct {
	PolicyID      string
	PartnerID     string
	AgentID       string
	PremiumAmount int64
	Currency      string
}

// PartnerRepository defines persistence contract for partner aggregate.
type PartnerRepository interface {
	Create(ctx context.Context, partner *partnerv1.Partner) error
	GetByID(ctx context.Context, id string) (*partnerv1.Partner, error)
	UpdateStatus(ctx context.Context, partnerID string, status partnerv1.PartnerStatus) error
	ListWithFilters(ctx context.Context, limit, offset int, filter, orderBy string) ([]*partnerv1.Partner, int32, error)
	Update(ctx context.Context, partnerID string, partner *partnerv1.Partner, updateMask []string) error
	SoftDelete(ctx context.Context, partnerID string) error
}

// AgentRepository defines persistence contract for agents.
type AgentRepository interface{}

// CommissionRepository defines commission persistence and lookup operations.
type CommissionRepository interface {
	Create(ctx context.Context, comm *partnerv1.Commission) error
	ListByPartnerAndDateRange(ctx context.Context, partnerID string, start, end *time.Time, limit, offset int) ([]*partnerv1.Commission, int32, error)
	SumByPartnerAndDateRange(ctx context.Context, partnerID string, start, end *time.Time) (int64, error)
	ExistsByPolicyAndType(ctx context.Context, policyID string, cType partnerv1.CommissionType) (bool, error)
	ResolvePolicyCommissionInput(ctx context.Context, policyID string) (*PolicyCommissionInput, error)
}

// EventPublisher defines partner domain event publishing contract.
type EventPublisher interface {
	PublishPartnerOnboarded(ctx context.Context, partner *partnerv1.Partner) error
	PublishPartnerVerified(ctx context.Context, partnerID string, verifiedBy string) error
	PublishAgentRegistered(ctx context.Context, agent *partnerv1.Agent) error
	PublishCommissionCalculated(ctx context.Context, commission *partnerv1.Commission) error
}

// AuthNClient defines required AuthN API key methods for partner integration.
type AuthNClient interface {
	ListAPIKeys(ctx context.Context, in *authnservicev1.ListAPIKeysRequest, opts ...grpc.CallOption) (*authnservicev1.ListAPIKeysResponse, error)
	CreateAPIKey(ctx context.Context, in *authnservicev1.CreateAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.CreateAPIKeyResponse, error)
	RotateAPIKey(ctx context.Context, in *authnservicev1.RotateAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.RotateAPIKeyResponse, error)
}

// PartnerService defines business contract consumed by transport and events.
type PartnerService interface {
	CreatePartner(ctx context.Context, req *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error)
	GetPartner(ctx context.Context, req *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error)
	UpdatePartner(ctx context.Context, req *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error)
	ListPartners(ctx context.Context, req *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error)
	DeletePartner(ctx context.Context, req *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error)
	VerifyPartner(ctx context.Context, req *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error)
	UpdatePartnerStatus(ctx context.Context, req *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error)
	GetPartnerCommission(ctx context.Context, req *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error)
	UpdateCommissionStructure(ctx context.Context, req *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error)
	GetPartnerAPICredentials(ctx context.Context, req *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error)
	RotatePartnerAPIKey(ctx context.Context, req *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error)
	ProcessPolicyCommissionEvent(ctx context.Context, policyID string, cType partnerv1.CommissionType) error
}
