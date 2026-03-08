package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeMetrics_IncrementAndSnapshot(t *testing.T) {
	m := NewRuntimeMetrics()

	m.IncPartnerCreated()
	m.IncPartnerCreated()
	m.IncPartnerFetched()
	m.IncPartnerUpdated()
	m.IncPartnerListed()
	m.IncPartnerDeleted()
	m.IncPartnerVerified()
	m.IncPartnerStatusUpdate()
	m.IncCommissionCalculated()
	m.IncAPIKeyRotated()

	snap := m.Snapshot()
	require.Equal(t, int64(2), snap["partner_created"])
	require.Equal(t, int64(1), snap["partner_fetched"])
	require.Equal(t, int64(1), snap["partner_updated"])
	require.Equal(t, int64(1), snap["partner_listed"])
	require.Equal(t, int64(1), snap["partner_deleted"])
	require.Equal(t, int64(1), snap["partner_verified"])
	require.Equal(t, int64(1), snap["partner_status_updates"])
	require.Equal(t, int64(1), snap["commission_calculated"])
	require.Equal(t, int64(1), snap["api_key_rotated"])
}

func TestRuntimeMetrics_SnapshotEmpty(t *testing.T) {
	m := NewRuntimeMetrics()
	snap := m.Snapshot()
	for k, v := range snap {
		require.Equal(t, int64(0), v, "expected zero for %s", k)
	}
}
