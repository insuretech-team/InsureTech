package service

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
)

var (
	fraudSvcDBOnce sync.Once
	fraudSvcDB     *gorm.DB
	fraudSvcDBErr  error
)

func testFraudServiceDB(t *testing.T) *gorm.DB {
	t.Helper()

	fraudSvcDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		fraudSvcDBErr = db.InitializeManagerForService(configPath)
		if fraudSvcDBErr != nil {
			return
		}
		fraudSvcDB = db.GetDB()
	})

	if fraudSvcDBErr != nil {
		t.Skipf("skipping live DB test: %v", fraudSvcDBErr)
	}
	if fraudSvcDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return fraudSvcDB
}

func newLiveFraudService(t *testing.T) (*FraudService, *gorm.DB) {
	t.Helper()
	dbConn := testFraudServiceDB(t)
	ruleRepo := repository.NewFraudRuleRepository(dbConn)
	alertRepo := repository.NewFraudAlertRepository(dbConn)
	caseRepo := repository.NewFraudCaseRepository(dbConn)
	publisher := events.NewPublisher(nil, "fraud.events") // nil producer: events silently dropped
	svc := NewFraudService(ruleRepo, alertRepo, caseRepo, publisher)
	return svc, dbConn
}

func cleanupFraudLiveRule(t *testing.T, db *gorm.DB, ruleID string) {
	t.Helper()
	_ = db.Exec(`DELETE FROM insurance_schema.fraud_rules WHERE fraud_rule_id = ?`, ruleID).Error
}

func cleanupFraudLiveAlert(t *testing.T, db *gorm.DB, alertID string) {
	t.Helper()
	_ = db.Exec(`DELETE FROM insurance_schema.fraud_cases WHERE fraud_alert_id = ?`, alertID).Error
	_ = db.Exec(`DELETE FROM insurance_schema.fraud_alerts WHERE alert_id = ?`, alertID).Error
}

// TestFraudService_Live_RuleCRUDLifecycle tests CreateFraudRule → ListFraudRules →
// UpdateFraudRule → Activate/Deactivate via the service layer with real DB.
func TestFraudService_Live_RuleCRUDLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveFraudService(t)

	ruleID := uuid.New().String()
	rule := &fraudv1.FraudRule{
		FraudRuleId: ruleID,
		Name:        "Service Live Rule " + ruleID,
		Category:    fraudv1.RuleCategory_RULE_CATEGORY_AMOUNT_ANOMALY,
		Conditions:  `{"field":"premium","op":"gt","value":50000}`,
		RiskLevel:   fraudv1.RiskLevel_RISK_LEVEL_MEDIUM,
		ScoreWeight: 20,
	}
	t.Cleanup(func() { cleanupFraudLiveRule(t, dbConn, ruleID) })

	// Create
	createResp, err := svc.CreateFraudRule(ctx, &fraudservicev1.CreateFraudRuleRequest{FraudRule: rule})
	require.NoError(t, err)
	require.Equal(t, ruleID, createResp.RuleId)

	// List
	listResp, err := svc.ListFraudRules(ctx, &fraudservicev1.ListFraudRulesRequest{
		PageSize:   50,
		ActiveOnly: true,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, listResp.TotalCount, int32(1))

	// Update
	updateResp, err := svc.UpdateFraudRule(ctx, &fraudservicev1.UpdateFraudRuleRequest{
		RuleId: ruleID,
		FraudRule: &fraudv1.FraudRule{
			Name:        "Updated Service Rule",
			ScoreWeight: 35,
		},
	})
	require.NoError(t, err)
	require.Equal(t, "fraud rule updated", updateResp.Message)

	// Deactivate
	deactResp, err := svc.DeactivateFraudRule(ctx, &fraudservicev1.DeactivateFraudRuleRequest{
		RuleId: ruleID,
		Reason: "live test deactivation",
	})
	require.NoError(t, err)
	require.Equal(t, "fraud rule deactivated", deactResp.Message)

	// Activate
	actResp, err := svc.ActivateFraudRule(ctx, &fraudservicev1.ActivateFraudRuleRequest{
		RuleId: ruleID,
	})
	require.NoError(t, err)
	require.Equal(t, "fraud rule activated", actResp.Message)
}

// TestFraudService_Live_CheckFraud_WithNoRules ensures CheckFraud returns safely when no rules match.
func TestFraudService_Live_CheckFraud_WithNoRules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, _ := newLiveFraudService(t)

	data, _ := structpb.NewStruct(map[string]any{
		"amount":    100,
		"entity_id": "test-entity",
	})
	resp, err := svc.CheckFraud(ctx, &fraudservicev1.CheckFraudRequest{
		EntityType: "TEST",
		EntityId:   uuid.New().String(),
		Data:       data,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.GreaterOrEqual(t, resp.FraudScore, int32(0))
}

// TestFraudService_Live_CheckFraud_WithMatchingRule creates a rule, then runs CheckFraud
// to trigger it, verifying alert creation.
func TestFraudService_Live_CheckFraud_WithMatchingRule(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveFraudService(t)

	ruleID := uuid.New().String()
	rule := &fraudv1.FraudRule{
		FraudRuleId: ruleID,
		Name:        "Trigger Rule " + ruleID,
		Category:    fraudv1.RuleCategory_RULE_CATEGORY_CLAIM_FREQUENCY,
		Conditions:  `{"field":"claim_amount","op":"gt","value":100000}`,
		RiskLevel:   fraudv1.RiskLevel_RISK_LEVEL_HIGH,
		ScoreWeight: 60,
	}
	t.Cleanup(func() { cleanupFraudLiveRule(t, dbConn, ruleID) })

	_, err := svc.CreateFraudRule(ctx, &fraudservicev1.CreateFraudRuleRequest{FraudRule: rule})
	require.NoError(t, err)

	// Now perform a check with data that should trigger the rule
	data, _ := structpb.NewStruct(map[string]any{
		"claim_amount": 250000,
	})
	checkResp, err := svc.CheckFraud(ctx, &fraudservicev1.CheckFraudRequest{
		EntityType: "CLAIM",
		EntityId:   uuid.New().String(),
		Data:       data,
	})
	require.NoError(t, err)
	require.True(t, checkResp.IsFraudDetected)
	require.GreaterOrEqual(t, checkResp.FraudScore, int32(60))
	require.NotEmpty(t, checkResp.FraudAlertId)
	require.NotEmpty(t, checkResp.TriggeredRules)

	// Clean up alert created by CheckFraud
	t.Cleanup(func() { cleanupFraudLiveAlert(t, dbConn, checkResp.FraudAlertId) })

	// Verify alert was persisted
	alertResp, err := svc.GetFraudAlert(ctx, &fraudservicev1.GetFraudAlertRequest{
		FraudAlertId: checkResp.FraudAlertId,
	})
	require.NoError(t, err)
	require.Equal(t, checkResp.FraudAlertId, alertResp.FraudAlert.Id)
}

// TestFraudService_Live_CaseLifecycle creates an alert, creates a case from it,
// updates the case, and verifies the lifecycle transitions.
func TestFraudService_Live_CaseLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveFraudService(t)

	// Setup: create a rule and trigger an alert
	ruleID := uuid.New().String()
	t.Cleanup(func() { cleanupFraudLiveRule(t, dbConn, ruleID) })

	_, err := svc.CreateFraudRule(ctx, &fraudservicev1.CreateFraudRuleRequest{
		FraudRule: &fraudv1.FraudRule{
			FraudRuleId: ruleID,
			Name:        "Case Lifecycle Rule " + ruleID,
			Category:    fraudv1.RuleCategory_RULE_CATEGORY_CLAIM_FREQUENCY,
			Conditions:  `{"field": "suspicious", "op": "eq", "value": "true"}`,
			RiskLevel:   fraudv1.RiskLevel_RISK_LEVEL_CRITICAL,
			ScoreWeight: 90,
		},
	})
	require.NoError(t, err)

	data, _ := structpb.NewStruct(map[string]any{"suspicious": "true"})
	checkResp, err := svc.CheckFraud(ctx, &fraudservicev1.CheckFraudRequest{
		EntityType: "CLAIM",
		EntityId:   uuid.New().String(),
		Data:       data,
	})
	require.NoError(t, err)
	require.True(t, checkResp.IsFraudDetected)
	alertID := checkResp.FraudAlertId
	t.Cleanup(func() { cleanupFraudLiveAlert(t, dbConn, alertID) })

	// Create case from alert
	caseResp, err := svc.CreateFraudCase(ctx, &fraudservicev1.CreateFraudCaseRequest{
		FraudAlertId:       alertID,
		Priority:           "HIGH",
		InvestigationNotes: "live test case investigation",
	})
	require.NoError(t, err)
	require.NotEmpty(t, caseResp.FraudCaseId)

	// Get case
	getResp, err := svc.GetFraudCase(ctx, &fraudservicev1.GetFraudCaseRequest{
		FraudCaseId: caseResp.FraudCaseId,
	})
	require.NoError(t, err)
	require.Equal(t, alertID, getResp.FraudCase.FraudAlertId)

	// Update case → close with confirmed fraud
	updateResp, err := svc.UpdateFraudCase(ctx, &fraudservicev1.UpdateFraudCaseRequest{
		FraudCaseId:        caseResp.FraudCaseId,
		Status:             "CASE_STATUS_CLOSED",
		Outcome:            "CASE_OUTCOME_FRAUD_CONFIRMED",
		InvestigationNotes: "confirmed via live test",
	})
	require.NoError(t, err)
	require.Equal(t, "fraud case updated", updateResp.Message)

	// Verify metrics snapshot
	snap := svc.MetricsSnapshot()
	require.GreaterOrEqual(t, snap["fraud_checks"], int64(1))
	require.GreaterOrEqual(t, snap["fraud_detections"], int64(1))
	require.GreaterOrEqual(t, snap["cases_created"], int64(1))
}

// TestFraudService_Live_ValidationErrors verifies service-level validation for invalid input.
func TestFraudService_Live_ValidationErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, _ := newLiveFraudService(t)

	// CheckFraud with empty entity fields
	_, err := svc.CheckFraud(ctx, &fraudservicev1.CheckFraudRequest{
		EntityType: "",
		EntityId:   "",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// CreateFraudRule with nil rule
	_, err = svc.CreateFraudRule(ctx, &fraudservicev1.CreateFraudRuleRequest{FraudRule: nil})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// GetFraudAlert with empty ID
	_, err = svc.GetFraudAlert(ctx, &fraudservicev1.GetFraudAlertRequest{FraudAlertId: ""})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// CreateFraudCase with empty alert ID
	_, err = svc.CreateFraudCase(ctx, &fraudservicev1.CreateFraudCaseRequest{FraudAlertId: ""})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// UpdateFraudCase with empty case ID
	_, err = svc.UpdateFraudCase(ctx, &fraudservicev1.UpdateFraudCaseRequest{FraudCaseId: ""})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)
}
