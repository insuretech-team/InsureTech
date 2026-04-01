package service

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	quotingv1 "github.com/newage-saint/insuretech/gen/go/insuretech/quoting/entity/v1"
	quotingservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/quoting/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const quotingTable = "quoting_schema.quotes"

type QuotingService struct {
	quotingservicev1.UnimplementedQuotingServiceServer
	db *gorm.DB
}

func NewQuotingService(db *gorm.DB) *QuotingService {
	return &QuotingService{db: db}
}

func (s *QuotingService) GenerateQuote(ctx context.Context, req *quotingservicev1.GenerateQuoteRequest) (*quotingservicev1.GenerateQuoteResponse, error) {
	if strings.TrimSpace(req.GetProductId()) == "" || strings.TrimSpace(req.GetCustomerId()) == "" || req.GetParameters() == nil {
		return &quotingservicev1.GenerateQuoteResponse{
			Error: errorResponse("INVALID_ARGUMENT", "product_id, customer_id, and parameters are required", 400),
		}, nil
	}

	calculation, coverages, discounts := calculateQuotePremium(req.GetProductId(), req.GetParameters())
	now := time.Now().UTC()
	validityDays := req.GetValidityDays()
	if validityDays <= 0 {
		validityDays = 30
	}

	quote := &quotingv1.Quote{
		QuoteId:                newID(),
		QuoteNumber:            "QT-" + now.Format("20060102") + "-" + strings.ToUpper(newID()[:8]),
		ProductId:              req.GetProductId(),
		CustomerId:             req.GetCustomerId(),
		AgentId:                req.GetAgentId(),
		Status:                 quotingv1.QuoteStatus_QUOTE_STATUS_GENERATED,
		ParametersJson:         marshalJSON(req.GetParameters()),
		PremiumCalculationJson: marshalJSON(calculation),
		CoveragesJson:          marshalJSON(coverages),
		DiscountsJson:          marshalJSON(discounts),
		TotalPremium:           calculation.GetTotalPremium(),
		ValidFrom:              timestamppb.New(now),
		ValidUntil:             timestamppb.New(now.AddDate(0, 0, int(validityDays))),
		RevisionNumber:         1,
		Metadata:               ensureMetadata(req.GetMetadata()),
		CreatedAt:              timestamppb.New(now),
		UpdatedAt:              timestamppb.New(now),
	}

	if err := s.db.WithContext(ctx).Table(quotingTable).Create(quote).Error; err != nil {
		return &quotingservicev1.GenerateQuoteResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &quotingservicev1.GenerateQuoteResponse{Quote: quote}, nil
}

func (s *QuotingService) GetQuote(ctx context.Context, req *quotingservicev1.GetQuoteRequest) (*quotingservicev1.GetQuoteResponse, error) {
	var quote quotingv1.Quote
	err := s.db.WithContext(ctx).
		Table(quotingTable).
		Where("quote_id = ? AND deleted_at IS NULL", req.GetQuoteId()).
		First(&quote).Error
	if err != nil {
		return &quotingservicev1.GetQuoteResponse{
			Error: errorResponse("QUOTE_NOT_FOUND", "quote not found", 404),
		}, nil
	}

	return &quotingservicev1.GetQuoteResponse{Quote: &quote}, nil
}

func (s *QuotingService) GetQuoteByNumber(ctx context.Context, req *quotingservicev1.GetQuoteByNumberRequest) (*quotingservicev1.GetQuoteResponse, error) {
	var quote quotingv1.Quote
	err := s.db.WithContext(ctx).
		Table(quotingTable).
		Where("quote_number = ? AND deleted_at IS NULL", req.GetQuoteNumber()).
		First(&quote).Error
	if err != nil {
		return &quotingservicev1.GetQuoteResponse{
			Error: errorResponse("QUOTE_NOT_FOUND", "quote not found", 404),
		}, nil
	}

	return &quotingservicev1.GetQuoteResponse{Quote: &quote}, nil
}

func (s *QuotingService) ListQuotes(ctx context.Context, req *quotingservicev1.ListQuotesRequest) (*quotingservicev1.ListQuotesResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(quotingTable).Where("deleted_at IS NULL")
	if id := strings.TrimSpace(req.GetCustomerId()); id != "" {
		query = query.Where("customer_id = ?", id)
	}
	if id := strings.TrimSpace(req.GetProductId()); id != "" {
		query = query.Where("product_id = ?", id)
	}
	if req.GetStatus() != quotingv1.QuoteStatus_QUOTE_STATUS_UNSPECIFIED {
		query = query.Where("status = ?", req.GetStatus())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &quotingservicev1.ListQuotesResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var quotes []*quotingv1.Quote
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&quotes).Error; err != nil {
		return &quotingservicev1.ListQuotesResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &quotingservicev1.ListQuotesResponse{
		Quotes:         quotes,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(quotes), int(total)),
	}, nil
}

func (s *QuotingService) ReviseQuote(ctx context.Context, req *quotingservicev1.ReviseQuoteRequest) (*quotingservicev1.ReviseQuoteResponse, error) {
	getResp, _ := s.GetQuote(ctx, &quotingservicev1.GetQuoteRequest{QuoteId: req.GetQuoteId()})
	if getResp.GetError() != nil || getResp.GetQuote() == nil {
		return &quotingservicev1.ReviseQuoteResponse{Error: errorResponse("QUOTE_NOT_FOUND", "quote not found", 404)}, nil
	}

	parent := getResp.GetQuote()
	calculation, coverages, discounts := calculateQuotePremium(parent.GetProductId(), req.GetNewParameters())
	now := time.Now().UTC()
	validityDays := req.GetValidityDays()
	if validityDays <= 0 {
		validityDays = 30
	}

	parent.Status = quotingv1.QuoteStatus_QUOTE_STATUS_EXPIRED
	parent.UpdatedAt = timestamppb.New(now)
	if err := s.db.WithContext(ctx).Table(quotingTable).
		Where("quote_id = ?", parent.GetQuoteId()).
		Updates(map[string]any{
			"status":     parent.GetStatus(),
			"updated_at": parent.GetUpdatedAt(),
		}).Error; err != nil {
		return &quotingservicev1.ReviseQuoteResponse{Error: errorResponse("UPDATE_FAILED", err.Error(), 500)}, nil
	}

	revised := &quotingv1.Quote{
		QuoteId:                newID(),
		QuoteNumber:            "QT-" + now.Format("20060102") + "-" + strings.ToUpper(newID()[:8]),
		ProductId:              parent.GetProductId(),
		CustomerId:             parent.GetCustomerId(),
		AgentId:                parent.GetAgentId(),
		Status:                 quotingv1.QuoteStatus_QUOTE_STATUS_GENERATED,
		ParametersJson:         marshalJSON(req.GetNewParameters()),
		PremiumCalculationJson: marshalJSON(calculation),
		CoveragesJson:          marshalJSON(coverages),
		DiscountsJson:          marshalJSON(discounts),
		TotalPremium:           calculation.GetTotalPremium(),
		ValidFrom:              timestamppb.New(now),
		ValidUntil:             timestamppb.New(now.AddDate(0, 0, int(validityDays))),
		RevisionNumber:         parent.GetRevisionNumber() + 1,
		ParentQuoteId:          parent.GetQuoteId(),
		RevisionReason:         req.GetRevisionReason(),
		Metadata:               ensureMetadata(parent.GetMetadata()),
		CreatedAt:              timestamppb.New(now),
		UpdatedAt:              timestamppb.New(now),
	}

	if err := s.db.WithContext(ctx).Table(quotingTable).Create(revised).Error; err != nil {
		return &quotingservicev1.ReviseQuoteResponse{Error: errorResponse("CREATE_FAILED", err.Error(), 500)}, nil
	}

	return &quotingservicev1.ReviseQuoteResponse{
		Quote:       revised,
		ParentQuote: parent,
	}, nil
}

func (s *QuotingService) CompareQuotes(ctx context.Context, req *quotingservicev1.CompareQuotesRequest) (*quotingservicev1.CompareQuotesResponse, error) {
	resp := &quotingservicev1.CompareQuotesResponse{}
	for _, id := range req.GetQuoteIds() {
		getResp, _ := s.GetQuote(ctx, &quotingservicev1.GetQuoteRequest{QuoteId: id})
		if getResp.GetQuote() == nil {
			continue
		}

		quote := getResp.GetQuote()
		resp.Comparisons = append(resp.Comparisons, &quotingservicev1.QuoteComparison{
			QuoteId:       quote.GetQuoteId(),
			QuoteNumber:   quote.GetQuoteNumber(),
			TotalPremium:  quote.GetTotalPremium(),
			ValidUntil:    quote.GetValidUntil(),
			Status:        quote.GetStatus(),
		})
	}

	return resp, nil
}

func (s *QuotingService) ConvertQuoteToPolicy(ctx context.Context, req *quotingservicev1.ConvertQuoteToPolicyRequest) (*quotingservicev1.ConvertQuoteToPolicyResponse, error) {
	now := nowTS()
	updates := map[string]any{
		"status":               quotingv1.QuoteStatus_QUOTE_STATUS_CONVERTED,
		"converted_policy_id":  req.GetPolicyId(),
		"converted_at":         now,
		"updated_at":           now,
	}

	result := s.db.WithContext(ctx).Table(quotingTable).
		Where("quote_id = ? AND deleted_at IS NULL", req.GetQuoteId()).
		Updates(updates)
	if result.Error != nil {
		return &quotingservicev1.ConvertQuoteToPolicyResponse{
			Error: errorResponse("CONVERT_FAILED", result.Error.Error(), 500),
		}, nil
	}
	if result.RowsAffected == 0 {
		return &quotingservicev1.ConvertQuoteToPolicyResponse{
			Error: errorResponse("QUOTE_NOT_FOUND", "quote not found", 404),
		}, nil
	}

	return &quotingservicev1.ConvertQuoteToPolicyResponse{
		QuoteId:     req.GetQuoteId(),
		PolicyId:    req.GetPolicyId(),
		ConvertedAt: now,
	}, nil
}

func (s *QuotingService) ExpireQuote(ctx context.Context, req *quotingservicev1.ExpireQuoteRequest) (*quotingservicev1.ExpireQuoteResponse, error) {
	result := s.db.WithContext(ctx).Table(quotingTable).
		Where("quote_id = ? AND deleted_at IS NULL", req.GetQuoteId()).
		Updates(map[string]any{
			"status":     quotingv1.QuoteStatus_QUOTE_STATUS_EXPIRED,
			"updated_at": nowTS(),
		})
	if result.Error != nil {
		return &quotingservicev1.ExpireQuoteResponse{
			Error: errorResponse("EXPIRE_FAILED", result.Error.Error(), 500),
		}, nil
	}

	return &quotingservicev1.ExpireQuoteResponse{Success: result.RowsAffected > 0}, nil
}

func (s *QuotingService) DeleteQuote(ctx context.Context, req *quotingservicev1.DeleteQuoteRequest) (*quotingservicev1.DeleteQuoteResponse, error) {
	query := s.db.WithContext(ctx).Table(quotingTable).Where("quote_id = ?", req.GetQuoteId())
	var result *gorm.DB
	if req.GetPermanent() {
		result = query.Delete(map[string]any{})
	} else {
		result = query.Update("deleted_at", nowTS())
	}
	if result.Error != nil {
		return &quotingservicev1.DeleteQuoteResponse{
			Error: errorResponse("DELETE_FAILED", result.Error.Error(), 500),
		}, nil
	}

	return &quotingservicev1.DeleteQuoteResponse{Success: result.RowsAffected > 0}, nil
}

func (s *QuotingService) GetQuoteStatistics(ctx context.Context, req *quotingservicev1.GetQuoteStatisticsRequest) (*quotingservicev1.GetQuoteStatisticsResponse, error) {
	listResp, _ := s.ListQuotes(ctx, &quotingservicev1.ListQuotesRequest{
		CustomerId: req.GetCustomerId(),
		ProductId:  req.GetProductId(),
		PageSize:   500,
	})
	if listResp.GetError() != nil {
		return &quotingservicev1.GetQuoteStatisticsResponse{Error: listResp.GetError()}, nil
	}

	quotes := listResp.GetQuotes()
	filtered := make([]*quotingv1.Quote, 0, len(quotes))
	for _, quote := range quotes {
		if quote == nil {
			continue
		}
		if req.GetStartDate() != nil && quote.GetCreatedAt() != nil && quote.GetCreatedAt().AsTime().Before(req.GetStartDate().AsTime()) {
			continue
		}
		if req.GetEndDate() != nil && quote.GetCreatedAt() != nil && quote.GetCreatedAt().AsTime().After(req.GetEndDate().AsTime()) {
			continue
		}
		filtered = append(filtered, quote)
	}

	var totalPremium int64
	var converted int32
	stats := &quotingservicev1.GetQuoteStatisticsResponse{TotalQuotes: int32(len(filtered))}
	for _, quote := range filtered {
		if quote.GetTotalPremium() != nil {
			totalPremium += quote.GetTotalPremium().GetAmount()
		}
		switch quote.GetStatus() {
		case quotingv1.QuoteStatus_QUOTE_STATUS_DRAFT:
			stats.DraftQuotes++
		case quotingv1.QuoteStatus_QUOTE_STATUS_SENT:
			stats.SentQuotes++
		case quotingv1.QuoteStatus_QUOTE_STATUS_ACCEPTED:
			stats.AcceptedQuotes++
		case quotingv1.QuoteStatus_QUOTE_STATUS_DECLINED:
			stats.DeclinedQuotes++
		case quotingv1.QuoteStatus_QUOTE_STATUS_EXPIRED:
			stats.ExpiredQuotes++
		case quotingv1.QuoteStatus_QUOTE_STATUS_CONVERTED:
			stats.ConvertedQuotes++
			converted++
		}
	}

	if len(filtered) > 0 {
		stats.AveragePremium = money(totalPremium/int64(len(filtered)), "BDT")
	}
	stats.TotalPremiumValue = money(totalPremium, "BDT")
	if stats.TotalQuotes > 0 {
		stats.ConversionRate = float64(converted) / float64(stats.TotalQuotes) * 100
	}

	return stats, nil
}

func (s *QuotingService) CalculatePremium(_ context.Context, req *quotingservicev1.CalculatePremiumRequest) (*quotingservicev1.CalculatePremiumResponse, error) {
	if strings.TrimSpace(req.GetProductId()) == "" || req.GetParameters() == nil {
		return &quotingservicev1.CalculatePremiumResponse{
			Error: errorResponse("INVALID_ARGUMENT", "product_id and parameters are required", 400),
		}, nil
	}

	calculation, coverages, discounts := calculateQuotePremium(req.GetProductId(), req.GetParameters())
	return &quotingservicev1.CalculatePremiumResponse{
		Calculation: calculation,
		Coverages:   coverages,
		Discounts:   discounts,
	}, nil
}

func calculateQuotePremium(productID string, params *quotingv1.QuoteParameters) (*quotingv1.PremiumCalculation, []*quotingv1.Coverage, []*quotingv1.Discount) {
	assetValue := params.GetAssetValue()
	if assetValue <= 0 {
		assetValue = 100000
	}

	basePremiumAmount := int64(math.Round(assetValue * 0.02))
	if basePremiumAmount < 1000 {
		basePremiumAmount = 1000
	}

	riskAdjustmentAmount := int64(math.Round(float64(basePremiumAmount) * 0.10))
	optionalAmount := int64(len(params.GetOptionalCoverages()) * 250)
	discountAmount := int64(0)
	if params.GetCoverageDurationMonths() >= 12 {
		discountAmount = int64(math.Round(float64(basePremiumAmount+riskAdjustmentAmount+optionalAmount) * 0.05))
	}
	taxAmount := int64(math.Round(float64(basePremiumAmount+riskAdjustmentAmount+optionalAmount-discountAmount) * 0.15))
	totalAmount := basePremiumAmount + riskAdjustmentAmount + optionalAmount - discountAmount + taxAmount

	calculation := &quotingv1.PremiumCalculation{
		BasePremium:            money(basePremiumAmount, "BDT"),
		RiskAdjustment:         money(riskAdjustmentAmount, "BDT"),
		OptionalCoveragesTotal: money(optionalAmount, "BDT"),
		DiscountsTotal:         money(discountAmount, "BDT"),
		Taxes:                  money(taxAmount, "BDT"),
		Fees:                   money(0, "BDT"),
		TotalPremium:           money(totalAmount, "BDT"),
		Currency:               "BDT",
		Breakdown: []*quotingv1.PremiumBreakdown{
			{Category: "Base", Description: "Base premium derived from asset value", Amount: money(basePremiumAmount, "BDT")},
			{Category: "Risk", Description: "Flat risk adjustment", Amount: money(riskAdjustmentAmount, "BDT")},
			{Category: "Optional", Description: "Selected optional coverages", Amount: money(optionalAmount, "BDT")},
			{Category: "Tax", Description: "Estimated taxes", Amount: money(taxAmount, "BDT")},
		},
	}

	coverages := []*quotingv1.Coverage{
		{
			CoverageId: "base-" + productID,
			Name:       "Base Coverage",
			Description:"Standard protection",
			Limit:      money(int64(assetValue), "BDT"),
			Deductible: money(5000, "BDT"),
			Premium:    money(basePremiumAmount, "BDT"),
			IsIncluded: true,
		},
	}

	for index, coverage := range params.GetOptionalCoverages() {
		coverages = append(coverages, &quotingv1.Coverage{
			CoverageId: "opt-" + strconv.Itoa(index+1),
			Name:       coverage.GetName(),
			Description:"Optional coverage",
			Limit:      money(int64(coverage.GetSelectedLimit()), "BDT"),
			Deductible: money(int64(coverage.GetSelectedDeductible()), "BDT"),
			Premium:    money(250, "BDT"),
			IsIncluded: true,
			IsOptional: true,
		})
	}

	discounts := []*quotingv1.Discount{}
	if discountAmount > 0 {
		discounts = append(discounts, &quotingv1.Discount{
			DiscountId:   "disc-annual",
			Name:         "Annual Term Discount",
			Description:  "Applied for 12+ month coverage",
			Amount:       money(discountAmount, "BDT"),
			Percentage:   5,
			DiscountType: "TERM",
		})
	}

	return calculation, coverages, discounts
}
