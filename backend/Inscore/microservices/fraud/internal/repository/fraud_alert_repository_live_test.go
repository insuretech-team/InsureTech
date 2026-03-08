package repository

import (
	"context"
	"testing"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestFraudAlertRepository_LiveDB_CreateGetListUpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testFraudDB(t)
	repo := NewFraudAlertRepository(dbConn)
	ruleRepo := NewFraudRuleRepository(dbConn)
	userID := liveUserID(t, dbConn)

	ruleID := newLiveID("rule_ref")
	rule := &fraudv1.FraudRule{
		FraudRuleId: ruleID,
		Name:        "Alert live rule " + ruleID,
		Category:    fraudv1.RuleCategory_RULE_CATEGORY_CLAIM_FREQUENCY,
		Description: "rule for alert live test",
		Conditions:  `{"field":"amount","op":"gt","value":10000}`,
		RiskLevel:   fraudv1.RiskLevel_RISK_LEVEL_HIGH,
		ScoreWeight: 20,
		IsActive:    true,
	}
	t.Cleanup(func() { cleanupFraudRule(t, dbConn, ruleID) })
	require.NoError(t, ruleRepo.Create(ctx, rule))

	alertID := newLiveID("alert")
	alert := &fraudv1.FraudAlert{
		Id:          alertID,
		EntityType:  "POLICY",
		EntityId:    newLiveID("policy"),
		FraudRuleId: ruleID,
		RiskLevel:   "RISK_LEVEL_HIGH",
		FraudScore:  65,
		Details:     `{"triggered_rules":["test_rule"]}`,
		Status:      fraudv1.AlertStatus_ALERT_STATUS_OPEN,
	}
	t.Cleanup(func() { cleanupFraudAlert(t, dbConn, alertID) })

	// Create
	err := repo.Create(ctx, alert)
	require.NoError(t, err)
	require.NotEmpty(t, alert.AlertNumber)

	// GetByID
	fetched, err := repo.GetByID(ctx, alertID)
	require.NoError(t, err)
	require.Equal(t, alertID, fetched.Id)
	require.Equal(t, "POLICY", fetched.EntityType)

	// List
	alerts, total, err := repo.List(ctx, "", "", nil, nil, 50, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int32(1))
	require.NotEmpty(t, alerts)

	// List with status filter
	openAlerts, openTotal, err := repo.List(ctx, "ALERT_STATUS_OPEN", "", nil, nil, 50, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, openTotal, int32(1))
	_ = openAlerts

	// UpdateStatus
	err = repo.UpdateStatus(ctx, alertID, fraudv1.AlertStatus_ALERT_STATUS_CLOSED, userID)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, alertID)
	require.NoError(t, err)
	require.Equal(t, "ALERT_STATUS_CLOSED", updated.Status.String())
}

func TestFraudAlertRepository_LiveDB_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testFraudDB(t)
	repo := NewFraudAlertRepository(dbConn)

	_, err := repo.GetByID(ctx, "nonexistent-alert-id")
	require.Error(t, err)
}
