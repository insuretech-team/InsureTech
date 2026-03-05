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

