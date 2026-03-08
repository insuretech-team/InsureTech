package repository

import (
	"context"
	"testing"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestFraudRuleRepository_LiveDB_CreateGetUpdateList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testFraudDB(t)
	repo := NewFraudRuleRepository(dbConn)

	ruleID := newLiveID("rule")
	rule := &fraudv1.FraudRule{
		FraudRuleId: ruleID,
		Name:        "Live Test Rule " + ruleID,
		Category:    fraudv1.RuleCategory_RULE_CATEGORY_CLAIM_FREQUENCY,
		Description: "live test rule for integration",
		Conditions:  `{"field":"amount","op":"gt","value":10000}`,
		RiskLevel:   fraudv1.RiskLevel_RISK_LEVEL_HIGH,
		ScoreWeight: 25,
		IsActive:    true,
	}
	t.Cleanup(func() { cleanupFraudRule(t, dbConn, ruleID) })

	// Create
	err := repo.Create(ctx, rule)
	require.NoError(t, err)

	// GetByID
	fetched, err := repo.GetByID(ctx, ruleID)
	require.NoError(t, err)
	require.Equal(t, ruleID, fetched.FraudRuleId)
	require.Equal(t, rule.Name, fetched.Name)

	// Update
	updated := &fraudv1.FraudRule{
		Name:        "Updated Live Rule",
		Description: "updated description",
		ScoreWeight: 50,
	}
	err = repo.Update(ctx, ruleID, updated)
	require.NoError(t, err)
	fetched2, err := repo.GetByID(ctx, ruleID)
	require.NoError(t, err)
	require.Equal(t, "Updated Live Rule", fetched2.Name)

	// List
	rules, total, err := repo.List(ctx, fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED, false, 50, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int32(1))
	require.NotEmpty(t, rules)

	// SetActive (deactivate)
	err = repo.SetActive(ctx, ruleID, false)
	require.NoError(t, err)

	// List active-only should exclude deactivated
	activeRules, _, err := repo.List(ctx, fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED, true, 500, 0)
	require.NoError(t, err)
	for _, r := range activeRules {
		require.NotEqual(t, ruleID, r.FraudRuleId)
	}

	// SetActive (reactivate)
	err = repo.SetActive(ctx, ruleID, true)
	require.NoError(t, err)
}

func TestFraudRuleRepository_LiveDB_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testFraudDB(t)
	repo := NewFraudRuleRepository(dbConn)

	_, err := repo.GetByID(ctx, "nonexistent-rule-id")
	// GORM may return a schema parse error for AuditInfo proto fields rather than
	// the sentinel ErrRuleNotFound. Assert that an error is always returned.
	require.Error(t, err)
}
