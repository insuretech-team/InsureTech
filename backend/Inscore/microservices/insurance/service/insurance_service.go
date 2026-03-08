package service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"go.uber.org/zap"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/insurance/internal/repository"
	insurancev1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurance/services/v1"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

type InsuranceService struct {
	insurancev1.UnimplementedInsuranceServiceServer
	productRepo                  *repository.ProductRepository
	productPlanRepo              *repository.ProductPlanRepository
	riderRepo                    *repository.RiderRepository
	pricingRepo                  *repository.PricingConfigRepository
	policyRepo                   *repository.PolicyRepository
	claimRepo                    *repository.ClaimRepository
	quoteRepo                    *repository.QuoteRepository
	underwritingDecisionRepo     *repository.UnderwritingDecisionRepository
	healthDeclarationRepo        *repository.HealthDeclarationRepository
	renewalScheduleRepo          *repository.RenewalScheduleRepository
	renewalReminderRepo          *repository.RenewalReminderRepository
	gracePeriodRepo              *repository.GracePeriodRepository
	insurerRepo                  *repository.InsurerRepository
	insurerConfigRepo            *repository.InsurerConfigRepository
	insurerProductRepo           *repository.InsurerProductRepository
	fraudRuleRepo                *repository.FraudRuleRepository
	fraudCaseRepo                *repository.FraudCaseRepository
	fraudAlertRepo               *repository.FraudAlertRepository
	beneficiaryRepo              *repository.BeneficiaryRepository
	individualBeneficiaryRepo    *repository.IndividualBeneficiaryRepository
	businessBeneficiaryRepo      *repository.BusinessBeneficiaryRepository
	endorsementRepo              *repository.EndorsementRepository
	quotationRepo                *repository.QuotationRepository
	policyServiceRequestRepo     *repository.PolicyServiceRequestRepository
	serviceProviderRepo          *repository.ServiceProviderRepository
}

func NewInsuranceService(db *gorm.DB) *InsuranceService {
	return &InsuranceService{
		productRepo:              repository.NewProductRepository(db),
		productPlanRepo:          repository.NewProductPlanRepository(db),
		riderRepo:                repository.NewRiderRepository(db),
		pricingRepo:              repository.NewPricingConfigRepository(db),
		policyRepo:               repository.NewPolicyRepository(db),
		claimRepo:                repository.NewClaimRepository(db),
		quoteRepo:                repository.NewQuoteRepository(db),
		underwritingDecisionRepo: repository.NewUnderwritingDecisionRepository(db),
		healthDeclarationRepo:    repository.NewHealthDeclarationRepository(db),
		renewalScheduleRepo:      repository.NewRenewalScheduleRepository(db),
		renewalReminderRepo:      repository.NewRenewalReminderRepository(db),
		gracePeriodRepo:          repository.NewGracePeriodRepository(db),
		insurerRepo:              repository.NewInsurerRepository(db),
		insurerConfigRepo:        repository.NewInsurerConfigRepository(db),
		insurerProductRepo:       repository.NewInsurerProductRepository(db),
		fraudRuleRepo:            repository.NewFraudRuleRepository(db),
		fraudCaseRepo:            repository.NewFraudCaseRepository(db),
		fraudAlertRepo:           repository.NewFraudAlertRepository(db),
		beneficiaryRepo:          repository.NewBeneficiaryRepository(db),
		individualBeneficiaryRepo: repository.NewIndividualBeneficiaryRepository(db),
		businessBeneficiaryRepo:  repository.NewBusinessBeneficiaryRepository(db),
		endorsementRepo:          repository.NewEndorsementRepository(db),
		quotationRepo:            repository.NewQuotationRepository(db),
		policyServiceRequestRepo: repository.NewPolicyServiceRequestRepository(db),
		serviceProviderRepo:      repository.NewServiceProviderRepository(db),
	}
}

// ========== PRODUCT CRUD ==========

func (s *InsuranceService) CreateProduct(ctx context.Context, req *insurancev1.CreateProductRequest) (*insurancev1.CreateProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	product, err := s.productRepo.Create(ctx, req.Product)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		// Return more specific error messages
		errMsg := err.Error()
		if strings.Contains(errMsg, "duplicate key") {
			return nil, status.Error(codes.AlreadyExists, "product with this code already exists")
		}
		if strings.Contains(errMsg, "check constraint") || strings.Contains(errMsg, "violates") {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("validation failed: %v", err))
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create product: %v", err))
	}

	return &insurancev1.CreateProductResponse{Product: product}, nil
}

func (s *InsuranceService) GetProduct(ctx context.Context, req *insurancev1.GetProductRequest) (*insurancev1.GetProductResponse, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product, err := s.productRepo.GetByID(ctx, req.ProductId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		logger.Error("Failed to get product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get product")
	}

	return &insurancev1.GetProductResponse{Product: product}, nil
}

func (s *InsuranceService) UpdateProduct(ctx context.Context, req *insurancev1.UpdateProductRequest) (*insurancev1.UpdateProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	product, err := s.productRepo.Update(ctx, req.Product)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update product")
	}

	return &insurancev1.UpdateProductResponse{Product: product}, nil
}

func (s *InsuranceService) DeleteProduct(ctx context.Context, req *insurancev1.DeleteProductRequest) (*emptypb.Empty, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	err := s.productRepo.Delete(ctx, req.ProductId)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete product")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListProducts(ctx context.Context, req *insurancev1.ListProductsRequest) (*insurancev1.ListProductsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	products, total, err := s.productRepo.List(ctx, req.TenantId, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list products", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	return &insurancev1.ListProductsResponse{
		Products: products,
		Total:    int32(total),
	}, nil
}

// ========== PRODUCT PLAN CRUD ==========

func (s *InsuranceService) CreateProductPlan(ctx context.Context, req *insurancev1.CreateProductPlanRequest) (*insurancev1.CreateProductPlanResponse, error) {
	if req.Plan == nil {
		return nil, status.Error(codes.InvalidArgument, "plan is required")
	}

	plan, err := s.productPlanRepo.Create(ctx, req.Plan)
	if err != nil {
		logger.Error("Failed to create product plan", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create product plan")
	}

	return &insurancev1.CreateProductPlanResponse{Plan: plan}, nil
}

func (s *InsuranceService) GetProductPlan(ctx context.Context, req *insurancev1.GetProductPlanRequest) (*insurancev1.GetProductPlanResponse, error) {
	if req.PlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "plan_id is required")
	}

	plan, err := s.productPlanRepo.GetByID(ctx, req.PlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "product plan not found")
		}
		logger.Error("Failed to get product plan", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get product plan")
	}

	return &insurancev1.GetProductPlanResponse{Plan: plan}, nil
}

func (s *InsuranceService) ListProductPlans(ctx context.Context, req *insurancev1.ListProductPlansRequest) (*insurancev1.ListProductPlansResponse, error) {
	plans, err := s.productPlanRepo.ListByProductID(ctx, req.ProductId)
	if err != nil {
		logger.Error("Failed to list product plans", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list product plans")
	}

	return &insurancev1.ListProductPlansResponse{Plans: plans}, nil
}

// ========== RIDER CRUD ==========

func (s *InsuranceService) CreateRider(ctx context.Context, req *insurancev1.CreateRiderRequest) (*insurancev1.CreateRiderResponse, error) {
	if req.Rider == nil {
		return nil, status.Error(codes.InvalidArgument, "rider is required")
	}

	rider, err := s.riderRepo.Create(ctx, req.Rider)
	if err != nil {
		logger.Error("Failed to create rider", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create rider")
	}

	return &insurancev1.CreateRiderResponse{Rider: rider}, nil
}

func (s *InsuranceService) GetRider(ctx context.Context, req *insurancev1.GetRiderRequest) (*insurancev1.GetRiderResponse, error) {
	if req.RiderId == "" {
		return nil, status.Error(codes.InvalidArgument, "rider_id is required")
	}

	rider, err := s.riderRepo.GetByID(ctx, req.RiderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "rider not found")
		}
		logger.Error("Failed to get rider", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get rider")
	}

	return &insurancev1.GetRiderResponse{Rider: rider}, nil
}

func (s *InsuranceService) ListRiders(ctx context.Context, req *insurancev1.ListRidersRequest) (*insurancev1.ListRidersResponse, error) {
	riders, err := s.riderRepo.ListByProductID(ctx, req.ProductId)
	if err != nil {
		logger.Error("Failed to list riders", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list riders")
	}

	return &insurancev1.ListRidersResponse{Riders: riders}, nil
}

// ========== PRICING CONFIG CRUD ==========

func (s *InsuranceService) CreatePricingConfig(ctx context.Context, req *insurancev1.CreatePricingConfigRequest) (*insurancev1.CreatePricingConfigResponse, error) {
	if req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "config is required")
	}

	config, err := s.pricingRepo.Create(ctx, req.Config)
	if err != nil {
		logger.Error("Failed to create pricing config", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create pricing config")
	}

	return &insurancev1.CreatePricingConfigResponse{Config: config}, nil
}

func (s *InsuranceService) GetPricingConfig(ctx context.Context, req *insurancev1.GetPricingConfigRequest) (*insurancev1.GetPricingConfigResponse, error) {
	if req.ConfigId == "" {
		return nil, status.Error(codes.InvalidArgument, "config_id is required")
	}

	config, err := s.pricingRepo.GetByID(ctx, req.ConfigId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "pricing config not found")
		}
		logger.Error("Failed to get pricing config", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get pricing config")
	}

	return &insurancev1.GetPricingConfigResponse{Config: config}, nil
}

// ========== POLICY CRUD ==========

func (s *InsuranceService) CreatePolicy(ctx context.Context, req *insurancev1.CreatePolicyRequest) (*insurancev1.CreatePolicyResponse, error) {
	if req.Policy == nil {
		return nil, status.Error(codes.InvalidArgument, "policy is required")
	}

	policy, err := s.policyRepo.Create(ctx, req.Policy)
	if err != nil {
		logger.Error("Failed to create policy", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create policy")
	}

	return &insurancev1.CreatePolicyResponse{Policy: policy}, nil
}

func (s *InsuranceService) GetPolicy(ctx context.Context, req *insurancev1.GetPolicyRequest) (*insurancev1.GetPolicyResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}

	policy, err := s.policyRepo.GetByID(ctx, req.PolicyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "policy not found")
		}
		logger.Error("Failed to get policy", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get policy")
	}

	return &insurancev1.GetPolicyResponse{Policy: policy}, nil
}

func (s *InsuranceService) UpdatePolicy(ctx context.Context, req *insurancev1.UpdatePolicyRequest) (*insurancev1.UpdatePolicyResponse, error) {
	if req.Policy == nil {
		return nil, status.Error(codes.InvalidArgument, "policy is required")
	}

	policy, err := s.policyRepo.Update(ctx, req.Policy)
	if err != nil {
		logger.Error("Failed to update policy", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update policy")
	}

	return &insurancev1.UpdatePolicyResponse{Policy: policy}, nil
}

func (s *InsuranceService) ListPolicies(ctx context.Context, req *insurancev1.ListPoliciesRequest) (*insurancev1.ListPoliciesResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	policies, total, err := s.policyRepo.List(ctx, req.TenantId, req.CustomerId, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list policies", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list policies")
	}

	return &insurancev1.ListPoliciesResponse{
		Policies: policies,
		Total:    int32(total),
	}, nil
}

// ========== CLAIM CRUD ==========

func (s *InsuranceService) CreateClaim(ctx context.Context, req *insurancev1.CreateClaimRequest) (*insurancev1.CreateClaimResponse, error) {
	if req.Claim == nil {
		return nil, status.Error(codes.InvalidArgument, "claim is required")
	}

	claim, err := s.claimRepo.Create(ctx, req.Claim)
	if err != nil {
		logger.Error("Failed to create claim", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create claim")
	}

	return &insurancev1.CreateClaimResponse{Claim: claim}, nil
}

func (s *InsuranceService) GetClaim(ctx context.Context, req *insurancev1.GetClaimRequest) (*insurancev1.GetClaimResponse, error) {
	if req.ClaimId == "" {
		return nil, status.Error(codes.InvalidArgument, "claim_id is required")
	}

	claim, err := s.claimRepo.GetByID(ctx, req.ClaimId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "claim not found")
		}
		logger.Error("Failed to get claim", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get claim")
	}

	return &insurancev1.GetClaimResponse{Claim: claim}, nil
}

func (s *InsuranceService) UpdateClaim(ctx context.Context, req *insurancev1.UpdateClaimRequest) (*insurancev1.UpdateClaimResponse, error) {
	if req.Claim == nil {
		return nil, status.Error(codes.InvalidArgument, "claim is required")
	}

	claim, err := s.claimRepo.Update(ctx, req.Claim)
	if err != nil {
		logger.Error("Failed to update claim", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update claim")
	}

	return &insurancev1.UpdateClaimResponse{Claim: claim}, nil
}

func (s *InsuranceService) ListClaims(ctx context.Context, req *insurancev1.ListClaimsRequest) (*insurancev1.ListClaimsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	claims, total, err := s.claimRepo.List(ctx, req.PolicyId, req.CustomerId, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list claims", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list claims")
	}

	return &insurancev1.ListClaimsResponse{
		Claims: claims,
		Total:  int32(total),
	}, nil
}

func (s *InsuranceService) DeletePolicy(ctx context.Context, req *insurancev1.DeletePolicyRequest) (*emptypb.Empty, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}

	err := s.policyRepo.Delete(ctx, req.PolicyId)
	if err != nil {
		logger.Error("Failed to delete policy", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete policy")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) DeleteClaim(ctx context.Context, req *insurancev1.DeleteClaimRequest) (*emptypb.Empty, error) {
	if req.ClaimId == "" {
		return nil, status.Error(codes.InvalidArgument, "claim_id is required")
	}

	err := s.claimRepo.Delete(ctx, req.ClaimId)
	if err != nil {
		logger.Error("Failed to delete claim", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete claim")
	}

	return &emptypb.Empty{}, nil
}

// ========== QUOTE CRUD ==========

func (s *InsuranceService) CreateQuote(ctx context.Context, req *insurancev1.CreateQuoteRequest) (*insurancev1.CreateQuoteResponse, error) {
	if req.Quote == nil {
		return nil, status.Error(codes.InvalidArgument, "quote is required")
	}

	quote, err := s.quoteRepo.Create(ctx, req.Quote)
	if err != nil {
		logger.Error("Failed to create quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create quote")
	}

	return &insurancev1.CreateQuoteResponse{Quote: quote}, nil
}

func (s *InsuranceService) GetQuote(ctx context.Context, req *insurancev1.GetQuoteRequest) (*insurancev1.GetQuoteResponse, error) {
	if req.QuoteId == "" {
		return nil, status.Error(codes.InvalidArgument, "quote_id is required")
	}

	quote, err := s.quoteRepo.GetByID(ctx, req.QuoteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "quote not found")
		}
		logger.Error("Failed to get quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get quote")
	}

	return &insurancev1.GetQuoteResponse{Quote: quote}, nil
}

func (s *InsuranceService) UpdateQuote(ctx context.Context, req *insurancev1.UpdateQuoteRequest) (*insurancev1.UpdateQuoteResponse, error) {
	if req.Quote == nil {
		return nil, status.Error(codes.InvalidArgument, "quote is required")
	}

	quote, err := s.quoteRepo.Update(ctx, req.Quote)
	if err != nil {
		logger.Error("Failed to update quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update quote")
	}

	return &insurancev1.UpdateQuoteResponse{Quote: quote}, nil
}

func (s *InsuranceService) DeleteQuote(ctx context.Context, req *insurancev1.DeleteQuoteRequest) (*emptypb.Empty, error) {
	if req.QuoteId == "" {
		return nil, status.Error(codes.InvalidArgument, "quote_id is required")
	}

	err := s.quoteRepo.Delete(ctx, req.QuoteId)
	if err != nil {
		logger.Error("Failed to delete quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete quote")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListQuotes(ctx context.Context, req *insurancev1.ListQuotesRequest) (*insurancev1.ListQuotesResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	quotes, total, err := s.quoteRepo.List(ctx, req.BeneficiaryId, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list quotes", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list quotes")
	}

	return &insurancev1.ListQuotesResponse{
		Quotes: quotes,
		Total:  int32(total),
	}, nil
}

// ========== UNDERWRITING DECISION CRUD ==========

func (s *InsuranceService) CreateUnderwritingDecision(ctx context.Context, req *insurancev1.CreateUnderwritingDecisionRequest) (*insurancev1.CreateUnderwritingDecisionResponse, error) {
	if req.Decision == nil {
		return nil, status.Error(codes.InvalidArgument, "decision is required")
	}

	decision, err := s.underwritingDecisionRepo.Create(ctx, req.Decision)
	if err != nil {
		logger.Error("Failed to create underwriting decision", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create underwriting decision")
	}

	return &insurancev1.CreateUnderwritingDecisionResponse{Decision: decision}, nil
}

func (s *InsuranceService) GetUnderwritingDecision(ctx context.Context, req *insurancev1.GetUnderwritingDecisionRequest) (*insurancev1.GetUnderwritingDecisionResponse, error) {
	if req.DecisionId == "" {
		return nil, status.Error(codes.InvalidArgument, "decision_id is required")
	}

	decision, err := s.underwritingDecisionRepo.GetByID(ctx, req.DecisionId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "underwriting decision not found")
		}
		logger.Error("Failed to get underwriting decision", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get underwriting decision")
	}

	return &insurancev1.GetUnderwritingDecisionResponse{Decision: decision}, nil
}

func (s *InsuranceService) UpdateUnderwritingDecision(ctx context.Context, req *insurancev1.UpdateUnderwritingDecisionRequest) (*insurancev1.UpdateUnderwritingDecisionResponse, error) {
	if req.Decision == nil {
		return nil, status.Error(codes.InvalidArgument, "decision is required")
	}

	decision, err := s.underwritingDecisionRepo.Update(ctx, req.Decision)
	if err != nil {
		logger.Error("Failed to update underwriting decision", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update underwriting decision")
	}

	return &insurancev1.UpdateUnderwritingDecisionResponse{Decision: decision}, nil
}

func (s *InsuranceService) DeleteUnderwritingDecision(ctx context.Context, req *insurancev1.DeleteUnderwritingDecisionRequest) (*emptypb.Empty, error) {
	if req.DecisionId == "" {
		return nil, status.Error(codes.InvalidArgument, "decision_id is required")
	}

	err := s.underwritingDecisionRepo.Delete(ctx, req.DecisionId)
	if err != nil {
		logger.Error("Failed to delete underwriting decision", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete underwriting decision")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListUnderwritingDecisions(ctx context.Context, req *insurancev1.ListUnderwritingDecisionsRequest) (*insurancev1.ListUnderwritingDecisionsResponse, error) {
	if req.QuoteId == "" {
		return nil, status.Error(codes.InvalidArgument, "quote_id is required")
	}

	decisions, err := s.underwritingDecisionRepo.ListByQuoteID(ctx, req.QuoteId)
	if err != nil {
		logger.Error("Failed to list underwriting decisions", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list underwriting decisions")
	}

	return &insurancev1.ListUnderwritingDecisionsResponse{Decisions: decisions}, nil
}

// ========== HEALTH DECLARATION CRUD ==========

func (s *InsuranceService) CreateHealthDeclaration(ctx context.Context, req *insurancev1.CreateHealthDeclarationRequest) (*insurancev1.CreateHealthDeclarationResponse, error) {
	if req.Declaration == nil {
		return nil, status.Error(codes.InvalidArgument, "declaration is required")
	}

	declaration, err := s.healthDeclarationRepo.Create(ctx, req.Declaration)
	if err != nil {
		logger.Error("Failed to create health declaration", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create health declaration")
	}

	return &insurancev1.CreateHealthDeclarationResponse{Declaration: declaration}, nil
}

func (s *InsuranceService) GetHealthDeclaration(ctx context.Context, req *insurancev1.GetHealthDeclarationRequest) (*insurancev1.GetHealthDeclarationResponse, error) {
	if req.DeclarationId == "" {
		return nil, status.Error(codes.InvalidArgument, "declaration_id is required")
	}

	declaration, err := s.healthDeclarationRepo.GetByID(ctx, req.DeclarationId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "health declaration not found")
		}
		logger.Error("Failed to get health declaration", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get health declaration")
	}

	return &insurancev1.GetHealthDeclarationResponse{Declaration: declaration}, nil
}

func (s *InsuranceService) UpdateHealthDeclaration(ctx context.Context, req *insurancev1.UpdateHealthDeclarationRequest) (*insurancev1.UpdateHealthDeclarationResponse, error) {
	if req.Declaration == nil {
		return nil, status.Error(codes.InvalidArgument, "declaration is required")
	}

	declaration, err := s.healthDeclarationRepo.Update(ctx, req.Declaration)
	if err != nil {
		logger.Error("Failed to update health declaration", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update health declaration")
	}

	return &insurancev1.UpdateHealthDeclarationResponse{Declaration: declaration}, nil
}

func (s *InsuranceService) DeleteHealthDeclaration(ctx context.Context, req *insurancev1.DeleteHealthDeclarationRequest) (*emptypb.Empty, error) {
	if req.DeclarationId == "" {
		return nil, status.Error(codes.InvalidArgument, "declaration_id is required")
	}

	err := s.healthDeclarationRepo.Delete(ctx, req.DeclarationId)
	if err != nil {
		logger.Error("Failed to delete health declaration", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete health declaration")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) GetHealthDeclarationByQuote(ctx context.Context, req *insurancev1.GetHealthDeclarationByQuoteRequest) (*insurancev1.GetHealthDeclarationByQuoteResponse, error) {
	if req.QuoteId == "" {
		return nil, status.Error(codes.InvalidArgument, "quote_id is required")
	}

	declaration, err := s.healthDeclarationRepo.GetByQuoteID(ctx, req.QuoteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "health declaration not found")
		}
		logger.Error("Failed to get health declaration by quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get health declaration by quote")
	}

	return &insurancev1.GetHealthDeclarationByQuoteResponse{Declaration: declaration}, nil
}

// ========== RENEWAL SCHEDULE CRUD ==========

func (s *InsuranceService) CreateRenewalSchedule(ctx context.Context, req *insurancev1.CreateRenewalScheduleRequest) (*insurancev1.CreateRenewalScheduleResponse, error) {
	if req.Schedule == nil {
		return nil, status.Error(codes.InvalidArgument, "schedule is required")
	}

	schedule, err := s.renewalScheduleRepo.Create(ctx, req.Schedule)
	if err != nil {
		logger.Error("Failed to create renewal schedule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create renewal schedule")
	}

	return &insurancev1.CreateRenewalScheduleResponse{Schedule: schedule}, nil
}

func (s *InsuranceService) GetRenewalSchedule(ctx context.Context, req *insurancev1.GetRenewalScheduleRequest) (*insurancev1.GetRenewalScheduleResponse, error) {
	if req.ScheduleId == "" {
		return nil, status.Error(codes.InvalidArgument, "schedule_id is required")
	}

	schedule, err := s.renewalScheduleRepo.GetByID(ctx, req.ScheduleId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "renewal schedule not found")
		}
		logger.Error("Failed to get renewal schedule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get renewal schedule")
	}

	return &insurancev1.GetRenewalScheduleResponse{Schedule: schedule}, nil
}

func (s *InsuranceService) UpdateRenewalSchedule(ctx context.Context, req *insurancev1.UpdateRenewalScheduleRequest) (*insurancev1.UpdateRenewalScheduleResponse, error) {
	if req.Schedule == nil {
		return nil, status.Error(codes.InvalidArgument, "schedule is required")
	}

	schedule, err := s.renewalScheduleRepo.Update(ctx, req.Schedule)
	if err != nil {
		logger.Error("Failed to update renewal schedule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update renewal schedule")
	}

	return &insurancev1.UpdateRenewalScheduleResponse{Schedule: schedule}, nil
}

func (s *InsuranceService) DeleteRenewalSchedule(ctx context.Context, req *insurancev1.DeleteRenewalScheduleRequest) (*emptypb.Empty, error) {
	if req.ScheduleId == "" {
		return nil, status.Error(codes.InvalidArgument, "schedule_id is required")
	}

	err := s.renewalScheduleRepo.Delete(ctx, req.ScheduleId)
	if err != nil {
		logger.Error("Failed to delete renewal schedule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete renewal schedule")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListRenewalSchedules(ctx context.Context, req *insurancev1.ListRenewalSchedulesRequest) (*insurancev1.ListRenewalSchedulesResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}

	schedules, err := s.renewalScheduleRepo.ListByPolicyID(ctx, req.PolicyId)
	if err != nil {
		logger.Error("Failed to list renewal schedules", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list renewal schedules")
	}

	return &insurancev1.ListRenewalSchedulesResponse{Schedules: schedules}, nil
}

// ========== RENEWAL REMINDER CRUD ==========

func (s *InsuranceService) CreateRenewalReminder(ctx context.Context, req *insurancev1.CreateRenewalReminderRequest) (*insurancev1.CreateRenewalReminderResponse, error) {
	if req.Reminder == nil {
		return nil, status.Error(codes.InvalidArgument, "reminder is required")
	}

	reminder, err := s.renewalReminderRepo.Create(ctx, req.Reminder)
	if err != nil {
		logger.Error("Failed to create renewal reminder", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create renewal reminder")
	}

	return &insurancev1.CreateRenewalReminderResponse{Reminder: reminder}, nil
}

func (s *InsuranceService) GetRenewalReminder(ctx context.Context, req *insurancev1.GetRenewalReminderRequest) (*insurancev1.GetRenewalReminderResponse, error) {
	if req.ReminderId == "" {
		return nil, status.Error(codes.InvalidArgument, "reminder_id is required")
	}

	reminder, err := s.renewalReminderRepo.GetByID(ctx, req.ReminderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "renewal reminder not found")
		}
		logger.Error("Failed to get renewal reminder", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get renewal reminder")
	}

	return &insurancev1.GetRenewalReminderResponse{Reminder: reminder}, nil
}

func (s *InsuranceService) UpdateRenewalReminder(ctx context.Context, req *insurancev1.UpdateRenewalReminderRequest) (*insurancev1.UpdateRenewalReminderResponse, error) {
	if req.Reminder == nil {
		return nil, status.Error(codes.InvalidArgument, "reminder is required")
	}

	reminder, err := s.renewalReminderRepo.Update(ctx, req.Reminder)
	if err != nil {
		logger.Error("Failed to update renewal reminder", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update renewal reminder")
	}

	return &insurancev1.UpdateRenewalReminderResponse{Reminder: reminder}, nil
}

func (s *InsuranceService) DeleteRenewalReminder(ctx context.Context, req *insurancev1.DeleteRenewalReminderRequest) (*emptypb.Empty, error) {
	if req.ReminderId == "" {
		return nil, status.Error(codes.InvalidArgument, "reminder_id is required")
	}

	err := s.renewalReminderRepo.Delete(ctx, req.ReminderId)
	if err != nil {
		logger.Error("Failed to delete renewal reminder", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete renewal reminder")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListRenewalReminders(ctx context.Context, req *insurancev1.ListRenewalRemindersRequest) (*insurancev1.ListRenewalRemindersResponse, error) {
	if req.ScheduleId == "" {
		return nil, status.Error(codes.InvalidArgument, "schedule_id is required")
	}

	reminders, err := s.renewalReminderRepo.ListByScheduleID(ctx, req.ScheduleId)
	if err != nil {
		logger.Error("Failed to list renewal reminders", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list renewal reminders")
	}

	return &insurancev1.ListRenewalRemindersResponse{Reminders: reminders}, nil
}

// ========== GRACE PERIOD CRUD ==========

func (s *InsuranceService) CreateGracePeriod(ctx context.Context, req *insurancev1.CreateGracePeriodRequest) (*insurancev1.CreateGracePeriodResponse, error) {
	if req.GracePeriod == nil {
		return nil, status.Error(codes.InvalidArgument, "grace_period is required")
	}

	gracePeriod, err := s.gracePeriodRepo.Create(ctx, req.GracePeriod)
	if err != nil {
		logger.Error("Failed to create grace period", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create grace period")
	}

	return &insurancev1.CreateGracePeriodResponse{GracePeriod: gracePeriod}, nil
}

func (s *InsuranceService) GetGracePeriod(ctx context.Context, req *insurancev1.GetGracePeriodRequest) (*insurancev1.GetGracePeriodResponse, error) {
	if req.GracePeriodId == "" {
		return nil, status.Error(codes.InvalidArgument, "grace_period_id is required")
	}

	gracePeriod, err := s.gracePeriodRepo.GetByID(ctx, req.GracePeriodId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "grace period not found")
		}
		logger.Error("Failed to get grace period", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get grace period")
	}

	return &insurancev1.GetGracePeriodResponse{GracePeriod: gracePeriod}, nil
}

func (s *InsuranceService) UpdateGracePeriod(ctx context.Context, req *insurancev1.UpdateGracePeriodRequest) (*insurancev1.UpdateGracePeriodResponse, error) {
	if req.GracePeriod == nil {
		return nil, status.Error(codes.InvalidArgument, "grace_period is required")
	}

	gracePeriod, err := s.gracePeriodRepo.Update(ctx, req.GracePeriod)
	if err != nil {
		logger.Error("Failed to update grace period", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update grace period")
	}

	return &insurancev1.UpdateGracePeriodResponse{GracePeriod: gracePeriod}, nil
}

func (s *InsuranceService) DeleteGracePeriod(ctx context.Context, req *insurancev1.DeleteGracePeriodRequest) (*emptypb.Empty, error) {
	if req.GracePeriodId == "" {
		return nil, status.Error(codes.InvalidArgument, "grace_period_id is required")
	}

	err := s.gracePeriodRepo.Delete(ctx, req.GracePeriodId)
	if err != nil {
		logger.Error("Failed to delete grace period", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete grace period")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) GetGracePeriodByPolicy(ctx context.Context, req *insurancev1.GetGracePeriodByPolicyRequest) (*insurancev1.GetGracePeriodByPolicyResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}

	gracePeriod, err := s.gracePeriodRepo.GetByPolicyID(ctx, req.PolicyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "grace period not found")
		}
		logger.Error("Failed to get grace period by policy", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get grace period by policy")
	}

	return &insurancev1.GetGracePeriodByPolicyResponse{GracePeriod: gracePeriod}, nil
}

func (s *InsuranceService) ListActiveGracePeriods(ctx context.Context, req *insurancev1.ListActiveGracePeriodsRequest) (*insurancev1.ListActiveGracePeriodsResponse, error) {
	gracePeriods, err := s.gracePeriodRepo.ListActive(ctx)
	if err != nil {
		logger.Error("Failed to list active grace periods", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list active grace periods")
	}

	return &insurancev1.ListActiveGracePeriodsResponse{GracePeriods: gracePeriods}, nil
}

// ========== INSURER CRUD ==========

func (s *InsuranceService) CreateInsurer(ctx context.Context, req *insurancev1.CreateInsurerRequest) (*insurancev1.CreateInsurerResponse, error) {
	if req.Insurer == nil {
		return nil, status.Error(codes.InvalidArgument, "insurer is required")
	}

	insurer, err := s.insurerRepo.Create(ctx, req.Insurer)
	if err != nil {
		logger.Error("Failed to create insurer", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create insurer")
	}

	return &insurancev1.CreateInsurerResponse{Insurer: insurer}, nil
}

func (s *InsuranceService) GetInsurer(ctx context.Context, req *insurancev1.GetInsurerRequest) (*insurancev1.GetInsurerResponse, error) {
	if req.InsurerId == "" {
		return nil, status.Error(codes.InvalidArgument, "insurer_id is required")
	}

	insurer, err := s.insurerRepo.GetByID(ctx, req.InsurerId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "insurer not found")
		}
		logger.Error("Failed to get insurer", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get insurer")
	}

	return &insurancev1.GetInsurerResponse{Insurer: insurer}, nil
}

func (s *InsuranceService) UpdateInsurer(ctx context.Context, req *insurancev1.UpdateInsurerRequest) (*insurancev1.UpdateInsurerResponse, error) {
	if req.Insurer == nil {
		return nil, status.Error(codes.InvalidArgument, "insurer is required")
	}

	insurer, err := s.insurerRepo.Update(ctx, req.Insurer)
	if err != nil {
		logger.Error("Failed to update insurer", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update insurer")
	}

	return &insurancev1.UpdateInsurerResponse{Insurer: insurer}, nil
}

func (s *InsuranceService) DeleteInsurer(ctx context.Context, req *insurancev1.DeleteInsurerRequest) (*emptypb.Empty, error) {
	if req.InsurerId == "" {
		return nil, status.Error(codes.InvalidArgument, "insurer_id is required")
	}

	err := s.insurerRepo.Delete(ctx, req.InsurerId)
	if err != nil {
		logger.Error("Failed to delete insurer", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete insurer")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListInsurers(ctx context.Context, req *insurancev1.ListInsurersRequest) (*insurancev1.ListInsurersResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	insurers, total, err := s.insurerRepo.List(ctx, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list insurers", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list insurers")
	}

	return &insurancev1.ListInsurersResponse{
		Insurers: insurers,
		Total:    int32(total),
	}, nil
}

// ========== INSURER CONFIG CRUD ==========

func (s *InsuranceService) CreateInsurerConfig(ctx context.Context, req *insurancev1.CreateInsurerConfigRequest) (*insurancev1.CreateInsurerConfigResponse, error) {
	if req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "config is required")
	}

	config, err := s.insurerConfigRepo.Create(ctx, req.Config)
	if err != nil {
		logger.Error("Failed to create insurer config", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create insurer config")
	}

	return &insurancev1.CreateInsurerConfigResponse{Config: config}, nil
}

func (s *InsuranceService) GetInsurerConfig(ctx context.Context, req *insurancev1.GetInsurerConfigRequest) (*insurancev1.GetInsurerConfigResponse, error) {
	if req.ConfigId == "" {
		return nil, status.Error(codes.InvalidArgument, "config_id is required")
	}

	config, err := s.insurerConfigRepo.GetByID(ctx, req.ConfigId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "insurer config not found")
		}
		logger.Error("Failed to get insurer config", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get insurer config")
	}

	return &insurancev1.GetInsurerConfigResponse{Config: config}, nil
}

func (s *InsuranceService) UpdateInsurerConfig(ctx context.Context, req *insurancev1.UpdateInsurerConfigRequest) (*insurancev1.UpdateInsurerConfigResponse, error) {
	if req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "config is required")
	}

	config, err := s.insurerConfigRepo.Update(ctx, req.Config)
	if err != nil {
		logger.Error("Failed to update insurer config", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update insurer config")
	}

	return &insurancev1.UpdateInsurerConfigResponse{Config: config}, nil
}

func (s *InsuranceService) DeleteInsurerConfig(ctx context.Context, req *insurancev1.DeleteInsurerConfigRequest) (*emptypb.Empty, error) {
	if req.ConfigId == "" {
		return nil, status.Error(codes.InvalidArgument, "config_id is required")
	}

	err := s.insurerConfigRepo.Delete(ctx, req.ConfigId)
	if err != nil {
		logger.Error("Failed to delete insurer config", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete insurer config")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) GetInsurerConfigByInsurer(ctx context.Context, req *insurancev1.GetInsurerConfigByInsurerRequest) (*insurancev1.GetInsurerConfigByInsurerResponse, error) {
	if req.InsurerId == "" {
		return nil, status.Error(codes.InvalidArgument, "insurer_id is required")
	}

	config, err := s.insurerConfigRepo.GetByInsurerID(ctx, req.InsurerId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "insurer config not found")
		}
		logger.Error("Failed to get insurer config by insurer", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get insurer config by insurer")
	}

	return &insurancev1.GetInsurerConfigByInsurerResponse{Config: config}, nil
}

// ========== INSURER PRODUCT CRUD ==========

func (s *InsuranceService) CreateInsurerProduct(ctx context.Context, req *insurancev1.CreateInsurerProductRequest) (*insurancev1.CreateInsurerProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	product, err := s.insurerProductRepo.Create(ctx, req.Product)
	if err != nil {
		logger.Error("Failed to create insurer product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create insurer product")
	}

	return &insurancev1.CreateInsurerProductResponse{Product: product}, nil
}

func (s *InsuranceService) GetInsurerProduct(ctx context.Context, req *insurancev1.GetInsurerProductRequest) (*insurancev1.GetInsurerProductResponse, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product, err := s.insurerProductRepo.GetByID(ctx, req.ProductId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "insurer product not found")
		}
		logger.Error("Failed to get insurer product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get insurer product")
	}

	return &insurancev1.GetInsurerProductResponse{Product: product}, nil
}

func (s *InsuranceService) UpdateInsurerProduct(ctx context.Context, req *insurancev1.UpdateInsurerProductRequest) (*insurancev1.UpdateInsurerProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	product, err := s.insurerProductRepo.Update(ctx, req.Product)
	if err != nil {
		logger.Error("Failed to update insurer product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update insurer product")
	}

	return &insurancev1.UpdateInsurerProductResponse{Product: product}, nil
}

func (s *InsuranceService) DeleteInsurerProduct(ctx context.Context, req *insurancev1.DeleteInsurerProductRequest) (*emptypb.Empty, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	err := s.insurerProductRepo.Delete(ctx, req.ProductId)
	if err != nil {
		logger.Error("Failed to delete insurer product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete insurer product")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListInsurerProducts(ctx context.Context, req *insurancev1.ListInsurerProductsRequest) (*insurancev1.ListInsurerProductsResponse, error) {
	if req.InsurerId == "" {
		return nil, status.Error(codes.InvalidArgument, "insurer_id is required")
	}

	products, err := s.insurerProductRepo.ListByInsurerID(ctx, req.InsurerId)
	if err != nil {
		logger.Error("Failed to list insurer products", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list insurer products")
	}

	return &insurancev1.ListInsurerProductsResponse{Products: products}, nil
}

// ========== FRAUD RULE CRUD ==========

func (s *InsuranceService) CreateFraudRule(ctx context.Context, req *insurancev1.CreateFraudRuleRequest) (*insurancev1.CreateFraudRuleResponse, error) {
	if req.Rule == nil {
		return nil, status.Error(codes.InvalidArgument, "rule is required")
	}

	rule, err := s.fraudRuleRepo.Create(ctx, req.Rule)
	if err != nil {
		logger.Error("Failed to create fraud rule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create fraud rule")
	}

	return &insurancev1.CreateFraudRuleResponse{Rule: rule}, nil
}

func (s *InsuranceService) GetFraudRule(ctx context.Context, req *insurancev1.GetFraudRuleRequest) (*insurancev1.GetFraudRuleResponse, error) {
	if req.FraudRuleId == "" {
		return nil, status.Error(codes.InvalidArgument, "fraud_rule_id is required")
	}

	rule, err := s.fraudRuleRepo.GetByID(ctx, req.FraudRuleId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "fraud rule not found")
		}
		logger.Error("Failed to get fraud rule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get fraud rule")
	}

	return &insurancev1.GetFraudRuleResponse{Rule: rule}, nil
}

func (s *InsuranceService) UpdateFraudRule(ctx context.Context, req *insurancev1.UpdateFraudRuleRequest) (*insurancev1.UpdateFraudRuleResponse, error) {
	if req.Rule == nil {
		return nil, status.Error(codes.InvalidArgument, "rule is required")
	}

	rule, err := s.fraudRuleRepo.Update(ctx, req.Rule)
	if err != nil {
		logger.Error("Failed to update fraud rule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update fraud rule")
	}

	return &insurancev1.UpdateFraudRuleResponse{Rule: rule}, nil
}

func (s *InsuranceService) DeleteFraudRule(ctx context.Context, req *insurancev1.DeleteFraudRuleRequest) (*emptypb.Empty, error) {
	if req.FraudRuleId == "" {
		return nil, status.Error(codes.InvalidArgument, "fraud_rule_id is required")
	}

	err := s.fraudRuleRepo.Delete(ctx, req.FraudRuleId)
	if err != nil {
		logger.Error("Failed to delete fraud rule", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete fraud rule")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListFraudRules(ctx context.Context, req *insurancev1.ListFraudRulesRequest) (*insurancev1.ListFraudRulesResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	rules, total, err := s.fraudRuleRepo.List(ctx, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list fraud rules", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list fraud rules")
	}

	return &insurancev1.ListFraudRulesResponse{
		Rules: rules,
		Total: int32(total),
	}, nil
}

func (s *InsuranceService) ListActiveFraudRules(ctx context.Context, req *insurancev1.ListActiveFraudRulesRequest) (*insurancev1.ListActiveFraudRulesResponse, error) {
	rules, err := s.fraudRuleRepo.ListActive(ctx)
	if err != nil {
		logger.Error("Failed to list active fraud rules", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list active fraud rules")
	}

	return &insurancev1.ListActiveFraudRulesResponse{Rules: rules}, nil
}

// ========== FRAUD CASE CRUD ==========

func (s *InsuranceService) CreateFraudCase(ctx context.Context, req *insurancev1.CreateFraudCaseRequest) (*insurancev1.CreateFraudCaseResponse, error) {
	if req.FraudCase == nil {
		return nil, status.Error(codes.InvalidArgument, "fraud_case is required")
	}

	fraudCase, err := s.fraudCaseRepo.Create(ctx, req.FraudCase)
	if err != nil {
		logger.Error("Failed to create fraud case", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create fraud case")
	}

	return &insurancev1.CreateFraudCaseResponse{FraudCase: fraudCase}, nil
}

func (s *InsuranceService) GetFraudCase(ctx context.Context, req *insurancev1.GetFraudCaseRequest) (*insurancev1.GetFraudCaseResponse, error) {
	if req.CaseId == "" {
		return nil, status.Error(codes.InvalidArgument, "case_id is required")
	}

	fraudCase, err := s.fraudCaseRepo.GetByID(ctx, req.CaseId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "fraud case not found")
		}
		logger.Error("Failed to get fraud case", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get fraud case")
	}

	return &insurancev1.GetFraudCaseResponse{FraudCase: fraudCase}, nil
}

func (s *InsuranceService) UpdateFraudCase(ctx context.Context, req *insurancev1.UpdateFraudCaseRequest) (*insurancev1.UpdateFraudCaseResponse, error) {
	if req.FraudCase == nil {
		return nil, status.Error(codes.InvalidArgument, "fraud_case is required")
	}

	fraudCase, err := s.fraudCaseRepo.Update(ctx, req.FraudCase)
	if err != nil {
		logger.Error("Failed to update fraud case", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update fraud case")
	}

	return &insurancev1.UpdateFraudCaseResponse{FraudCase: fraudCase}, nil
}

func (s *InsuranceService) DeleteFraudCase(ctx context.Context, req *insurancev1.DeleteFraudCaseRequest) (*emptypb.Empty, error) {
	if req.CaseId == "" {
		return nil, status.Error(codes.InvalidArgument, "case_id is required")
	}

	err := s.fraudCaseRepo.Delete(ctx, req.CaseId)
	if err != nil {
		logger.Error("Failed to delete fraud case", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete fraud case")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListFraudCasesByAlert(ctx context.Context, req *insurancev1.ListFraudCasesByAlertRequest) (*insurancev1.ListFraudCasesByAlertResponse, error) {
	if req.AlertId == "" {
		return nil, status.Error(codes.InvalidArgument, "alert_id is required")
	}

	cases, err := s.fraudCaseRepo.ListByAlertID(ctx, req.AlertId)
	if err != nil {
		logger.Error("Failed to list fraud cases", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list fraud cases")
	}

	return &insurancev1.ListFraudCasesByAlertResponse{Cases: cases}, nil
}

// ========== FRAUD ALERT CRUD ==========

func (s *InsuranceService) CreateFraudAlert(ctx context.Context, req *insurancev1.CreateFraudAlertRequest) (*insurancev1.CreateFraudAlertResponse, error) {
	if req.Alert == nil {
		return nil, status.Error(codes.InvalidArgument, "alert is required")
	}

	alert, err := s.fraudAlertRepo.Create(ctx, req.Alert)
	if err != nil {
		logger.Error("Failed to create fraud alert", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create fraud alert")
	}

	return &insurancev1.CreateFraudAlertResponse{Alert: alert}, nil
}

func (s *InsuranceService) GetFraudAlert(ctx context.Context, req *insurancev1.GetFraudAlertRequest) (*insurancev1.GetFraudAlertResponse, error) {
	if req.AlertId == "" {
		return nil, status.Error(codes.InvalidArgument, "alert_id is required")
	}

	alert, err := s.fraudAlertRepo.GetByID(ctx, req.AlertId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "fraud alert not found")
		}
		logger.Error("Failed to get fraud alert", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get fraud alert")
	}

	return &insurancev1.GetFraudAlertResponse{Alert: alert}, nil
}

func (s *InsuranceService) UpdateFraudAlert(ctx context.Context, req *insurancev1.UpdateFraudAlertRequest) (*insurancev1.UpdateFraudAlertResponse, error) {
	if req.Alert == nil {
		return nil, status.Error(codes.InvalidArgument, "alert is required")
	}

	alert, err := s.fraudAlertRepo.Update(ctx, req.Alert)
	if err != nil {
		logger.Error("Failed to update fraud alert", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update fraud alert")
	}

	return &insurancev1.UpdateFraudAlertResponse{Alert: alert}, nil
}

func (s *InsuranceService) DeleteFraudAlert(ctx context.Context, req *insurancev1.DeleteFraudAlertRequest) (*emptypb.Empty, error) {
	if req.AlertId == "" {
		return nil, status.Error(codes.InvalidArgument, "alert_id is required")
	}

	err := s.fraudAlertRepo.Delete(ctx, req.AlertId)
	if err != nil {
		logger.Error("Failed to delete fraud alert", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete fraud alert")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListFraudAlertsByEntity(ctx context.Context, req *insurancev1.ListFraudAlertsByEntityRequest) (*insurancev1.ListFraudAlertsByEntityResponse, error) {
	if req.EntityId == "" {
		return nil, status.Error(codes.InvalidArgument, "entity_id is required")
	}

	alerts, err := s.fraudAlertRepo.ListByEntityID(ctx, req.EntityId)
	if err != nil {
		logger.Error("Failed to list fraud alerts", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list fraud alerts")
	}

	return &insurancev1.ListFraudAlertsByEntityResponse{Alerts: alerts}, nil
}

// ========== BENEFICIARY CRUD ==========

func (s *InsuranceService) CreateBeneficiary(ctx context.Context, req *insurancev1.CreateBeneficiaryRequest) (*insurancev1.CreateBeneficiaryResponse, error) {
	if req.Beneficiary == nil {
		return nil, status.Error(codes.InvalidArgument, "beneficiary is required")
	}

	beneficiary, err := s.beneficiaryRepo.Create(ctx, req.Beneficiary)
	if err != nil {
		logger.Error("Failed to create beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create beneficiary")
	}

	return &insurancev1.CreateBeneficiaryResponse{Beneficiary: beneficiary}, nil
}

func (s *InsuranceService) GetBeneficiary(ctx context.Context, req *insurancev1.GetBeneficiaryRequest) (*insurancev1.GetBeneficiaryResponse, error) {
	if req.BeneficiaryId == "" {
		return nil, status.Error(codes.InvalidArgument, "beneficiary_id is required")
	}

	beneficiary, err := s.beneficiaryRepo.GetByID(ctx, req.BeneficiaryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "beneficiary not found")
		}
		logger.Error("Failed to get beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get beneficiary")
	}

	return &insurancev1.GetBeneficiaryResponse{Beneficiary: beneficiary}, nil
}

func (s *InsuranceService) UpdateBeneficiary(ctx context.Context, req *insurancev1.UpdateBeneficiaryRequest) (*insurancev1.UpdateBeneficiaryResponse, error) {
	if req.Beneficiary == nil {
		return nil, status.Error(codes.InvalidArgument, "beneficiary is required")
	}

	beneficiary, err := s.beneficiaryRepo.Update(ctx, req.Beneficiary)
	if err != nil {
		logger.Error("Failed to update beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update beneficiary")
	}

	return &insurancev1.UpdateBeneficiaryResponse{Beneficiary: beneficiary}, nil
}

func (s *InsuranceService) DeleteBeneficiary(ctx context.Context, req *insurancev1.DeleteBeneficiaryRequest) (*emptypb.Empty, error) {
	if req.BeneficiaryId == "" {
		return nil, status.Error(codes.InvalidArgument, "beneficiary_id is required")
	}

	err := s.beneficiaryRepo.Delete(ctx, req.BeneficiaryId)
	if err != nil {
		logger.Error("Failed to delete beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete beneficiary")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListBeneficiaries(ctx context.Context, req *insurancev1.ListBeneficiariesRequest) (*insurancev1.ListBeneficiariesResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	beneficiaries, total, err := s.beneficiaryRepo.List(ctx, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list beneficiaries", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list beneficiaries")
	}

	return &insurancev1.ListBeneficiariesResponse{
		Beneficiaries: beneficiaries,
		Total:         int32(total),
	}, nil
}

// ========== INDIVIDUAL BENEFICIARY CRUD ==========

func (s *InsuranceService) CreateIndividualBeneficiary(ctx context.Context, req *insurancev1.CreateIndividualBeneficiaryRequest) (*insurancev1.CreateIndividualBeneficiaryResponse, error) {
	if req.Individual == nil {
		return nil, status.Error(codes.InvalidArgument, "individual is required")
	}

	individual, err := s.individualBeneficiaryRepo.Create(ctx, req.Individual)
	if err != nil {
		logger.Error("Failed to create individual beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create individual beneficiary")
	}

	return &insurancev1.CreateIndividualBeneficiaryResponse{Individual: individual}, nil
}

func (s *InsuranceService) GetIndividualBeneficiary(ctx context.Context, req *insurancev1.GetIndividualBeneficiaryRequest) (*insurancev1.GetIndividualBeneficiaryResponse, error) {
	if req.BeneficiaryId == "" {
		return nil, status.Error(codes.InvalidArgument, "beneficiary_id is required")
	}

	individual, err := s.individualBeneficiaryRepo.GetByID(ctx, req.BeneficiaryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "individual beneficiary not found")
		}
		logger.Error("Failed to get individual beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get individual beneficiary")
	}

	return &insurancev1.GetIndividualBeneficiaryResponse{Individual: individual}, nil
}

func (s *InsuranceService) UpdateIndividualBeneficiary(ctx context.Context, req *insurancev1.UpdateIndividualBeneficiaryRequest) (*insurancev1.UpdateIndividualBeneficiaryResponse, error) {
	if req.Individual == nil {
		return nil, status.Error(codes.InvalidArgument, "individual is required")
	}

	individual, err := s.individualBeneficiaryRepo.Update(ctx, req.Individual)
	if err != nil {
		logger.Error("Failed to update individual beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update individual beneficiary")
	}

	return &insurancev1.UpdateIndividualBeneficiaryResponse{Individual: individual}, nil
}

func (s *InsuranceService) DeleteIndividualBeneficiary(ctx context.Context, req *insurancev1.DeleteIndividualBeneficiaryRequest) (*emptypb.Empty, error) {
	if req.BeneficiaryId == "" {
		return nil, status.Error(codes.InvalidArgument, "beneficiary_id is required")
	}

	err := s.individualBeneficiaryRepo.Delete(ctx, req.BeneficiaryId)
	if err != nil {
		logger.Error("Failed to delete individual beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete individual beneficiary")
	}

	return &emptypb.Empty{}, nil
}

// ========== BUSINESS BENEFICIARY CRUD ==========

func (s *InsuranceService) CreateBusinessBeneficiary(ctx context.Context, req *insurancev1.CreateBusinessBeneficiaryRequest) (*insurancev1.CreateBusinessBeneficiaryResponse, error) {
	if req.Business == nil {
		return nil, status.Error(codes.InvalidArgument, "business is required")
	}

	business, err := s.businessBeneficiaryRepo.Create(ctx, req.Business)
	if err != nil {
		logger.Error("Failed to create business beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create business beneficiary")
	}

	return &insurancev1.CreateBusinessBeneficiaryResponse{Business: business}, nil
}

func (s *InsuranceService) GetBusinessBeneficiary(ctx context.Context, req *insurancev1.GetBusinessBeneficiaryRequest) (*insurancev1.GetBusinessBeneficiaryResponse, error) {
	if req.BeneficiaryId == "" {
		return nil, status.Error(codes.InvalidArgument, "beneficiary_id is required")
	}

	business, err := s.businessBeneficiaryRepo.GetByID(ctx, req.BeneficiaryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "business beneficiary not found")
		}
		logger.Error("Failed to get business beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get business beneficiary")
	}

	return &insurancev1.GetBusinessBeneficiaryResponse{Business: business}, nil
}

func (s *InsuranceService) UpdateBusinessBeneficiary(ctx context.Context, req *insurancev1.UpdateBusinessBeneficiaryRequest) (*insurancev1.UpdateBusinessBeneficiaryResponse, error) {
	if req.Business == nil {
		return nil, status.Error(codes.InvalidArgument, "business is required")
	}

	business, err := s.businessBeneficiaryRepo.Update(ctx, req.Business)
	if err != nil {
		logger.Error("Failed to update business beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update business beneficiary")
	}

	return &insurancev1.UpdateBusinessBeneficiaryResponse{Business: business}, nil
}

func (s *InsuranceService) DeleteBusinessBeneficiary(ctx context.Context, req *insurancev1.DeleteBusinessBeneficiaryRequest) (*emptypb.Empty, error) {
	if req.BeneficiaryId == "" {
		return nil, status.Error(codes.InvalidArgument, "beneficiary_id is required")
	}

	err := s.businessBeneficiaryRepo.Delete(ctx, req.BeneficiaryId)
	if err != nil {
		logger.Error("Failed to delete business beneficiary", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete business beneficiary")
	}

	return &emptypb.Empty{}, nil
}

// ========== ENDORSEMENT CRUD ==========

func (s *InsuranceService) CreateEndorsement(ctx context.Context, req *insurancev1.CreateEndorsementRequest) (*insurancev1.CreateEndorsementResponse, error) {
	if req.Endorsement == nil {
		return nil, status.Error(codes.InvalidArgument, "endorsement is required")
	}

	endorsement, err := s.endorsementRepo.Create(ctx, req.Endorsement)
	if err != nil {
		logger.Error("Failed to create endorsement", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create endorsement")
	}

	return &insurancev1.CreateEndorsementResponse{Endorsement: endorsement}, nil
}

func (s *InsuranceService) GetEndorsement(ctx context.Context, req *insurancev1.GetEndorsementRequest) (*insurancev1.GetEndorsementResponse, error) {
	if req.EndorsementId == "" {
		return nil, status.Error(codes.InvalidArgument, "endorsement_id is required")
	}

	endorsement, err := s.endorsementRepo.GetByID(ctx, req.EndorsementId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "endorsement not found")
		}
		logger.Error("Failed to get endorsement", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get endorsement")
	}

	return &insurancev1.GetEndorsementResponse{Endorsement: endorsement}, nil
}

func (s *InsuranceService) UpdateEndorsement(ctx context.Context, req *insurancev1.UpdateEndorsementRequest) (*insurancev1.UpdateEndorsementResponse, error) {
	if req.Endorsement == nil {
		return nil, status.Error(codes.InvalidArgument, "endorsement is required")
	}

	endorsement, err := s.endorsementRepo.Update(ctx, req.Endorsement)
	if err != nil {
		logger.Error("Failed to update endorsement", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update endorsement")
	}

	return &insurancev1.UpdateEndorsementResponse{Endorsement: endorsement}, nil
}

func (s *InsuranceService) DeleteEndorsement(ctx context.Context, req *insurancev1.DeleteEndorsementRequest) (*emptypb.Empty, error) {
	if req.EndorsementId == "" {
		return nil, status.Error(codes.InvalidArgument, "endorsement_id is required")
	}

	err := s.endorsementRepo.Delete(ctx, req.EndorsementId)
	if err != nil {
		logger.Error("Failed to delete endorsement", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete endorsement")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListEndorsementsByPolicy(ctx context.Context, req *insurancev1.ListEndorsementsByPolicyRequest) (*insurancev1.ListEndorsementsByPolicyResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}

	endorsements, err := s.endorsementRepo.ListByPolicyID(ctx, req.PolicyId)
	if err != nil {
		logger.Error("Failed to list endorsements", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list endorsements")
	}

	return &insurancev1.ListEndorsementsByPolicyResponse{Endorsements: endorsements}, nil
}

// ========== QUOTATION CRUD ==========

func (s *InsuranceService) CreateQuotation(ctx context.Context, req *insurancev1.CreateQuotationRequest) (*insurancev1.CreateQuotationResponse, error) {
	if req.Quotation == nil {
		return nil, status.Error(codes.InvalidArgument, "quotation is required")
	}

	quotation, err := s.quotationRepo.Create(ctx, req.Quotation)
	if err != nil {
		logger.Error("Failed to create quotation", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create quotation")
	}

	return &insurancev1.CreateQuotationResponse{Quotation: quotation}, nil
}

func (s *InsuranceService) GetQuotation(ctx context.Context, req *insurancev1.GetQuotationRequest) (*insurancev1.GetQuotationResponse, error) {
	if req.QuotationId == "" {
		return nil, status.Error(codes.InvalidArgument, "quotation_id is required")
	}

	quotation, err := s.quotationRepo.GetByID(ctx, req.QuotationId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "quotation not found")
		}
		logger.Error("Failed to get quotation", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get quotation")
	}

	return &insurancev1.GetQuotationResponse{Quotation: quotation}, nil
}

func (s *InsuranceService) UpdateQuotation(ctx context.Context, req *insurancev1.UpdateQuotationRequest) (*insurancev1.UpdateQuotationResponse, error) {
	if req.Quotation == nil {
		return nil, status.Error(codes.InvalidArgument, "quotation is required")
	}

	quotation, err := s.quotationRepo.Update(ctx, req.Quotation)
	if err != nil {
		logger.Error("Failed to update quotation", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update quotation")
	}

	return &insurancev1.UpdateQuotationResponse{Quotation: quotation}, nil
}

func (s *InsuranceService) DeleteQuotation(ctx context.Context, req *insurancev1.DeleteQuotationRequest) (*emptypb.Empty, error) {
	if req.QuotationId == "" {
		return nil, status.Error(codes.InvalidArgument, "quotation_id is required")
	}

	err := s.quotationRepo.Delete(ctx, req.QuotationId)
	if err != nil {
		logger.Error("Failed to delete quotation", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete quotation")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListQuotations(ctx context.Context, req *insurancev1.ListQuotationsRequest) (*insurancev1.ListQuotationsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	quotations, total, err := s.quotationRepo.List(ctx, req.BusinessId, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list quotations", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list quotations")
	}

	return &insurancev1.ListQuotationsResponse{
		Quotations: quotations,
		Total:      int32(total),
	}, nil
}

// ========== POLICY SERVICE REQUEST CRUD ==========

func (s *InsuranceService) CreatePolicyServiceRequest(ctx context.Context, req *insurancev1.CreatePolicyServiceRequestRequest) (*insurancev1.CreatePolicyServiceRequestResponse, error) {
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	request, err := s.policyServiceRequestRepo.Create(ctx, req.Request)
	if err != nil {
		logger.Error("Failed to create policy service request", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create policy service request")
	}

	return &insurancev1.CreatePolicyServiceRequestResponse{Request: request}, nil
}

func (s *InsuranceService) GetPolicyServiceRequest(ctx context.Context, req *insurancev1.GetPolicyServiceRequestRequest) (*insurancev1.GetPolicyServiceRequestResponse, error) {
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "request_id is required")
	}

	request, err := s.policyServiceRequestRepo.GetByID(ctx, req.RequestId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "policy service request not found")
		}
		logger.Error("Failed to get policy service request", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get policy service request")
	}

	return &insurancev1.GetPolicyServiceRequestResponse{Request: request}, nil
}

func (s *InsuranceService) UpdatePolicyServiceRequest(ctx context.Context, req *insurancev1.UpdatePolicyServiceRequestRequest) (*insurancev1.UpdatePolicyServiceRequestResponse, error) {
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	request, err := s.policyServiceRequestRepo.Update(ctx, req.Request)
	if err != nil {
		logger.Error("Failed to update policy service request", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update policy service request")
	}

	return &insurancev1.UpdatePolicyServiceRequestResponse{Request: request}, nil
}

func (s *InsuranceService) DeletePolicyServiceRequest(ctx context.Context, req *insurancev1.DeletePolicyServiceRequestRequest) (*emptypb.Empty, error) {
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "request_id is required")
	}

	err := s.policyServiceRequestRepo.Delete(ctx, req.RequestId)
	if err != nil {
		logger.Error("Failed to delete policy service request", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete policy service request")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListPolicyServiceRequestsByPolicy(ctx context.Context, req *insurancev1.ListPolicyServiceRequestsByPolicyRequest) (*insurancev1.ListPolicyServiceRequestsByPolicyResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}

	requests, err := s.policyServiceRequestRepo.ListByPolicyID(ctx, req.PolicyId)
	if err != nil {
		logger.Error("Failed to list policy service requests", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list policy service requests")
	}

	return &insurancev1.ListPolicyServiceRequestsByPolicyResponse{Requests: requests}, nil
}

// ========== SERVICE PROVIDER CRUD ==========

func (s *InsuranceService) CreateServiceProvider(ctx context.Context, req *insurancev1.CreateServiceProviderRequest) (*insurancev1.CreateServiceProviderResponse, error) {
	if req.Provider == nil {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}

	provider, err := s.serviceProviderRepo.Create(ctx, req.Provider)
	if err != nil {
		logger.Error("Failed to create service provider", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create service provider")
	}

	return &insurancev1.CreateServiceProviderResponse{Provider: provider}, nil
}

func (s *InsuranceService) GetServiceProvider(ctx context.Context, req *insurancev1.GetServiceProviderRequest) (*insurancev1.GetServiceProviderResponse, error) {
	if req.ProviderId == "" {
		return nil, status.Error(codes.InvalidArgument, "provider_id is required")
	}

	provider, err := s.serviceProviderRepo.GetByID(ctx, req.ProviderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "service provider not found")
		}
		logger.Error("Failed to get service provider", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get service provider")
	}

	return &insurancev1.GetServiceProviderResponse{Provider: provider}, nil
}

func (s *InsuranceService) UpdateServiceProvider(ctx context.Context, req *insurancev1.UpdateServiceProviderRequest) (*insurancev1.UpdateServiceProviderResponse, error) {
	if req.Provider == nil {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}

	provider, err := s.serviceProviderRepo.Update(ctx, req.Provider)
	if err != nil {
		logger.Error("Failed to update service provider", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update service provider")
	}

	return &insurancev1.UpdateServiceProviderResponse{Provider: provider}, nil
}

func (s *InsuranceService) DeleteServiceProvider(ctx context.Context, req *insurancev1.DeleteServiceProviderRequest) (*emptypb.Empty, error) {
	if req.ProviderId == "" {
		return nil, status.Error(codes.InvalidArgument, "provider_id is required")
	}

	err := s.serviceProviderRepo.Delete(ctx, req.ProviderId)
	if err != nil {
		logger.Error("Failed to delete service provider", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete service provider")
	}

	return &emptypb.Empty{}, nil
}

func (s *InsuranceService) ListServiceProviders(ctx context.Context, req *insurancev1.ListServiceProvidersRequest) (*insurancev1.ListServiceProvidersResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	providers, total, err := s.serviceProviderRepo.List(ctx, req.ProviderType, req.City, int(page), int(pageSize))
	if err != nil {
		logger.Error("Failed to list service providers", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list service providers")
	}

	return &insurancev1.ListServiceProvidersResponse{
		Providers: providers,
		Total:     int32(total),
	}, nil
}
