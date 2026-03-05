package repository

import (
	"context"
	"testing"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestCasbinRuleRepo_LiveDB_UpsertListDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewCasbinRuleRepo(dbConn)

	domain := "agent:" + newLiveID("casbin")
	subject := "role:" + newLiveID("role")
	object := "svc:claim/view"
	action := "GET"

	t.Cleanup(func() { cleanupCasbinByDomainPrefix(t, dbConn, domain) })

	_, err := repo.Upsert(ctx, &authzentityv1.CasbinRule{
		Ptype: "p",
		V0:    subject,
		V1:    domain,
		V2:    object,
		V3:    action,
		V4:    "allow",
		V5:    "",
	})
	require.NoError(t, err)

	_, err = repo.Upsert(ctx, &authzentityv1.CasbinRule{
		Ptype: "p",
		V0:    subject,
		V1:    domain,
		V2:    object,
		V3:    action,
		V4:    "deny",
		V5:    "ctx.owner == true",
	})
	require.NoError(t, err)

	rules, err := repo.ListByDomain(ctx, domain)
	require.NoError(t, err)
	require.NotEmpty(t, rules)
	found := false
	for _, r := range rules {
		if r.Ptype == "p" && r.V0 == subject && r.V1 == domain && r.V2 == object && r.V3 == action {
			found = true
			require.Equal(t, "deny", r.V4)
			require.Equal(t, "ctx.owner == true", r.V5)
		}
	}
	require.True(t, found, "expected upserted casbin policy in ListByDomain result")

	require.NoError(t, repo.DeleteByDomainAndSubject(ctx, domain, subject))
	rulesAfterDeleteSubject, err := repo.ListByDomain(ctx, domain)
	require.NoError(t, err)
	for _, r := range rulesAfterDeleteSubject {
		require.False(t, r.V0 == subject && r.V1 == domain)
	}

	rule2 := &authzentityv1.CasbinRule{
		Ptype: "p",
		V0:    "role:" + newLiveID("role2"),
		V1:    domain,
		V2:    "svc:claim/edit",
		V3:    "PUT",
		V4:    "allow",
	}
	_, err = repo.Upsert(ctx, rule2)
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, rule2))

	rulesFinal, err := repo.ListByDomain(ctx, domain)
	require.NoError(t, err)
	for _, r := range rulesFinal {
		require.False(t, r.V0 == rule2.V0 && r.V2 == rule2.V2 && r.V3 == rule2.V3)
	}
}

