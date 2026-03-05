package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/metrics"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/repository"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
)

// FraudService contains fraud domain business logic.
type FraudService struct {
	ruleRepo  domain.RuleRepository
	alertRepo domain.AlertRepository
	caseRepo  domain.CaseRepository
	publisher *events.Publisher
	metrics   *metrics.RuntimeMetrics
}

func NewFraudService(
	ruleRepo domain.RuleRepository,
	alertRepo domain.AlertRepository,
	caseRepo domain.CaseRepository,
	publisher *events.Publisher,
) *FraudService {
	return &FraudService{
		ruleRepo:  ruleRepo,
		alertRepo: alertRepo,
		caseRepo:  caseRepo,
		publisher: publisher,
		metrics:   metrics.NewRuntimeMetrics(),
	}
}

func (s *FraudService) CheckFraud(ctx context.Context, req *fraudservicev1.CheckFraudRequest) (*fraudservicev1.CheckFraudResponse, error) {
	s.metrics.IncFraudChecks()
	if strings.TrimSpace(req.EntityType) == "" || strings.TrimSpace(req.EntityId) == "" {
		return nil, fmt.Errorf("%w: entity_type and entity_id are required", ErrInvalidArgument)
	}

	rules, _, err := s.ruleRepo.List(ctx, fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED, true, 500, 0)
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	if req.Data != nil {
		data = req.Data.AsMap()
	}

	triggeredRuleNames := make([]string, 0, 8)
	totalScore := int32(0)
	primaryRuleID := ""
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		matched := evaluateRule(rule, data)
		if !matched {
			continue
		}
		if primaryRuleID == "" {
			primaryRuleID = rule.FraudRuleId
		}
		name := strings.TrimSpace(rule.Name)
		if name == "" {
			name = rule.FraudRuleId
		}
		triggeredRuleNames = append(triggeredRuleNames, name)
		totalScore += rule.ScoreWeight
	}

	riskLevel := riskLevelFromScore(totalScore)
	detected := len(triggeredRuleNames) > 0
	fraudAlertID := ""
	if detected {
		s.metrics.IncFraudDetections()
		detailsJSON, _ := json.Marshal(map[string]any{
			"triggered_rules": triggeredRuleNames,
			"entity_type":     req.EntityType,
			"entity_id":       req.EntityId,
			"data":            data,
		})
		alert := &fraudv1.FraudAlert{
			EntityType:  req.EntityType,
			EntityId:    req.EntityId,
			FraudRuleId: primaryRuleID,
			RiskLevel:   riskLevel,
			FraudScore:  totalScore,
			Details:     string(detailsJSON),
			Status:      fraudv1.AlertStatus_ALERT_STATUS_OPEN,
		}
		if err := s.alertRepo.Create(ctx, alert); err != nil {
			return nil, err
		}
		s.metrics.IncAlertsCreated()
		fraudAlertID = alert.Id
		if s.publisher != nil {
			_ = s.publisher.PublishFraudAlertTriggered(ctx, alert, correlationIDFromContext(ctx))
		}
	}

	return &fraudservicev1.CheckFraudResponse{
		IsFraudDetected: detected,
		FraudScore:      totalScore,
		RiskLevel:       riskLevel,
		TriggeredRules:  triggeredRuleNames,
		FraudAlertId:    fraudAlertID,
	}, nil
}

func (s *FraudService) GetFraudAlert(ctx context.Context, req *fraudservicev1.GetFraudAlertRequest) (*fraudservicev1.GetFraudAlertResponse, error) {
	if strings.TrimSpace(req.FraudAlertId) == "" {
		return nil, fmt.Errorf("%w: fraud_alert_id is required", ErrInvalidArgument)
	}
	alert, err := s.alertRepo.GetByID(ctx, req.FraudAlertId)
	if err != nil {
		if errors.Is(err, repository.ErrAlertNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &fraudservicev1.GetFraudAlertResponse{FraudAlert: alert}, nil
}

func (s *FraudService) ListFraudAlerts(ctx context.Context, req *fraudservicev1.ListFraudAlertsRequest) (*fraudservicev1.ListFraudAlertsResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	start, end := parseDateRange(req.StartDate, req.EndDate)
	alerts, total, err := s.alertRepo.List(ctx, req.Status, req.RiskLevel, start, end, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return &fraudservicev1.ListFraudAlertsResponse{
		FraudAlerts: alerts,
		TotalCount:  total,
	}, nil
}

func (s *FraudService) CreateFraudCase(ctx context.Context, req *fraudservicev1.CreateFraudCaseRequest) (*fraudservicev1.CreateFraudCaseResponse, error) {
	if strings.TrimSpace(req.FraudAlertId) == "" {
		return nil, fmt.Errorf("%w: fraud_alert_id is required", ErrInvalidArgument)
	}
	if _, err := s.alertRepo.GetByID(ctx, req.FraudAlertId); err != nil {
		if errors.Is(err, repository.ErrAlertNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	priority, ok := parseCasePriority(req.Priority)
	if !ok {
		return nil, fmt.Errorf("%w: invalid priority", ErrInvalidArgument)
	}

	fraudCase := &fraudv1.FraudCase{
		FraudAlertId:       req.FraudAlertId,
		Priority:           priority,
		InvestigationNotes: req.InvestigationNotes,
		Status:             fraudv1.CaseStatus_CASE_STATUS_OPEN,
		InvestigatorId:     req.InvestigatorId,
	}
	if err := s.caseRepo.Create(ctx, fraudCase); err != nil {
		return nil, err
	}
	s.metrics.IncCasesCreated()
	_ = s.alertRepo.UpdateStatus(ctx, req.FraudAlertId, fraudv1.AlertStatus_ALERT_STATUS_INVESTIGATING, req.InvestigatorId)

	if s.publisher != nil {
		_ = s.publisher.PublishFraudCaseCreated(ctx, fraudCase, correlationIDFromContext(ctx))
	}

	return &fraudservicev1.CreateFraudCaseResponse{
		FraudCaseId: fraudCase.Id,
		CaseNumber:  fraudCase.CaseNumber,
		Message:     "fraud case created",
	}, nil
}

func (s *FraudService) GetFraudCase(ctx context.Context, req *fraudservicev1.GetFraudCaseRequest) (*fraudservicev1.GetFraudCaseResponse, error) {
	if strings.TrimSpace(req.FraudCaseId) == "" {
		return nil, fmt.Errorf("%w: fraud_case_id is required", ErrInvalidArgument)
	}
	fraudCase, err := s.caseRepo.GetByID(ctx, req.FraudCaseId)
	if err != nil {
		if errors.Is(err, repository.ErrCaseNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &fraudservicev1.GetFraudCaseResponse{FraudCase: fraudCase}, nil
}

func (s *FraudService) UpdateFraudCase(ctx context.Context, req *fraudservicev1.UpdateFraudCaseRequest) (*fraudservicev1.UpdateFraudCaseResponse, error) {
	if strings.TrimSpace(req.FraudCaseId) == "" {
		return nil, fmt.Errorf("%w: fraud_case_id is required", ErrInvalidArgument)
	}

	status, ok := parseCaseStatus(req.Status)
	if !ok {
		return nil, fmt.Errorf("%w: invalid status", ErrInvalidArgument)
	}
	outcome, ok := parseCaseOutcome(req.Outcome)
	if !ok {
		return nil, fmt.Errorf("%w: invalid outcome", ErrInvalidArgument)
	}

	evidence := ""
	if req.Evidence != nil {
		b, _ := json.Marshal(req.Evidence.AsMap())
		evidence = string(b)
	}

	if err := s.caseRepo.Update(ctx, req.FraudCaseId, status, outcome, req.InvestigationNotes, evidence); err != nil {
		if errors.Is(err, repository.ErrCaseNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	updatedCase, err := s.caseRepo.GetByID(ctx, req.FraudCaseId)
	if err == nil {
		if status == fraudv1.CaseStatus_CASE_STATUS_CLOSED {
			_ = s.alertRepo.UpdateStatus(ctx, updatedCase.FraudAlertId, fraudv1.AlertStatus_ALERT_STATUS_CLOSED, updatedCase.InvestigatorId)
		}
		if outcome == fraudv1.CaseOutcome_CASE_OUTCOME_FRAUD_CONFIRMED {
			_ = s.alertRepo.UpdateStatus(ctx, updatedCase.FraudAlertId, fraudv1.AlertStatus_ALERT_STATUS_CONFIRMED, updatedCase.InvestigatorId)
			alert, alertErr := s.alertRepo.GetByID(ctx, updatedCase.FraudAlertId)
			if alertErr == nil && s.publisher != nil {
				_ = s.publisher.PublishFraudConfirmed(ctx, updatedCase, alert.EntityType, alert.EntityId, correlationIDFromContext(ctx))
			}
		}
	}

	return &fraudservicev1.UpdateFraudCaseResponse{Message: "fraud case updated"}, nil
}

func (s *FraudService) ListFraudRules(ctx context.Context, req *fraudservicev1.ListFraudRulesRequest) (*fraudservicev1.ListFraudRulesResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	offset := decodePageToken(req.PageToken)

	category, ok := parseRuleCategory(req.Category)
	if !ok {
		return nil, fmt.Errorf("%w: invalid category", ErrInvalidArgument)
	}

	rules, total, err := s.ruleRepo.List(ctx, category, req.ActiveOnly, limit, offset)
	if err != nil {
		return nil, err
	}

	nextToken := ""
	if int32(offset+len(rules)) < total {
		nextToken = strconv.Itoa(offset + len(rules))
	}

	return &fraudservicev1.ListFraudRulesResponse{
		FraudRules:    rules,
		NextPageToken: nextToken,
		TotalCount:    total,
	}, nil
}

func (s *FraudService) CreateFraudRule(ctx context.Context, req *fraudservicev1.CreateFraudRuleRequest) (*fraudservicev1.CreateFraudRuleResponse, error) {
	if req.FraudRule == nil {
		return nil, fmt.Errorf("%w: fraud_rule is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.FraudRule.Name) == "" {
		return nil, fmt.Errorf("%w: fraud_rule.name is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.FraudRule.Conditions) == "" {
		return nil, fmt.Errorf("%w: fraud_rule.conditions is required", ErrInvalidArgument)
	}
	if req.FraudRule.Category == fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED {
		return nil, fmt.Errorf("%w: fraud_rule.category is required", ErrInvalidArgument)
	}
	if req.FraudRule.RiskLevel == fraudv1.RiskLevel_RISK_LEVEL_UNSPECIFIED {
		req.FraudRule.RiskLevel = fraudv1.RiskLevel_RISK_LEVEL_MEDIUM
	}
	if req.FraudRule.ScoreWeight <= 0 {
		req.FraudRule.ScoreWeight = 10
	}
	if err := s.ruleRepo.Create(ctx, req.FraudRule); err != nil {
		return nil, err
	}
	return &fraudservicev1.CreateFraudRuleResponse{RuleId: req.FraudRule.FraudRuleId, Message: "fraud rule created"}, nil
}

func (s *FraudService) UpdateFraudRule(ctx context.Context, req *fraudservicev1.UpdateFraudRuleRequest) (*fraudservicev1.UpdateFraudRuleResponse, error) {
	if strings.TrimSpace(req.RuleId) == "" || req.FraudRule == nil {
		return nil, fmt.Errorf("%w: rule_id and fraud_rule are required", ErrInvalidArgument)
	}
	if err := s.ruleRepo.Update(ctx, req.RuleId, req.FraudRule); err != nil {
		if errors.Is(err, repository.ErrRuleNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &fraudservicev1.UpdateFraudRuleResponse{Message: "fraud rule updated"}, nil
}

func (s *FraudService) ActivateFraudRule(ctx context.Context, req *fraudservicev1.ActivateFraudRuleRequest) (*fraudservicev1.ActivateFraudRuleResponse, error) {
	if strings.TrimSpace(req.RuleId) == "" {
		return nil, fmt.Errorf("%w: rule_id is required", ErrInvalidArgument)
	}
	if err := s.ruleRepo.SetActive(ctx, req.RuleId, true); err != nil {
		if errors.Is(err, repository.ErrRuleNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	s.metrics.IncRulesActivated()
	return &fraudservicev1.ActivateFraudRuleResponse{Message: "fraud rule activated"}, nil
}

func (s *FraudService) DeactivateFraudRule(ctx context.Context, req *fraudservicev1.DeactivateFraudRuleRequest) (*fraudservicev1.DeactivateFraudRuleResponse, error) {
	if strings.TrimSpace(req.RuleId) == "" {
		return nil, fmt.Errorf("%w: rule_id is required", ErrInvalidArgument)
	}
	if err := s.ruleRepo.SetActive(ctx, req.RuleId, false); err != nil {
		if errors.Is(err, repository.ErrRuleNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	s.metrics.IncRulesDeactivated()
	if strings.TrimSpace(req.Reason) != "" {
		appLogger.Infof("Fraud rule deactivated (rule_id=%s, reason=%s, trace=%s)", req.RuleId, req.Reason, correlationIDFromContext(ctx))
	}
	return &fraudservicev1.DeactivateFraudRuleResponse{Message: "fraud rule deactivated"}, nil
}

// MetricsSnapshot returns current fraud runtime counters.
func (s *FraudService) MetricsSnapshot() map[string]int64 {
	if s.metrics == nil {
		return map[string]int64{}
	}
	return s.metrics.Snapshot()
}

func decodePageToken(token string) int {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0
	}
	n, err := strconv.Atoi(token)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func parseDateRange(startDate, endDate string) (*time.Time, *time.Time) {
	parse := func(v string) *time.Time {
		v = strings.TrimSpace(v)
		if v == "" {
			return nil
		}
		for _, layout := range []string{time.RFC3339, "2006-01-02"} {
			if t, err := time.Parse(layout, v); err == nil {
				u := t.UTC()
				return &u
			}
		}
		return nil
	}
	return parse(startDate), parse(endDate)
}

func parseCasePriority(raw string) (fraudv1.CasePriority, bool) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fraudv1.CasePriority_CASE_PRIORITY_MEDIUM, true
	}
	if iv, ok := fraudv1.CasePriority_value[v]; ok {
		return fraudv1.CasePriority(iv), true
	}
	k := "CASE_PRIORITY_" + strings.ToUpper(v)
	if iv, ok := fraudv1.CasePriority_value[k]; ok {
		return fraudv1.CasePriority(iv), true
	}
	return fraudv1.CasePriority_CASE_PRIORITY_UNSPECIFIED, false
}

func parseCaseStatus(raw string) (fraudv1.CaseStatus, bool) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fraudv1.CaseStatus_CASE_STATUS_UNSPECIFIED, true
	}
	if iv, ok := fraudv1.CaseStatus_value[v]; ok {
		return fraudv1.CaseStatus(iv), true
	}
	k := "CASE_STATUS_" + strings.ToUpper(v)
	if iv, ok := fraudv1.CaseStatus_value[k]; ok {
		return fraudv1.CaseStatus(iv), true
	}
	return fraudv1.CaseStatus_CASE_STATUS_UNSPECIFIED, false
}

func parseCaseOutcome(raw string) (fraudv1.CaseOutcome, bool) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fraudv1.CaseOutcome_CASE_OUTCOME_UNSPECIFIED, true
	}
	if iv, ok := fraudv1.CaseOutcome_value[v]; ok {
		return fraudv1.CaseOutcome(iv), true
	}
	k := "CASE_OUTCOME_" + strings.ToUpper(v)
	if iv, ok := fraudv1.CaseOutcome_value[k]; ok {
		return fraudv1.CaseOutcome(iv), true
	}
	return fraudv1.CaseOutcome_CASE_OUTCOME_UNSPECIFIED, false
}

func parseRuleCategory(raw string) (fraudv1.RuleCategory, bool) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED, true
	}
	if iv, ok := fraudv1.RuleCategory_value[v]; ok {
		return fraudv1.RuleCategory(iv), true
	}
	k := "RULE_CATEGORY_" + strings.ToUpper(v)
	if iv, ok := fraudv1.RuleCategory_value[k]; ok {
		return fraudv1.RuleCategory(iv), true
	}
	return fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED, false
}

func riskLevelFromScore(score int32) string {
	switch {
	case score >= 80:
		return "RISK_LEVEL_CRITICAL"
	case score >= 55:
		return "RISK_LEVEL_HIGH"
	case score >= 30:
		return "RISK_LEVEL_MEDIUM"
	case score > 0:
		return "RISK_LEVEL_LOW"
	default:
		return "RISK_LEVEL_LOW"
	}
}

func evaluateRule(rule *fraudv1.FraudRule, payload map[string]any) bool {
	if rule == nil || strings.TrimSpace(rule.Conditions) == "" {
		return false
	}

	var cond any
	if err := json.Unmarshal([]byte(rule.Conditions), &cond); err != nil {
		return false
	}
	return evalCondition(cond, payload)
}

func evalCondition(cond any, payload map[string]any) bool {
	switch c := cond.(type) {
	case map[string]any:
		if allRaw, ok := c["all"]; ok {
			arr, ok := allRaw.([]any)
			if !ok {
				return false
			}
			for _, item := range arr {
				if !evalCondition(item, payload) {
					return false
				}
			}
			return true
		}
		if anyRaw, ok := c["any"]; ok {
			arr, ok := anyRaw.([]any)
			if !ok {
				return false
			}
			for _, item := range arr {
				if evalCondition(item, payload) {
					return true
				}
			}
			return false
		}

		field, _ := c["field"].(string)
		op, _ := c["op"].(string)
		target := c["value"]
		if strings.TrimSpace(field) == "" {
			return false
		}
		current, ok := lookupPath(payload, field)
		if !ok {
			return false
		}
		return compareValues(current, target, op)
	default:
		return false
	}
}

func lookupPath(payload map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var current any = payload
	for _, p := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		v, ok := m[p]
		if !ok {
			return nil, false
		}
		current = v
	}
	return current, true
}

func compareValues(current any, target any, op string) bool {
	op = strings.ToLower(strings.TrimSpace(op))
	if op == "" {
		op = "eq"
	}
	cNum, cOK := toFloat(current)
	tNum, tOK := toFloat(target)
	if cOK && tOK {
		switch op {
		case "gt":
			return cNum > tNum
		case "gte", "ge":
			return cNum >= tNum
		case "lt":
			return cNum < tNum
		case "lte", "le":
			return cNum <= tNum
		case "ne", "neq":
			return cNum != tNum
		default:
			return cNum == tNum
		}
	}

	cStr := strings.TrimSpace(fmt.Sprint(current))
	tStr := strings.TrimSpace(fmt.Sprint(target))
	switch op {
	case "contains":
		return strings.Contains(strings.ToLower(cStr), strings.ToLower(tStr))
	case "ne", "neq":
		return !strings.EqualFold(cStr, tStr)
	default:
		return strings.EqualFold(cStr, tStr)
	}
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		if err == nil {
			return f, true
		}
		return 0, false
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		if err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func correlationIDFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	for _, k := range []string{"x-correlation-id", "x-request-id", "x-trace-id"} {
		vals := md.Get(k)
		if len(vals) == 0 {
			continue
		}
		if v := strings.TrimSpace(vals[0]); v != "" {
			return v
		}
	}
	return ""
}

func structToJSON(in *structpb.Struct) string {
	if in == nil {
		return ""
	}
	b, err := json.Marshal(in.AsMap())
	if err != nil {
		return ""
	}
	return string(b)
}
