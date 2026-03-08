package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeMetrics_IncrementAndSnapshot(t *testing.T) {
	m := NewRuntimeMetrics()

	m.IncFraudChecks()
	m.IncFraudChecks()
	m.IncFraudDetections()
	m.IncAlertsCreated()
	m.IncCasesCreated()
	m.IncRulesActivated()
	m.IncRulesDeactivated()

	snap := m.Snapshot()
	require.Equal(t, int64(2), snap["fraud_checks"])
	require.Equal(t, int64(1), snap["fraud_detections"])
	require.Equal(t, int64(1), snap["alerts_created"])
	require.Equal(t, int64(1), snap["cases_created"])
	require.Equal(t, int64(1), snap["rules_activated"])
	require.Equal(t, int64(1), snap["rules_deactivated"])
}

func TestRuntimeMetrics_SnapshotEmpty(t *testing.T) {
	m := NewRuntimeMetrics()
	snap := m.Snapshot()
	for k, v := range snap {
		require.Equal(t, int64(0), v, "expected zero for %s", k)
	}
}
