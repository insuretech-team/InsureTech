package service

import (
	"context"
	"math"
	"strings"
	"time"

	actuarialv1 "github.com/newage-saint/insuretech/gen/go/insuretech/actuarial/entity/v1"
	actuarialservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/actuarial/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	actuarialCalculationTable = "actuarial_schema.actuarial_calculations"
	ratingFormulaTable        = "actuarial_schema.rating_formulas"
	reserveCalculationTable   = "actuarial_schema.reserve_calculations"
	lossRatioTable            = "actuarial_schema.loss_ratio_calculations"
)

type ActuarialService struct {
	actuarialservicev1.UnimplementedActuarialServiceServer
	db *gorm.DB
}

func NewActuarialService(db *gorm.DB) *ActuarialService {
	return &ActuarialService{db: db}
}

func (s *ActuarialService) CalculatePremium(ctx context.Context, req *actuarialservicev1.CalculatePremiumRequest) (*actuarialservicev1.CalculatePremiumResponse, error) {
	result := calculateActuarialPremium(req.GetInput())
	reference := req.GetCalculationReference()
	if strings.TrimSpace(reference) == "" {
		reference = actuarialReference("ACT")
	}

	calculation := buildPremiumCalculationRecord(req.GetInput(), reference, req.GetCalculatedBy(), result)
	if err := s.db.WithContext(ctx).Table(actuarialCalculationTable).Create(calculation).Error; err != nil {
		return &actuarialservicev1.CalculatePremiumResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.CalculatePremiumResponse{
		CalculationId:        calculation.GetCalculationId(),
		CalculationReference: reference,
		Result:               result,
		Success:              true,
		CalculatedAt:         calculation.GetCalculatedAt(),
	}, nil
}

func (s *ActuarialService) CalculatePurePremium(ctx context.Context, req *actuarialservicev1.CalculatePurePremiumRequest) (*actuarialservicev1.CalculatePremiumResponse, error) {
	purePremium := 0.0
	if req.GetExposureUnits() > 0 {
		purePremium = (req.GetExpectedClaims() * req.GetClaimSeverity() / req.GetExposureUnits()) * req.GetRiskAdjustmentFactor()
	}

	result := &actuarialv1.PremiumCalculationResult{
		BasePremium: purePremium,
		NetPremium:  purePremium,
		GrossPremium: purePremium,
		Currency:    "BDT",
	}

	reference := req.GetCalculationReference()
	if strings.TrimSpace(reference) == "" {
		reference = actuarialReference("ACT")
	}

	input := &actuarialv1.PremiumCalculationInput{
		ProductId:         req.GetProductId(),
		CoverageType:      "PURE_PREMIUM",
		RatingFactors:     map[string]float64{"risk_adjustment_factor": req.GetRiskAdjustmentFactor()},
		RiskCharacteristics: map[string]string{},
	}

	calculation := buildPremiumCalculationRecord(input, reference, req.GetCalculatedBy(), result)
	if err := s.db.WithContext(ctx).Table(actuarialCalculationTable).Create(calculation).Error; err != nil {
		return &actuarialservicev1.CalculatePremiumResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.CalculatePremiumResponse{
		CalculationId:        calculation.GetCalculationId(),
		CalculationReference: reference,
		Result:               result,
		Success:              true,
		CalculatedAt:         calculation.GetCalculatedAt(),
	}, nil
}

func (s *ActuarialService) CalculateReserves(ctx context.Context, req *actuarialservicev1.CalculateReservesRequest) (*actuarialservicev1.CalculateReservesResponse, error) {
	result := calculateReserveResult(req.GetInput())
	now := nowTS()
	reserve := &actuarialv1.ReserveCalculation{
		ReserveId:          newID(),
		ClaimId:            req.GetClaimId(),
		PolicyId:           "",
		ReserveType:        actuarialv1.ReserveType_RESERVE_TYPE_TOTAL,
		CaseReserve:        money(int64(math.Round(result.GetCaseReserve())), "BDT"),
		IbnrReserve:        money(int64(math.Round(result.GetIbnrReserve())), "BDT"),
		ExpenseReserve:     money(int64(math.Round(result.GetExpenseReserve())), "BDT"),
		TotalReserve:       money(int64(math.Round(result.GetTotalReserve())), "BDT"),
		CalculationMethod:  req.GetInput().GetCalculationMethod(),
		ConfidenceLevel:    req.GetInput().GetConfidenceLevel(),
		LowerBound:         money(int64(math.Round(result.GetLowerBound())), "BDT"),
		UpperBound:         money(int64(math.Round(result.GetUpperBound())), "BDT"),
		TriangleDataJson:   req.GetInput().GetTriangleDataJson(),
		Status:             actuarialv1.ReserveStatus_RESERVE_STATUS_CALCULATED,
		Metadata:           map[string]string{},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.db.WithContext(ctx).Table(reserveCalculationTable).Create(reserve).Error; err != nil {
		return &actuarialservicev1.CalculateReservesResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.CalculateReservesResponse{
		ReserveId:            reserve.GetReserveId(),
		CalculationReference: req.GetCalculationReference(),
		Result:               result,
		Success:              true,
		CalculatedAt:         now,
	}, nil
}

func (s *ActuarialService) GetReserveCalculation(ctx context.Context, req *actuarialservicev1.GetReserveCalculationRequest) (*actuarialservicev1.GetReserveCalculationResponse, error) {
	query := s.db.WithContext(ctx).Table(reserveCalculationTable).Where("deleted_at IS NULL")
	if id := strings.TrimSpace(req.GetReserveId()); id != "" {
		query = query.Where("reserve_id = ?", id)
	} else if id := strings.TrimSpace(req.GetClaimId()); id != "" {
		query = query.Where("claim_id = ?", id)
	}

	var reserve actuarialv1.ReserveCalculation
	if err := query.First(&reserve).Error; err != nil {
		return &actuarialservicev1.GetReserveCalculationResponse{Found: false}, nil
	}

	return &actuarialservicev1.GetReserveCalculationResponse{
		Reserve: &reserve,
		Found:   true,
	}, nil
}

func (s *ActuarialService) SaveReserveCalculation(ctx context.Context, req *actuarialservicev1.SaveReserveCalculationRequest) (*actuarialservicev1.SaveReserveCalculationResponse, error) {
	if req.GetReserve() == nil {
		return &actuarialservicev1.SaveReserveCalculationResponse{
			Success: false,
			Errors:  []string{"reserve is required"},
		}, nil
	}

	reserve := req.GetReserve()
	now := nowTS()
	if strings.TrimSpace(reserve.GetReserveId()) == "" {
		reserve.ReserveId = newID()
		reserve.CreatedAt = now
	}
	reserve.UpdatedAt = now

	if err := s.db.WithContext(ctx).Table(reserveCalculationTable).Save(reserve).Error; err != nil {
		return &actuarialservicev1.SaveReserveCalculationResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.SaveReserveCalculationResponse{
		Reserve: reserve,
		Success: true,
	}, nil
}

func (s *ActuarialService) CalculateLossRatio(ctx context.Context, req *actuarialservicev1.CalculateLossRatioRequest) (*actuarialservicev1.CalculateLossRatioResponse, error) {
	result := calculateLossRatioResult(req.GetInput())
	now := nowTS()
	record := &actuarialv1.LossRatioCalculation{
		LossRatioId:             newID(),
		ProductId:               req.GetProductId(),
		LineOfBusiness:          req.GetLineOfBusiness(),
		PeriodStart:             req.GetPeriodStart(),
		PeriodEnd:               req.GetPeriodEnd(),
		EarnedPremium:           money(int64(math.Round(req.GetInput().GetEarnedPremium())), "BDT"),
		WrittenPremium:          money(int64(math.Round(req.GetInput().GetWrittenPremium())), "BDT"),
		IncurredLosses:          money(int64(math.Round(req.GetInput().GetIncurredLosses())), "BDT"),
		LossAdjustmentExpenses:  money(int64(math.Round(req.GetInput().GetLossAdjustmentExpenses())), "BDT"),
		TotalIncurred:           money(int64(math.Round(req.GetInput().GetIncurredLosses()+req.GetInput().GetLossAdjustmentExpenses())), "BDT"),
		LossRatio:               result.GetLossRatio(),
		ExpenseRatio:            result.GetExpenseRatio(),
		CombinedRatio:           result.GetCombinedRatio(),
		Metadata:                map[string]string{},
		CreatedAt:               now,
		UpdatedAt:               now,
	}
	if err := s.db.WithContext(ctx).Table(lossRatioTable).Create(record).Error; err != nil {
		return &actuarialservicev1.CalculateLossRatioResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	calculation := &actuarialv1.ActuarialCalculation{
		CalculationId:         newID(),
		CalculationReference:  fallbackString(req.GetCalculationReference(), actuarialReference("LR")),
		CalculationType:       actuarialv1.ActuarialCalculationType_ACTUARIAL_CALCULATION_TYPE_LOSS_RATIO,
		EntityType:            "PRODUCT",
		EntityId:              req.GetProductId(),
		ParametersJson:        marshalJSON(req.GetInput()),
		ResultsJson:           marshalJSON(result),
		Status:                actuarialv1.CalculationStatus_CALCULATION_STATUS_COMPLETED,
		LossRatio:             result.GetLossRatio(),
		CombinedRatio:         result.GetCombinedRatio(),
		CalculatedBy:          req.GetCalculatedBy(),
		CalculatedAt:          now,
		EffectiveDate:         now,
		Metadata:              map[string]string{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	_ = s.db.WithContext(ctx).Table(actuarialCalculationTable).Create(calculation).Error

	return &actuarialservicev1.CalculateLossRatioResponse{
		LossRatioId:           record.GetLossRatioId(),
		CalculationReference:  calculation.GetCalculationReference(),
		Result:                result,
		Success:               true,
		CalculatedAt:          now,
	}, nil
}

func (s *ActuarialService) GetLossRatioCalculation(ctx context.Context, req *actuarialservicev1.GetLossRatioCalculationRequest) (*actuarialservicev1.GetLossRatioCalculationResponse, error) {
	var record actuarialv1.LossRatioCalculation
	if err := s.db.WithContext(ctx).Table(lossRatioTable).
		Where("loss_ratio_id = ? AND deleted_at IS NULL", req.GetLossRatioId()).
		First(&record).Error; err != nil {
		return &actuarialservicev1.GetLossRatioCalculationResponse{Found: false}, nil
	}

	return &actuarialservicev1.GetLossRatioCalculationResponse{
		LossRatio: &record,
		Found:     true,
	}, nil
}

func (s *ActuarialService) ListLossRatioCalculations(ctx context.Context, req *actuarialservicev1.ListLossRatioCalculationsRequest) (*actuarialservicev1.ListLossRatioCalculationsResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(lossRatioTable).Where("deleted_at IS NULL")
	if productID := strings.TrimSpace(req.GetProductId()); productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if lob := strings.TrimSpace(req.GetLineOfBusiness()); lob != "" {
		query = query.Where("line_of_business = ?", lob)
	}
	if req.GetPeriodStart() != nil {
		query = query.Where("period_start >= ?", req.GetPeriodStart().AsTime())
	}
	if req.GetPeriodEnd() != nil {
		query = query.Where("period_end <= ?", req.GetPeriodEnd().AsTime())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &actuarialservicev1.ListLossRatioCalculationsResponse{}, nil
	}

	var records []*actuarialv1.LossRatioCalculation
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return &actuarialservicev1.ListLossRatioCalculationsResponse{}, nil
	}

	return &actuarialservicev1.ListLossRatioCalculationsResponse{
		LossRatios:    records,
		TotalCount:    int32(total),
		NextPageToken: nextToken(offset, len(records), int(total)),
	}, nil
}

func (s *ActuarialService) SaveLossRatioCalculation(ctx context.Context, req *actuarialservicev1.SaveLossRatioCalculationRequest) (*actuarialservicev1.SaveLossRatioCalculationResponse, error) {
	if req.GetLossRatio() == nil {
		return &actuarialservicev1.SaveLossRatioCalculationResponse{
			Success: false,
			Errors:  []string{"loss_ratio is required"},
		}, nil
	}

	record := req.GetLossRatio()
	now := nowTS()
	if strings.TrimSpace(record.GetLossRatioId()) == "" {
		record.LossRatioId = newID()
		record.CreatedAt = now
	}
	record.UpdatedAt = now

	if err := s.db.WithContext(ctx).Table(lossRatioTable).Save(record).Error; err != nil {
		return &actuarialservicev1.SaveLossRatioCalculationResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.SaveLossRatioCalculationResponse{
		LossRatio: record,
		Success:   true,
	}, nil
}

func (s *ActuarialService) CreateRatingFormula(ctx context.Context, req *actuarialservicev1.CreateRatingFormulaRequest) (*actuarialservicev1.CreateRatingFormulaResponse, error) {
	now := nowTS()
	formula := &actuarialv1.RatingFormula{
		FormulaId:         newID(),
		FormulaCode:       req.GetFormulaCode(),
		FormulaName:       req.GetFormulaName(),
		Description:       req.GetDescription(),
		Category:          req.GetCategory(),
		InsuranceType:     req.GetInsuranceType(),
		FormulaExpression: req.GetFormulaExpression(),
		VariablesJson:     marshalJSON(req.GetVariables()),
		SortOrder:         req.GetSortOrder(),
		Version:           1,
		Status:            actuarialv1.FormulaStatus_FORMULA_STATUS_DRAFT,
		ValidFrom:         fallbackTimestamp(req.GetValidFrom(), now),
		ValidUntil:        req.GetValidUntil(),
		Metadata:          ensureMetadata(req.GetMetadata()),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := s.db.WithContext(ctx).Table(ratingFormulaTable).Create(formula).Error; err != nil {
		return &actuarialservicev1.CreateRatingFormulaResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.CreateRatingFormulaResponse{
		FormulaId: formula.GetFormulaId(),
		Success:   true,
		Formula:   formula,
	}, nil
}

func (s *ActuarialService) GetRatingFormula(ctx context.Context, req *actuarialservicev1.GetRatingFormulaRequest) (*actuarialservicev1.GetRatingFormulaResponse, error) {
	query := s.db.WithContext(ctx).Table(ratingFormulaTable).Where("deleted_at IS NULL")
	if id := strings.TrimSpace(req.GetFormulaId()); id != "" {
		query = query.Where("formula_id = ?", id)
	} else if code := strings.TrimSpace(req.GetFormulaCode()); code != "" {
		query = query.Where("formula_code = ?", code)
	}

	var formula actuarialv1.RatingFormula
	if err := query.First(&formula).Error; err != nil {
		return &actuarialservicev1.GetRatingFormulaResponse{Found: false}, nil
	}

	return &actuarialservicev1.GetRatingFormulaResponse{
		Formula: &formula,
		Found:   true,
	}, nil
}

func (s *ActuarialService) UpdateRatingFormula(ctx context.Context, req *actuarialservicev1.UpdateRatingFormulaRequest) (*actuarialservicev1.UpdateRatingFormulaResponse, error) {
	var formula actuarialv1.RatingFormula
	if err := s.db.WithContext(ctx).Table(ratingFormulaTable).
		Where("formula_id = ? AND deleted_at IS NULL", req.GetFormulaId()).
		First(&formula).Error; err != nil {
		return &actuarialservicev1.UpdateRatingFormulaResponse{
			Success: false,
			Errors:  []string{"formula not found"},
		}, nil
	}

	formula.FormulaName = req.GetFormulaName()
	formula.Description = req.GetDescription()
	formula.Category = req.GetCategory()
	formula.FormulaExpression = req.GetFormulaExpression()
	formula.VariablesJson = marshalJSON(req.GetVariables())
	formula.SortOrder = req.GetSortOrder()
	formula.ValidUntil = req.GetValidUntil()
	if req.GetMetadata() != nil {
		formula.Metadata = ensureMetadata(req.GetMetadata())
	}
	formula.UpdatedAt = nowTS()

	if err := s.db.WithContext(ctx).Table(ratingFormulaTable).Where("formula_id = ?", formula.GetFormulaId()).Save(&formula).Error; err != nil {
		return &actuarialservicev1.UpdateRatingFormulaResponse{
			Success: false,
			Errors:  []string{err.Error()},
		}, nil
	}

	return &actuarialservicev1.UpdateRatingFormulaResponse{
		Formula: &formula,
		Success: true,
	}, nil
}

func (s *ActuarialService) DeleteRatingFormula(ctx context.Context, req *actuarialservicev1.DeleteRatingFormulaRequest) (*actuarialservicev1.DeleteRatingFormulaResponse, error) {
	query := s.db.WithContext(ctx).Table(ratingFormulaTable).Where("formula_id = ?", req.GetFormulaId())
	var result *gorm.DB
	if req.GetPermanent() {
		result = query.Delete(&actuarialv1.RatingFormula{})
	} else {
		result = query.Update("deleted_at", nowTS())
	}
	if result.Error != nil {
		return &actuarialservicev1.DeleteRatingFormulaResponse{Success: false}, nil
	}

	return &actuarialservicev1.DeleteRatingFormulaResponse{
		Success: true,
		Deleted: result.RowsAffected > 0,
	}, nil
}

func (s *ActuarialService) ListRatingFormulas(ctx context.Context, req *actuarialservicev1.ListRatingFormulasRequest) (*actuarialservicev1.ListRatingFormulasResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(ratingFormulaTable).Where("deleted_at IS NULL")
	if insuranceType := strings.TrimSpace(req.GetInsuranceType()); insuranceType != "" {
		query = query.Where("insurance_type = ?", insuranceType)
	}
	if req.GetCategory() != actuarialv1.FormulaCategory_FORMULA_CATEGORY_UNSPECIFIED {
		query = query.Where("category = ?", req.GetCategory())
	}
	if req.GetStatus() != actuarialv1.FormulaStatus_FORMULA_STATUS_UNSPECIFIED {
		query = query.Where("status = ?", req.GetStatus())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &actuarialservicev1.ListRatingFormulasResponse{}, nil
	}

	var formulas []*actuarialv1.RatingFormula
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&formulas).Error; err != nil {
		return &actuarialservicev1.ListRatingFormulasResponse{}, nil
	}

	return &actuarialservicev1.ListRatingFormulasResponse{
		Formulas:      formulas,
		TotalCount:    int32(total),
		NextPageToken: nextToken(offset, len(formulas), int(total)),
	}, nil
}

func (s *ActuarialService) ActivateRatingFormula(ctx context.Context, req *actuarialservicev1.ActivateRatingFormulaRequest) (*actuarialservicev1.ActivateRatingFormulaResponse, error) {
	var formula actuarialv1.RatingFormula
	if err := s.db.WithContext(ctx).Table(ratingFormulaTable).
		Where("formula_id = ? AND deleted_at IS NULL", req.GetFormulaId()).
		First(&formula).Error; err != nil {
		return &actuarialservicev1.ActivateRatingFormulaResponse{Success: false}, nil
	}

	formula.Status = actuarialv1.FormulaStatus_FORMULA_STATUS_ACTIVE
	formula.ValidFrom = fallbackTimestamp(req.GetEffectiveDate(), nowTS())
	formula.UpdatedAt = nowTS()
	if err := s.db.WithContext(ctx).Table(ratingFormulaTable).Where("formula_id = ?", formula.GetFormulaId()).Save(&formula).Error; err != nil {
		return &actuarialservicev1.ActivateRatingFormulaResponse{Success: false}, nil
	}

	return &actuarialservicev1.ActivateRatingFormulaResponse{
		Success: true,
		Formula: &formula,
	}, nil
}

func (s *ActuarialService) GetCalculation(ctx context.Context, req *actuarialservicev1.GetCalculationRequest) (*actuarialservicev1.GetCalculationResponse, error) {
	query := s.db.WithContext(ctx).Table(actuarialCalculationTable).Where("deleted_at IS NULL")
	if id := strings.TrimSpace(req.GetCalculationId()); id != "" {
		query = query.Where("calculation_id = ?", id)
	} else if ref := strings.TrimSpace(req.GetCalculationReference()); ref != "" {
		query = query.Where("calculation_reference = ?", ref)
	}

	var calculation actuarialv1.ActuarialCalculation
	if err := query.First(&calculation).Error; err != nil {
		return &actuarialservicev1.GetCalculationResponse{Found: false}, nil
	}

	return &actuarialservicev1.GetCalculationResponse{
		Calculation: &calculation,
		Found:       true,
	}, nil
}

func (s *ActuarialService) ListCalculations(ctx context.Context, req *actuarialservicev1.ListCalculationsRequest) (*actuarialservicev1.ListCalculationsResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(actuarialCalculationTable).Where("deleted_at IS NULL")
	if req.GetCalculationType() != actuarialv1.ActuarialCalculationType_ACTUARIAL_CALCULATION_TYPE_UNSPECIFIED {
		query = query.Where("calculation_type = ?", req.GetCalculationType())
	}
	if entityType := strings.TrimSpace(req.GetEntityType()); entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	if entityID := strings.TrimSpace(req.GetEntityId()); entityID != "" {
		query = query.Where("entity_id = ?", entityID)
	}
	if req.GetDateFrom() != nil {
		query = query.Where("calculated_at >= ?", req.GetDateFrom().AsTime())
	}
	if req.GetDateTo() != nil {
		query = query.Where("calculated_at <= ?", req.GetDateTo().AsTime())
	}
	if by := strings.TrimSpace(req.GetCalculatedBy()); by != "" {
		query = query.Where("calculated_by = ?", by)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &actuarialservicev1.ListCalculationsResponse{}, nil
	}

	var calculations []*actuarialv1.ActuarialCalculation
	if err := query.Order("calculated_at DESC").Offset(offset).Limit(pageSize).Find(&calculations).Error; err != nil {
		return &actuarialservicev1.ListCalculationsResponse{}, nil
	}

	return &actuarialservicev1.ListCalculationsResponse{
		Calculations:   calculations,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(calculations), int(total)),
	}, nil
}

func buildPremiumCalculationRecord(input *actuarialv1.PremiumCalculationInput, reference, calculatedBy string, result *actuarialv1.PremiumCalculationResult) *actuarialv1.ActuarialCalculation {
	now := nowTS()
	return &actuarialv1.ActuarialCalculation{
		CalculationId:        newID(),
		CalculationReference: reference,
		CalculationType:      actuarialv1.ActuarialCalculationType_ACTUARIAL_CALCULATION_TYPE_PREMIUM,
		EntityType:           "PRODUCT",
		EntityId:             input.GetProductId(),
		ParametersJson:       marshalJSON(input),
		ResultsJson:          marshalJSON(result),
		Status:               actuarialv1.CalculationStatus_CALCULATION_STATUS_COMPLETED,
		CalculatedPremium:    result.GetGrossPremium(),
		CalculatedBy:         calculatedBy,
		CalculatedAt:         now,
		EffectiveDate:        now,
		Metadata:             map[string]string{},
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

func calculateActuarialPremium(input *actuarialv1.PremiumCalculationInput) *actuarialv1.PremiumCalculationResult {
	base := input.GetSumInsured() * 0.015
	loadingsTotal := 0.0
	discountsTotal := 0.0
	breakdown := make([]*actuarialv1.FactorBreakdown, 0)

	for name, value := range input.GetRatingFactors() {
		amount := base * value * 0.01
		loadingsTotal += amount
		breakdown = append(breakdown, &actuarialv1.FactorBreakdown{
			FactorName:  name,
			FactorType:  "RATING_FACTOR",
			FactorValue: value,
			Amount:      amount,
			Description: "Rating factor adjustment",
		})
	}

	for _, loading := range input.GetLoadings() {
		amount := base * 0.02
		loadingsTotal += amount
		breakdown = append(breakdown, &actuarialv1.FactorBreakdown{
			FactorName:  loading,
			FactorType:  "LOADING",
			FactorValue: 0.02,
			Amount:      amount,
			Description: "Loading applied",
		})
	}

	for _, discount := range input.GetDiscounts() {
		amount := base * 0.01
		discountsTotal += amount
		breakdown = append(breakdown, &actuarialv1.FactorBreakdown{
			FactorName:  discount,
			FactorType:  "DISCOUNT",
			FactorValue: 0.01,
			Amount:      amount,
			Description: "Discount applied",
		})
	}

	netPremium := base + loadingsTotal
	grossPremium := math.Max(0, netPremium-discountsTotal)
	return &actuarialv1.PremiumCalculationResult{
		BasePremium:    base,
		NetPremium:     netPremium,
		GrossPremium:   grossPremium,
		TotalLoadings:  loadingsTotal,
		TotalDiscounts: discountsTotal,
		FactorBreakdown: breakdown,
		Currency:       "BDT",
	}
}

func calculateReserveResult(input *actuarialv1.ReserveInput) *actuarialv1.ReserveResult {
	caseReserve := input.GetCaseReserve()
	reportedClaims := input.GetReportedClaims()
	expectedUltimate := reportedClaims * 1.2
	ibnrReserve := math.Max(0, expectedUltimate-reportedClaims)
	ibnerReserve := caseReserve * 0.1
	expenseReserve := (caseReserve + ibnrReserve) * 0.12
	totalReserve := caseReserve + ibnrReserve + ibnerReserve + expenseReserve
	stdDev := totalReserve * 0.15
	zScore := 1.96

	return &actuarialv1.ReserveResult{
		CaseReserve:   caseReserve,
		IbnrReserve:   ibnrReserve,
		IbnerReserve:  ibnerReserve,
		ExpenseReserve: expenseReserve,
		TotalReserve:  totalReserve,
		LowerBound:    math.Max(0, totalReserve-zScore*stdDev),
		UpperBound:    totalReserve + zScore*stdDev,
		MethodUsed:    input.GetCalculationMethod(),
	}
}

func calculateLossRatioResult(input *actuarialv1.LossRatioInput) *actuarialv1.LossRatioResult {
	totalIncurred := input.GetIncurredLosses() + input.GetLossAdjustmentExpenses()
	lossRatio := 0.0
	if input.GetEarnedPremium() > 0 {
		lossRatio = totalIncurred / input.GetEarnedPremium()
	}

	expenseRatio := 0.0
	if input.GetWrittenPremium() > 0 {
		expenseRatio = input.GetOperatingExpenses() / input.GetWrittenPremium()
	}

	combinedRatio := lossRatio + expenseRatio
	interpretation := "LOSS_MAKING"
	switch {
	case combinedRatio < 0.95:
		interpretation = "PROFITABLE"
	case combinedRatio <= 1.05:
		interpretation = "BREAK_EVEN"
	}

	return &actuarialv1.LossRatioResult{
		LossRatio:                math.Round(lossRatio*10000) / 10000,
		ExpenseRatio:             math.Round(expenseRatio*10000) / 10000,
		CombinedRatio:            math.Round(combinedRatio*10000) / 10000,
		UnderwritingProfitMargin: math.Round((1-combinedRatio)*10000) / 10000,
		Interpretation:           interpretation,
	}
}

func actuarialReference(prefix string) string {
	return prefix + "-" + time.Now().UTC().Format("2006") + "-" + strings.ToUpper(newID()[:6])
}

func fallbackString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func fallbackTimestamp(value, fallback *timestamppb.Timestamp) *timestamppb.Timestamp {
	if value != nil {
		return value
	}
	return fallback
}
