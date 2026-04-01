package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestRecordDecisionAndCacheRatio(t *testing.T) {
	allowBefore := testutil.ToFloat64(DecisionsTotal.WithLabelValues("system", "allow"))
	denyBefore := testutil.ToFloat64(DecisionsTotal.WithLabelValues("agent", "deny"))

	RecordDecision("system:root", true, 4.2)
	RecordDecision("agent:tenant-1", false, 7.1)
	UpdateCacheHitRatio("L1", 0.75)

	require.GreaterOrEqual(t, testutil.ToFloat64(DecisionsTotal.WithLabelValues("system", "allow")), allowBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(DecisionsTotal.WithLabelValues("agent", "deny")), denyBefore+1)
	require.Equal(t, 0.75, testutil.ToFloat64(CacheHitRatio.WithLabelValues("L1")))
}

func TestExtractPortal(t *testing.T) {
	require.Equal(t, "system", extractPortal("system:root"))
	require.Equal(t, "unknown", extractPortal("unknown"))
}

func TestScopeAndPortalMetrics(t *testing.T) {
	scopeAllowBefore := testutil.ToFloat64(APIScopeValidations.WithLabelValues("system", "allowed"))
	scopeDenyBefore := testutil.ToFloat64(APIScopeValidations.WithLabelValues("system", "denied"))
	noScopesBefore := testutil.ToFloat64(APIScopeDenialReasons.WithLabelValues("no_scopes"))
	mismatchBefore := testutil.ToFloat64(APIScopeDenialReasons.WithLabelValues("scope_mismatch"))
	portalReqBefore := testutil.ToFloat64(PortalConfigRequests.WithLabelValues("system", "success"))

	RecordAPIScopeValidation("system:root", true, 0.4)
	RecordAPIScopeValidation("system:root", false, 0.6)
	RecordAPIScopeDenial("API key has no scopes defined")
	RecordAPIScopeDenial("scope does not match")
	RecordPortalConfigRequest("system", "success")
	RecordCacheHit(true)
	RecordCacheHit(false)

	require.GreaterOrEqual(t, testutil.ToFloat64(APIScopeValidations.WithLabelValues("system", "allowed")), scopeAllowBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(APIScopeValidations.WithLabelValues("system", "denied")), scopeDenyBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(APIScopeDenialReasons.WithLabelValues("no_scopes")), noScopesBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(APIScopeDenialReasons.WithLabelValues("scope_mismatch")), mismatchBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(PortalConfigRequests.WithLabelValues("system", "success")), portalReqBefore+1)
}
