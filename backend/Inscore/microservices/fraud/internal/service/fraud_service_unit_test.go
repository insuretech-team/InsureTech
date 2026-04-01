package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/repository"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

type fakeRuleRepo struct {
	listRules      []*fraudv1.FraudRule
	listTotal      int32
	listErr        error
	createErr      error
	updateErr      error
	setActiveErr   error
	gotCategory    fraudv1.RuleCategory
	gotActiveOnly  bool
	gotLimit       int
	gotOffset      int
	createdRule    *fraudv1.FraudRule
	updatedRuleID  string
	updatedRule    *fraudv1.FraudRule
	setActiveID    string
	setActiveValue bool
}

func (f *fakeRuleRepo) Create(_ context.Context, rule *fraudv1.FraudRule) error {
	f.createdRule = cloneRule(rule)
	if f.createdRule != nil && f.createdRule.FraudRuleId == "" {
		f.createdRule.FraudRuleId = uuid.NewString()
		rule.FraudRuleId = f.createdRule.FraudRuleId
	}
	return f.createErr
}

func (f *fakeRuleRepo) GetByID(_ context.Context, _ string) (*fraudv1.FraudRule, error) {
	return nil, repository.ErrRuleNotFound
}

func (f *fakeRuleRepo) Update(_ context.Context, ruleID string, rule *fraudv1.FraudRule) error {
	f.updatedRuleID = ruleID
	f.updatedRule = cloneRule(rule)
	return f.updateErr
}

func (f *fakeRuleRepo) List(_ context.Context, category fraudv1.RuleCategory, activeOnly bool, limit, offset int) ([]*fraudv1.FraudRule, int32, error) {
	f.gotCategory = category
	f.gotActiveOnly = activeOnly
	f.gotLimit = limit
	f.gotOffset = offset
	return f.listRules, f.listTotal, f.listErr
}

func (f *fakeRuleRepo) SetActive(_ context.Context, ruleID string, active bool) error {
	f.setActiveID = ruleID
	f.setActiveValue = active
	return f.setActiveErr
}

type fakeAlertRepo struct {
	getAlert          *fraudv1.FraudAlert
	getErr            error
	listAlerts        []*fraudv1.FraudAlert
	listTotal         int32
	listErr           error
	createErr         error
	updateErr         error
	createdAlert      *fraudv1.FraudAlert
	gotStatus         string
	gotRiskLevel      string
	gotStart          *time.Time
	gotEnd            *time.Time
	gotLimit          int
	gotOffset         int
	updateStatusCalls []alertStatusCall
}

type alertStatusCall struct {
	alertID    string
	status     fraudv1.AlertStatus
	assignedTo string
}

func (f *fakeAlertRepo) Create(_ context.Context, alert *fraudv1.FraudAlert) error {
	f.createdAlert = cloneAlert(alert)
	if f.createdAlert != nil && f.createdAlert.Id == "" {
		f.createdAlert.Id = uuid.NewString()
		alert.Id = f.createdAlert.Id
	}
	return f.createErr
}

func (f *fakeAlertRepo) GetByID(_ context.Context, _ string) (*fraudv1.FraudAlert, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return cloneAlert(f.getAlert), nil
}

func (f *fakeAlertRepo) List(_ context.Context, status string, riskLevel string, start, end *time.Time, limit, offset int) ([]*fraudv1.FraudAlert, int32, error) {
	f.gotStatus = status
	f.gotRiskLevel = riskLevel
	f.gotStart = start
	f.gotEnd = end
	f.gotLimit = limit
	f.gotOffset = offset
	return f.listAlerts, f.listTotal, f.listErr
}

func (f *fakeAlertRepo) UpdateStatus(_ context.Context, alertID string, status fraudv1.AlertStatus, assignedTo string) error {
	f.updateStatusCalls = append(f.updateStatusCalls, alertStatusCall{
		alertID:    alertID,
		status:     status,
		assignedTo: assignedTo,
	})
	return f.updateErr
}

type fakeCaseRepo struct {
	getCase     *fraudv1.FraudCase
	getErr      error
	createErr   error
	updateErr   error
	createdCase *fraudv1.FraudCase
	updatedCase updateCaseCall
}

type updateCaseCall struct {
	caseID   string
	status   fraudv1.CaseStatus
	outcome  fraudv1.CaseOutcome
	notes    string
	evidence string
}

func (f *fakeCaseRepo) Create(_ context.Context, fraudCase *fraudv1.FraudCase) error {
	f.createdCase = cloneCase(fraudCase)
	if f.createdCase != nil {
		if f.createdCase.Id == "" {
			f.createdCase.Id = uuid.NewString()
			fraudCase.Id = f.createdCase.Id
		}
		if f.createdCase.CaseNumber == "" {
			f.createdCase.CaseNumber = "FRC-TEST-001"
			fraudCase.CaseNumber = f.createdCase.CaseNumber
		}
	}
	return f.createErr
}

func (f *fakeCaseRepo) GetByID(_ context.Context, _ string) (*fraudv1.FraudCase, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return cloneCase(f.getCase), nil
}

func (f *fakeCaseRepo) Update(_ context.Context, caseID string, status fraudv1.CaseStatus, outcome fraudv1.CaseOutcome, notes string, evidence string) error {
	f.updatedCase = updateCaseCall{
		caseID:   caseID,
		status:   status,
		outcome:  outcome,
		notes:    notes,
		evidence: evidence,
	}
	return f.updateErr
}

func TestFraudServiceCheckFraudCreatesAlertAndMetrics(t *testing.T) {
	ruleRepo := &fakeRuleRepo{
		listRules: []*fraudv1.FraudRule{
			nil,
			{FraudRuleId: "bad-json", Conditions: "{"},
			{
				FraudRuleId: "rule-1",
				Name:        "High amount",
				Conditions:  `{"field":"claim.amount","op":"gt","value":1000}`,
				ScoreWeight: 35,
			},
			{
				FraudRuleId: "rule-2",
				Conditions:  `{"all":[{"field":"customer.segment","op":"eq","value":"vip"},{"field":"flags.count","op":"gte","value":2}]}`,
				ScoreWeight: 25,
			},
		},
		listTotal: 4,
	}
	alertRepo := &fakeAlertRepo{}
	svc := NewFraudService(ruleRepo, alertRepo, &fakeCaseRepo{}, nil)

	data, err := structpb.NewStruct(map[string]any{
		"claim":    map[string]any{"amount": 1500},
		"customer": map[string]any{"segment": "VIP"},
		"flags":    map[string]any{"count": 2},
	})
	require.NoError(t, err)

	resp, err := svc.CheckFraud(context.Background(), &fraudservicev1.CheckFraudRequest{
		EntityType: "CLAIM",
		EntityId:   "claim-1",
		Data:       data,
	})
	require.NoError(t, err)
	require.True(t, resp.IsFraudDetected)
	require.Equal(t, int32(60), resp.FraudScore)
	require.Equal(t, "RISK_LEVEL_HIGH", resp.RiskLevel)
	require.Equal(t, []string{"High amount", "rule-2"}, resp.TriggeredRules)
	require.NotEmpty(t, resp.FraudAlertId)
	require.NotNil(t, alertRepo.createdAlert)
	require.Equal(t, "CLAIM", alertRepo.createdAlert.EntityType)
	require.Equal(t, "claim-1", alertRepo.createdAlert.EntityId)
	require.Equal(t, "rule-1", alertRepo.createdAlert.FraudRuleId)
	require.Equal(t, int32(60), alertRepo.createdAlert.FraudScore)

	snap := svc.MetricsSnapshot()
	require.Equal(t, int64(1), snap["fraud_checks"])
	require.Equal(t, int64(1), snap["fraud_detections"])
	require.Equal(t, int64(1), snap["alerts_created"])
}

func TestFraudServiceCheckFraudValidationAndNoMatch(t *testing.T) {
	svc := NewFraudService(&fakeRuleRepo{}, &fakeAlertRepo{}, &fakeCaseRepo{}, nil)

	_, err := svc.CheckFraud(context.Background(), &fraudservicev1.CheckFraudRequest{})
	require.ErrorIs(t, err, ErrInvalidArgument)

	resp, err := svc.CheckFraud(context.Background(), &fraudservicev1.CheckFraudRequest{
		EntityType: "POLICY",
		EntityId:   "policy-1",
	})
	require.NoError(t, err)
	require.False(t, resp.IsFraudDetected)
	require.Empty(t, resp.FraudAlertId)
}

func TestFraudServiceGetAndListFraudAlerts(t *testing.T) {
	expectedAlert := &fraudv1.FraudAlert{Id: "alert-1"}
	alertRepo := &fakeAlertRepo{
		getAlert:   expectedAlert,
		listAlerts: []*fraudv1.FraudAlert{expectedAlert},
		listTotal:  3,
	}
	svc := NewFraudService(&fakeRuleRepo{}, alertRepo, &fakeCaseRepo{}, nil)

	getResp, err := svc.GetFraudAlert(context.Background(), &fraudservicev1.GetFraudAlertRequest{FraudAlertId: "alert-1"})
	require.NoError(t, err)
	require.Equal(t, "alert-1", getResp.FraudAlert.Id)

	listResp, err := svc.ListFraudAlerts(context.Background(), &fraudservicev1.ListFraudAlertsRequest{
		Page:      -1,
		PageSize:  999,
		Status:    "OPEN",
		RiskLevel: "RISK_LEVEL_HIGH",
		StartDate: "2026-03-01",
		EndDate:   "2026-03-13T12:00:00Z",
	})
	require.NoError(t, err)
	require.Equal(t, int32(3), listResp.TotalCount)
	require.Len(t, listResp.FraudAlerts, 1)
	require.Equal(t, 200, alertRepo.gotLimit)
	require.Equal(t, 0, alertRepo.gotOffset)
	require.NotNil(t, alertRepo.gotStart)
	require.NotNil(t, alertRepo.gotEnd)

	_, err = svc.GetFraudAlert(context.Background(), &fraudservicev1.GetFraudAlertRequest{FraudAlertId: ""})
	require.ErrorIs(t, err, ErrInvalidArgument)

	alertRepo.getErr = repository.ErrAlertNotFound
	_, err = svc.GetFraudAlert(context.Background(), &fraudservicev1.GetFraudAlertRequest{FraudAlertId: "missing"})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestFraudServiceCreateAndUpdateFraudCase(t *testing.T) {
	alertRepo := &fakeAlertRepo{
		getAlert: &fraudv1.FraudAlert{
			Id:         "alert-1",
			EntityType: "CLAIM",
			EntityId:   "claim-1",
		},
	}
	caseRepo := &fakeCaseRepo{
		getCase: &fraudv1.FraudCase{
			Id:             "case-1",
			FraudAlertId:   "alert-1",
			InvestigatorId: "investigator-1",
		},
	}
	svc := NewFraudService(&fakeRuleRepo{}, alertRepo, caseRepo, nil)

	createResp, err := svc.CreateFraudCase(context.Background(), &fraudservicev1.CreateFraudCaseRequest{
		FraudAlertId:       "alert-1",
		Priority:           "high",
		InvestigatorId:     "investigator-1",
		InvestigationNotes: "needs review",
	})
	require.NoError(t, err)
	require.Equal(t, "fraud case created", createResp.Message)
	require.Equal(t, "FRC-TEST-001", createResp.CaseNumber)
	require.Equal(t, fraudv1.CasePriority_CASE_PRIORITY_HIGH, caseRepo.createdCase.Priority)
	require.Len(t, alertRepo.updateStatusCalls, 1)
	require.Equal(t, fraudv1.AlertStatus_ALERT_STATUS_INVESTIGATING, alertRepo.updateStatusCalls[0].status)

	evidence, err := structpb.NewStruct(map[string]any{"proof": "document"})
	require.NoError(t, err)
	updateResp, err := svc.UpdateFraudCase(context.Background(), &fraudservicev1.UpdateFraudCaseRequest{
		FraudCaseId:        "case-1",
		Status:             "closed",
		Outcome:            "fraud_confirmed",
		InvestigationNotes: "confirmed",
		Evidence:           evidence,
	})
	require.NoError(t, err)
	require.Equal(t, "fraud case updated", updateResp.Message)
	require.Equal(t, fraudv1.CaseStatus_CASE_STATUS_CLOSED, caseRepo.updatedCase.status)
	require.Equal(t, fraudv1.CaseOutcome_CASE_OUTCOME_FRAUD_CONFIRMED, caseRepo.updatedCase.outcome)
	require.JSONEq(t, `{"proof":"document"}`, caseRepo.updatedCase.evidence)
	require.Len(t, alertRepo.updateStatusCalls, 3)
	require.Equal(t, fraudv1.AlertStatus_ALERT_STATUS_CLOSED, alertRepo.updateStatusCalls[1].status)
	require.Equal(t, fraudv1.AlertStatus_ALERT_STATUS_CONFIRMED, alertRepo.updateStatusCalls[2].status)

	_, err = svc.CreateFraudCase(context.Background(), &fraudservicev1.CreateFraudCaseRequest{
		FraudAlertId: "alert-1",
		Priority:     "bad-priority",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)

	alertRepo.getErr = repository.ErrAlertNotFound
	_, err = svc.CreateFraudCase(context.Background(), &fraudservicev1.CreateFraudCaseRequest{
		FraudAlertId: "missing",
	})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestFraudServiceGetAndUpdateFraudCaseErrors(t *testing.T) {
	caseRepo := &fakeCaseRepo{getErr: repository.ErrCaseNotFound}
	svc := NewFraudService(&fakeRuleRepo{}, &fakeAlertRepo{}, caseRepo, nil)

	_, err := svc.GetFraudCase(context.Background(), &fraudservicev1.GetFraudCaseRequest{FraudCaseId: "missing"})
	require.ErrorIs(t, err, ErrNotFound)

	_, err = svc.UpdateFraudCase(context.Background(), &fraudservicev1.UpdateFraudCaseRequest{
		FraudCaseId: "case-1",
		Status:      "invalid",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)

	_, err = svc.UpdateFraudCase(context.Background(), &fraudservicev1.UpdateFraudCaseRequest{
		FraudCaseId: "case-1",
		Status:      "closed",
		Outcome:     "invalid",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)
}

func TestFraudServiceRuleLifecycle(t *testing.T) {
	ruleRepo := &fakeRuleRepo{
		listRules: []*fraudv1.FraudRule{{FraudRuleId: "rule-1"}},
		listTotal: 3,
	}
	svc := NewFraudService(ruleRepo, &fakeAlertRepo{}, &fakeCaseRepo{}, nil)

	createResp, err := svc.CreateFraudRule(context.Background(), &fraudservicev1.CreateFraudRuleRequest{
		FraudRule: &fraudv1.FraudRule{
			Name:       "Rule",
			Category:   fraudv1.RuleCategory_RULE_CATEGORY_AMOUNT_ANOMALY,
			Conditions: `{"field":"amount","value":10}`,
		},
	})
	require.NoError(t, err)
	require.Equal(t, "fraud rule created", createResp.Message)
	require.NotNil(t, ruleRepo.createdRule)
	require.Equal(t, fraudv1.RiskLevel_RISK_LEVEL_MEDIUM, ruleRepo.createdRule.RiskLevel)
	require.Equal(t, int32(10), ruleRepo.createdRule.ScoreWeight)

	listResp, err := svc.ListFraudRules(context.Background(), &fraudservicev1.ListFraudRulesRequest{
		Category:   "claim_frequency",
		ActiveOnly: true,
		PageSize:   2,
		PageToken:  "1",
	})
	require.NoError(t, err)
	require.Equal(t, int32(3), listResp.TotalCount)
	require.Equal(t, "2", listResp.NextPageToken)
	require.Equal(t, fraudv1.RuleCategory_RULE_CATEGORY_CLAIM_FREQUENCY, ruleRepo.gotCategory)
	require.True(t, ruleRepo.gotActiveOnly)
	require.Equal(t, 2, ruleRepo.gotLimit)
	require.Equal(t, 1, ruleRepo.gotOffset)

	updateResp, err := svc.UpdateFraudRule(context.Background(), &fraudservicev1.UpdateFraudRuleRequest{
		RuleId:    "rule-1",
		FraudRule: &fraudv1.FraudRule{Name: "Updated"},
	})
	require.NoError(t, err)
	require.Equal(t, "fraud rule updated", updateResp.Message)
	require.Equal(t, "rule-1", ruleRepo.updatedRuleID)

	activateResp, err := svc.ActivateFraudRule(context.Background(), &fraudservicev1.ActivateFraudRuleRequest{RuleId: "rule-1"})
	require.NoError(t, err)
	require.Equal(t, "fraud rule activated", activateResp.Message)
	require.True(t, ruleRepo.setActiveValue)

	deactivateResp, err := svc.DeactivateFraudRule(context.Background(), &fraudservicev1.DeactivateFraudRuleRequest{
		RuleId: "rule-1",
		Reason: "expired",
	})
	require.NoError(t, err)
	require.Equal(t, "fraud rule deactivated", deactivateResp.Message)
	require.False(t, ruleRepo.setActiveValue)

	_, err = svc.ListFraudRules(context.Background(), &fraudservicev1.ListFraudRulesRequest{Category: "bad"})
	require.ErrorIs(t, err, ErrInvalidArgument)
}

func TestFraudServiceHelpers(t *testing.T) {
	require.Equal(t, 0, decodePageToken(""))
	require.Equal(t, 0, decodePageToken("-4"))
	require.Equal(t, 7, decodePageToken("7"))

	start, end := parseDateRange("2026-03-01", "2026-03-13T08:00:00Z")
	require.NotNil(t, start)
	require.NotNil(t, end)
	require.Nil(t, func() *time.Time {
		s, _ := parseDateRange("bad", "")
		return s
	}())

	priority, ok := parseCasePriority("high")
	require.True(t, ok)
	require.Equal(t, fraudv1.CasePriority_CASE_PRIORITY_HIGH, priority)

	status, ok := parseCaseStatus("CASE_STATUS_CLOSED")
	require.True(t, ok)
	require.Equal(t, fraudv1.CaseStatus_CASE_STATUS_CLOSED, status)

	outcome, ok := parseCaseOutcome("false_positive")
	require.True(t, ok)
	require.Equal(t, fraudv1.CaseOutcome_CASE_OUTCOME_FALSE_POSITIVE, outcome)

	category, ok := parseRuleCategory("amount_anomaly")
	require.True(t, ok)
	require.Equal(t, fraudv1.RuleCategory_RULE_CATEGORY_AMOUNT_ANOMALY, category)

	require.Equal(t, "RISK_LEVEL_LOW", riskLevelFromScore(0))
	require.Equal(t, "RISK_LEVEL_MEDIUM", riskLevelFromScore(30))
	require.Equal(t, "RISK_LEVEL_HIGH", riskLevelFromScore(55))
	require.Equal(t, "RISK_LEVEL_CRITICAL", riskLevelFromScore(80))

	require.True(t, evaluateRule(&fraudv1.FraudRule{Conditions: `{"any":[{"field":"amount","op":"gt","value":50},{"field":"status","value":"ok"}]}`}, map[string]any{
		"amount": 60,
	}))
	require.False(t, evaluateRule(&fraudv1.FraudRule{Conditions: `{"field":"missing","value":"x"}`}, map[string]any{}))

	current, found := lookupPath(map[string]any{"outer": map[string]any{"inner": "value"}}, "outer.inner")
	require.True(t, found)
	require.Equal(t, "value", current)
	require.False(t, compareValues("Hello World", "bye", "contains"))
	require.True(t, compareValues("Hello World", "world", "contains"))
	require.True(t, compareValues(10, "9", "gt"))

	f, ok := toFloat(json.Number("12.5"))
	require.True(t, ok)
	require.Equal(t, 12.5, f)
	_, ok = toFloat(struct{}{})
	require.False(t, ok)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-correlation-id", "corr-1"))
	require.Equal(t, "corr-1", correlationIDFromContext(ctx))
	require.Equal(t, "", correlationIDFromContext(context.Background()))

	st, err := structpb.NewStruct(map[string]any{"ok": true})
	require.NoError(t, err)
	require.JSONEq(t, `{"ok":true}`, structToJSON(st))
	require.Equal(t, "", structToJSON(nil))
}

func cloneRule(in *fraudv1.FraudRule) *fraudv1.FraudRule {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

func cloneAlert(in *fraudv1.FraudAlert) *fraudv1.FraudAlert {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

func cloneCase(in *fraudv1.FraudCase) *fraudv1.FraudCase {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}
