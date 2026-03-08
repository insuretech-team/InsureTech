package repository

import (
	"context"
	"testing"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestFraudCaseRepository_LiveDB_CreateGetUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testFraudDB(t)
	alertRepo := NewFraudAlertRepository(dbConn)
	caseRepo := NewFraudCaseRepository(dbConn)
	ruleRepo := NewFraudRuleRepository(dbConn)
	userID := liveUserID(t, dbConn)

	ruleID := newLiveID("rule_ref")
	rule := &fraudv1.FraudRule{
		FraudRuleId: ruleID,
		Name:        "Case live rule " + ruleID,
		Category:    fraudv1.RuleCategory_RULE_CATEGORY_CLAIM_FREQUENCY,
		Description: "rule for case live test",
		Conditions:  `{"field":"claims_count","op":"gt","value":3}`,
		RiskLevel:   fraudv1.RiskLevel_RISK_LEVEL_MEDIUM,
		ScoreWeight: 15,
		IsActive:    true,
	}
	t.Cleanup(func() { cleanupFraudRule(t, dbConn, ruleID) })
	require.NoError(t, ruleRepo.Create(ctx, rule))

	// First, create an alert (cases reference alerts via fraud_alert_id).
	alertID := newLiveID("case_alert")
	alert := &fraudv1.FraudAlert{
		Id:          alertID,
		EntityType:  "CLAIM",
		EntityId:    newLiveID("claim"),
		FraudRuleId: ruleID,
		RiskLevel:   "RISK_LEVEL_MEDIUM",
		FraudScore:  40,
		Details:     "{}",
		Status:      fraudv1.AlertStatus_ALERT_STATUS_OPEN,
	}
	t.Cleanup(func() { cleanupFraudAlert(t, dbConn, alertID) })
	require.NoError(t, alertRepo.Create(ctx, alert))

	// Create case
	caseID := newLiveID("case")
	fc := &fraudv1.FraudCase{
		Id:                 caseID,
		FraudAlertId:       alertID,
		Priority:           fraudv1.CasePriority_CASE_PRIORITY_HIGH,
		InvestigationNotes: "live test case notes",
		Status:             fraudv1.CaseStatus_CASE_STATUS_OPEN,
		InvestigatorId:     userID,
	}
	t.Cleanup(func() { cleanupFraudCase(t, dbConn, caseID) })
	err := caseRepo.Create(ctx, fc)
	require.NoError(t, err)
	require.NotEmpty(t, fc.CaseNumber)

	// GetByID
	fetched, err := caseRepo.GetByID(ctx, caseID)
	require.NoError(t, err)
	require.Equal(t, caseID, fetched.Id)
	require.Equal(t, alertID, fetched.FraudAlertId)

	// Update
	err = caseRepo.Update(ctx, caseID,
		fraudv1.CaseStatus_CASE_STATUS_CLOSED,
		fraudv1.CaseOutcome_CASE_OUTCOME_FRAUD_CONFIRMED,
		"confirmed via live DB test", `{"evidence":"test"}`,
	)
	require.NoError(t, err)

	updated, err := caseRepo.GetByID(ctx, caseID)
	require.NoError(t, err)
	require.Equal(t, "CASE_STATUS_CLOSED", updated.Status.String())
	require.Equal(t, "CASE_OUTCOME_FRAUD_CONFIRMED", updated.Outcome.String())
}

func TestFraudCaseRepository_LiveDB_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testFraudDB(t)
	caseRepo := NewFraudCaseRepository(dbConn)

	_, err := caseRepo.GetByID(ctx, "nonexistent-case-id")
	require.Error(t, err)
}
