package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	lifev1 "github.com/newage-saint/insuretech/gen/go/insuretech/life/entity/v1"
	lifeservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/life/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	lifeProductTable = "life_schema.life_products"
	lifeQuoteTable   = "life_schema.life_quotes"
)

type LifeInsuranceService struct {
	lifeservicev1.UnimplementedLifeInsuranceServiceServer
	db *gorm.DB
}

type lifePremiumOutcome struct {
	basePremium        int64
	ageAddition        int64
	conditionMultiplier float32
	conditionAddition  int64
	bonusDiscount      int64
	totalPremium       int64
	breakdown          []*lifev1.PremiumBreakdown
	appliedConditions  []string
	appliedBonuses     []string
}

func NewLifeInsuranceService(db *gorm.DB) *LifeInsuranceService {
	return &LifeInsuranceService{db: db}
}

func (s *LifeInsuranceService) GetLifeProduct(ctx context.Context, req *lifeservicev1.GetLifeProductRequest) (*lifeservicev1.GetLifeProductResponse, error) {
	var product lifev1.LifeProduct
	err := s.db.WithContext(ctx).
		Table(lifeProductTable).
		Where("product_id = ? AND deleted_at IS NULL", req.GetProductId()).
		First(&product).Error
	if err != nil {
		return &lifeservicev1.GetLifeProductResponse{
			Error: errorResponse("PRODUCT_NOT_FOUND", "life product not found", 404),
		}, nil
	}

	return &lifeservicev1.GetLifeProductResponse{Product: &product}, nil
}

func (s *LifeInsuranceService) ListLifeProducts(ctx context.Context, req *lifeservicev1.ListLifeProductsRequest) (*lifeservicev1.ListLifeProductsResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(lifeProductTable).Where("deleted_at IS NULL")
	if req.GetOnlyActive() {
		query = query.Where("is_active = ?", true)
	}
	if req.GetProductType() != lifev1.LifeProductType_LIFE_PRODUCT_TYPE_UNSPECIFIED {
		query = query.Where("product_type = ?", req.GetProductType())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &lifeservicev1.ListLifeProductsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var products []*lifev1.LifeProduct
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return &lifeservicev1.ListLifeProductsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &lifeservicev1.ListLifeProductsResponse{
		Products:       products,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(products), int(total)),
	}, nil
}

func (s *LifeInsuranceService) CreateLifeProduct(ctx context.Context, req *lifeservicev1.CreateLifeProductRequest) (*lifeservicev1.CreateLifeProductResponse, error) {
	if strings.TrimSpace(req.GetProductCode()) == "" || strings.TrimSpace(req.GetProductName()) == "" {
		return &lifeservicev1.CreateLifeProductResponse{
			Error: errorResponse("INVALID_ARGUMENT", "product_code and product_name are required", 400),
		}, nil
	}

	now := time.Now().UTC()
	product := &lifev1.LifeProduct{
		ProductId:              newID(),
		ProductCode:            req.GetProductCode(),
		ProductName:            req.GetProductName(),
		ProductType:            req.GetProductType(),
		Description:            req.GetDescription(),
		BaseRate:               req.GetBaseRate(),
		AgeAdditionConfig:      req.GetAgeAdditionConfig(),
		ConditionMultipliersJson: marshalJSON(req.GetConditionMultipliers()),
		BonusConfigJson:        marshalJSON(req.GetBonuses()),
		MinSumAssured:          req.GetMinSumAssured(),
		MaxSumAssured:          req.GetMaxSumAssured(),
		MinEntryAge:            req.GetMinEntryAge(),
		MaxEntryAge:            req.GetMaxEntryAge(),
		MinPolicyTerm:          req.GetMinPolicyTerm(),
		MaxPolicyTerm:          req.GetMaxPolicyTerm(),
		IsActive:               true,
		Metadata:               ensureMetadata(req.GetMetadata()),
		CreatedAt:              timestamppb.New(now),
		UpdatedAt:              timestamppb.New(now),
	}

	if product.GetAgeAdditionConfig() == nil {
		product.AgeAdditionConfig = &lifev1.AgeAdditionConfig{}
	}

	if err := s.db.WithContext(ctx).Table(lifeProductTable).Create(product).Error; err != nil {
		return &lifeservicev1.CreateLifeProductResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &lifeservicev1.CreateLifeProductResponse{Product: product}, nil
}

func (s *LifeInsuranceService) UpdateLifeProduct(ctx context.Context, req *lifeservicev1.UpdateLifeProductRequest) (*lifeservicev1.UpdateLifeProductResponse, error) {
	var product lifev1.LifeProduct
	err := s.db.WithContext(ctx).
		Table(lifeProductTable).
		Where("product_id = ? AND deleted_at IS NULL", req.GetProductId()).
		First(&product).Error
	if err != nil {
		return &lifeservicev1.UpdateLifeProductResponse{
			Error: errorResponse("PRODUCT_NOT_FOUND", "life product not found", 404),
		}, nil
	}

	product.ProductName = req.GetProductName()
	product.Description = req.GetDescription()
	product.BaseRate = req.GetBaseRate()
	if req.GetAgeAdditionConfig() != nil {
		product.AgeAdditionConfig = req.GetAgeAdditionConfig()
	}
	product.ConditionMultipliersJson = marshalJSON(req.GetConditionMultipliers())
	product.BonusConfigJson = marshalJSON(req.GetBonuses())
	product.MinSumAssured = req.GetMinSumAssured()
	product.MaxSumAssured = req.GetMaxSumAssured()
	product.MinEntryAge = req.GetMinEntryAge()
	product.MaxEntryAge = req.GetMaxEntryAge()
	product.MinPolicyTerm = req.GetMinPolicyTerm()
	product.MaxPolicyTerm = req.GetMaxPolicyTerm()
	product.IsActive = req.GetIsActive()
	if req.GetMetadata() != nil {
		product.Metadata = ensureMetadata(req.GetMetadata())
	}
	product.UpdatedAt = nowTS()

	if err := s.db.WithContext(ctx).Table(lifeProductTable).
		Where("product_id = ?", product.GetProductId()).
		Save(&product).Error; err != nil {
		return &lifeservicev1.UpdateLifeProductResponse{
			Error: errorResponse("UPDATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &lifeservicev1.UpdateLifeProductResponse{Product: &product}, nil
}

func (s *LifeInsuranceService) DeleteLifeProduct(ctx context.Context, req *lifeservicev1.DeleteLifeProductRequest) (*lifeservicev1.DeleteLifeProductResponse, error) {
	query := s.db.WithContext(ctx).Table(lifeProductTable).Where("product_id = ?", req.GetProductId())
	var result *gorm.DB
	if req.GetPermanent() {
		result = query.Delete(&lifev1.LifeProduct{})
	} else {
		result = query.Update("deleted_at", nowTS())
	}
	if result.Error != nil {
		return &lifeservicev1.DeleteLifeProductResponse{
			Error: errorResponse("DELETE_FAILED", result.Error.Error(), 500),
		}, nil
	}

	return &lifeservicev1.DeleteLifeProductResponse{Success: result.RowsAffected > 0}, nil
}

func (s *LifeInsuranceService) CalculatePremium(ctx context.Context, req *lifeservicev1.CalculatePremiumRequest) (*lifeservicev1.CalculatePremiumResponse, error) {
	product, errResp := s.getLifeProduct(ctx, req.GetProductId())
	if errResp != nil {
		return &lifeservicev1.CalculatePremiumResponse{Error: errResp}, nil
	}

	outcome := calculateLifePremium(product, req.GetInsuredPerson(), req.GetAgeAtEntry(), req.GetPolicyTermYears(), req.GetSumAssured(), req.GetBonusCodes())
	return &lifeservicev1.CalculatePremiumResponse{
		BasePremium:        outcome.basePremium,
		AgeAddition:        outcome.ageAddition,
		ConditionMultiplier: outcome.conditionMultiplier,
		ConditionAddition:  outcome.conditionAddition,
		BonusDiscount:      outcome.bonusDiscount,
		TotalPremium:       outcome.totalPremium,
		Breakdown:          outcome.breakdown,
		AppliedConditions:  outcome.appliedConditions,
		AppliedBonuses:     outcome.appliedBonuses,
	}, nil
}

func (s *LifeInsuranceService) GenerateQuote(ctx context.Context, req *lifeservicev1.GenerateQuoteRequest) (*lifeservicev1.GenerateQuoteResponse, error) {
	product, errResp := s.getLifeProduct(ctx, req.GetProductId())
	if errResp != nil {
		return &lifeservicev1.GenerateQuoteResponse{Error: errResp}, nil
	}

	outcome := calculateLifePremium(product, req.GetInsuredPerson(), req.GetAgeAtEntry(), req.GetPolicyTermYears(), req.GetSumAssured(), req.GetBonusCodes())
	now := time.Now().UTC()
	validityDays := req.GetValidityDays()
	if validityDays <= 0 {
		validityDays = 30
	}

	quote := &lifev1.LifeQuote{
		QuoteId:              newID(),
		QuoteNumber:          "LQ-" + now.Format("20060102") + "-" + strings.ToUpper(newID()[:8]),
		ProductId:            req.GetProductId(),
		CustomerId:           req.GetCustomerId(),
		AgentId:              req.GetAgentId(),
		Status:               lifev1.LifeQuoteStatus_LIFE_QUOTE_STATUS_GENERATED,
		InsuredPersonJson:    marshalJSON(req.GetInsuredPerson()),
		AgeAtEntry:           req.GetAgeAtEntry(),
		PolicyTermYears:      req.GetPolicyTermYears(),
		SumAssured:           req.GetSumAssured(),
		BasePremium:          outcome.basePremium,
		AgeAddition:          outcome.ageAddition,
		ConditionMultiplier:  outcome.conditionMultiplier,
		ConditionAddition:    outcome.conditionAddition,
		BonusDiscount:        outcome.bonusDiscount,
		TotalPremium:         outcome.totalPremium,
		HealthConditionsJson: marshalJSON(req.GetInsuredPerson().GetHealthConditions()),
		BonusesAppliedJson:   marshalJSON(outcome.appliedBonuses),
		ValidUntil:           timestamppb.New(now.AddDate(0, 0, int(validityDays))),
		Metadata:             ensureMetadata(req.GetMetadata()),
		CreatedAt:            timestamppb.New(now),
		UpdatedAt:            timestamppb.New(now),
	}

	if err := s.db.WithContext(ctx).Table(lifeQuoteTable).Create(quote).Error; err != nil {
		return &lifeservicev1.GenerateQuoteResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &lifeservicev1.GenerateQuoteResponse{Quote: quote}, nil
}

func (s *LifeInsuranceService) GetQuote(ctx context.Context, req *lifeservicev1.GetQuoteRequest) (*lifeservicev1.GetQuoteResponse, error) {
	var quote lifev1.LifeQuote
	err := s.db.WithContext(ctx).
		Table(lifeQuoteTable).
		Where("quote_id = ? AND deleted_at IS NULL", req.GetQuoteId()).
		First(&quote).Error
	if err != nil {
		return &lifeservicev1.GetQuoteResponse{
			Error: errorResponse("QUOTE_NOT_FOUND", "life quote not found", 404),
		}, nil
	}

	return &lifeservicev1.GetQuoteResponse{Quote: &quote}, nil
}

func (s *LifeInsuranceService) GetQuoteByNumber(ctx context.Context, req *lifeservicev1.GetQuoteByNumberRequest) (*lifeservicev1.GetQuoteResponse, error) {
	var quote lifev1.LifeQuote
	err := s.db.WithContext(ctx).
		Table(lifeQuoteTable).
		Where("quote_number = ? AND deleted_at IS NULL", req.GetQuoteNumber()).
		First(&quote).Error
	if err != nil {
		return &lifeservicev1.GetQuoteResponse{
			Error: errorResponse("QUOTE_NOT_FOUND", "life quote not found", 404),
		}, nil
	}

	return &lifeservicev1.GetQuoteResponse{Quote: &quote}, nil
}

func (s *LifeInsuranceService) ListQuotes(ctx context.Context, req *lifeservicev1.ListQuotesRequest) (*lifeservicev1.ListQuotesResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(lifeQuoteTable).Where("deleted_at IS NULL")
	if id := strings.TrimSpace(req.GetCustomerId()); id != "" {
		query = query.Where("customer_id = ?", id)
	}
	if id := strings.TrimSpace(req.GetProductId()); id != "" {
		query = query.Where("product_id = ?", id)
	}
	if req.GetStatus() != lifev1.LifeQuoteStatus_LIFE_QUOTE_STATUS_UNSPECIFIED {
		query = query.Where("status = ?", req.GetStatus())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &lifeservicev1.ListQuotesResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var quotes []*lifev1.LifeQuote
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&quotes).Error; err != nil {
		return &lifeservicev1.ListQuotesResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &lifeservicev1.ListQuotesResponse{
		Quotes:         quotes,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(quotes), int(total)),
	}, nil
}

func (s *LifeInsuranceService) ConvertQuoteToPolicy(ctx context.Context, req *lifeservicev1.ConvertQuoteToPolicyRequest) (*lifeservicev1.ConvertQuoteToPolicyResponse, error) {
	now := nowTS()
	result := s.db.WithContext(ctx).Table(lifeQuoteTable).
		Where("quote_id = ? AND deleted_at IS NULL", req.GetQuoteId()).
		Updates(map[string]any{
			"status":               lifev1.LifeQuoteStatus_LIFE_QUOTE_STATUS_CONVERTED,
			"converted_policy_id":  req.GetPolicyId(),
			"converted_at":         now,
			"updated_at":           now,
		})
	if result.Error != nil {
		return &lifeservicev1.ConvertQuoteToPolicyResponse{
			Error: errorResponse("CONVERT_FAILED", result.Error.Error(), 500),
		}, nil
	}
	if result.RowsAffected == 0 {
		return &lifeservicev1.ConvertQuoteToPolicyResponse{
			Error: errorResponse("QUOTE_NOT_FOUND", "life quote not found", 404),
		}, nil
	}

	return &lifeservicev1.ConvertQuoteToPolicyResponse{
		QuoteId:     req.GetQuoteId(),
		PolicyId:    req.GetPolicyId(),
		ConvertedAt: now,
	}, nil
}

func (s *LifeInsuranceService) GetHealthConditions(ctx context.Context, req *lifeservicev1.GetHealthConditionsRequest) (*lifeservicev1.GetHealthConditionsResponse, error) {
	product, errResp := s.getLifeProduct(ctx, req.GetProductId())
	if errResp != nil {
		return &lifeservicev1.GetHealthConditionsResponse{Error: errResp}, nil
	}

	var conditions []*lifev1.ConditionMultiplier
	_ = json.Unmarshal([]byte(product.GetConditionMultipliersJson()), &conditions)

	return &lifeservicev1.GetHealthConditionsResponse{Conditions: conditions}, nil
}

func (s *LifeInsuranceService) getLifeProduct(ctx context.Context, productID string) (*lifev1.LifeProduct, *commonv1.Error) {
	var product lifev1.LifeProduct
	err := s.db.WithContext(ctx).
		Table(lifeProductTable).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		First(&product).Error
	if err != nil {
		return nil, errorResponse("PRODUCT_NOT_FOUND", "life product not found", 404)
	}
	return &product, nil
}

func calculateLifePremium(product *lifev1.LifeProduct, insuredPerson *lifev1.InsuredPerson, ageAtEntry int32, policyTermYears int32, _ int64, bonusCodes []string) lifePremiumOutcome {
	var conditionMultipliers []*lifev1.ConditionMultiplier
	var bonusConfigs []*lifev1.BonusConfig
	_ = json.Unmarshal([]byte(product.GetConditionMultipliersJson()), &conditionMultipliers)
	_ = json.Unmarshal([]byte(product.GetBonusConfigJson()), &bonusConfigs)

	basePremium := product.GetBaseRate() * int64(policyTermYears)
	ageAddition := int64(0)
	ageConfig := product.GetAgeAdditionConfig()
	if ageConfig != nil && ageAtEntry > ageConfig.GetStartAge() && ageConfig.GetAgeIncrement() > 0 {
		ageDiff := ageAtEntry - ageConfig.GetStartAge()
		incrementCount := ageDiff / ageConfig.GetAgeIncrement()
		ageAddition = int64(incrementCount) * ageConfig.GetPriceToAdd() * int64(policyTermYears)
	}

	breakdown := []*lifev1.PremiumBreakdown{
		{
			Component:   "Base Premium",
			Amount:      basePremium,
			Description: "Base rate multiplied by policy term",
			IsDiscount:  false,
		},
		{
			Component:   "Age Addition",
			Amount:      ageAddition,
			Description: "Age bracket loading",
			IsDiscount:  false,
		},
	}

	totalConditionMultiplier := float32(1.0)
	appliedConditions := make([]string, 0)
	for _, condition := range insuredPerson.GetHealthConditions() {
		for _, configured := range conditionMultipliers {
			if strings.EqualFold(configured.GetConditionCode(), condition.GetConditionCode()) {
				totalConditionMultiplier += configured.GetMultiplier()
				appliedConditions = append(appliedConditions, condition.GetConditionName())
				break
			}
		}
	}

	afterAge := basePremium + ageAddition
	conditionAddition := int64(float32(afterAge) * (totalConditionMultiplier - 1.0))
	breakdown = append(breakdown, &lifev1.PremiumBreakdown{
		Component:   "Condition Addition",
		Amount:      conditionAddition,
		Description: "Health condition multiplier adjustment",
		IsDiscount:  false,
	})

	subtotal := afterAge + conditionAddition
	bonusDiscount := int64(0)
	appliedBonuses := make([]string, 0)
	for _, code := range bonusCodes {
		for _, bonus := range bonusConfigs {
			if !strings.EqualFold(bonus.GetBonusCode(), code) {
				continue
			}

			amount := int64(0)
			switch strings.ToUpper(bonus.GetBonusType()) {
			case "PERCENTAGE":
				amount = int64(float64(subtotal) * float64(bonus.GetPercentage()))
			case "FIXED_AMOUNT":
				amount = bonus.GetFixedAmount()
			}

			bonusDiscount += amount
			appliedBonuses = append(appliedBonuses, bonus.GetBonusName())
			breakdown = append(breakdown, &lifev1.PremiumBreakdown{
				Component:   bonus.GetBonusName(),
				Amount:      amount,
				Description: bonus.GetDescription(),
				IsDiscount:  true,
			})
			break
		}
	}

	totalPremium := subtotal + conditionAddition - bonusDiscount
	return lifePremiumOutcome{
		basePremium:        basePremium,
		ageAddition:        ageAddition,
		conditionMultiplier: totalConditionMultiplier,
		conditionAddition:  conditionAddition,
		bonusDiscount:      bonusDiscount,
		totalPremium:       totalPremium,
		breakdown:          breakdown,
		appliedConditions:  appliedConditions,
		appliedBonuses:     appliedBonuses,
	}
}
