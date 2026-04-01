package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func histogramSampleCount(t *testing.T, metric prometheus.Metric) uint64 {
	t.Helper()

	pb := &dto.Metric{}
	require.NoError(t, metric.Write(pb))
	require.NotNil(t, pb.Histogram)
	return pb.GetHistogram().GetSampleCount()
}

func TestAuthNMetrics_Recorders(t *testing.T) {
	rotationBefore := testutil.ToFloat64(APIKeyRotationsTotal.WithLabelValues("B2B", "success"))
	rotationFailureBefore := testutil.ToFloat64(APIKeyRotationsTotal.WithLabelValues("B2C", "failure"))
	cacheHitBefore := testutil.ToFloat64(PortalConfigCacheHits.WithLabelValues("agent", "true"))
	cacheMissBefore := testutil.ToFloat64(PortalConfigCacheHits.WithLabelValues("agent", "false"))
	passwordBefore := testutil.ToFloat64(PasswordValidationFailures.WithLabelValues("system", "too_short"))
	sessionCreateBefore := testutil.ToFloat64(SessionsCreated.WithLabelValues("system", "JWT"))
	mfaChallengeBefore := testutil.ToFloat64(MFAChallengesIssued.WithLabelValues("system", "TOTP"))
	mfaVerifyBefore := testutil.ToFloat64(MFAVerificationResults.WithLabelValues("system", "EMAIL", "failure"))

	rotationDurationBefore := histogramSampleCount(t, APIKeyRotationDuration)
	portalLoadBefore := histogramSampleCount(t, PortalConfigLoadDuration.WithLabelValues("agent").(prometheus.Metric))
	sessionLatencyBefore := histogramSampleCount(t, SessionValidationLatency.WithLabelValues("JWT", "true").(prometheus.Metric))

	RecordAPIKeyRotation("B2B", true, 0.25)
	RecordAPIKeyRotation("B2C", false, 0.5)
	UpdateAPIKeyCount("active", "B2B", 7)
	RecordPortalConfigCache("agent", true)
	RecordPortalConfigCache("agent", false)
	RecordPortalConfigLoad("agent", 0.12)
	RecordPasswordValidationFailure("system", "too_short")
	RecordSessionCreation("system", "JWT")
	RecordSessionValidation("JWT", true, 0.02)
	UpdateActiveSessions("system", "JWT", 5)
	RecordMFAChallenge("system", "TOTP")
	RecordMFAVerification("system", "EMAIL", false)

	require.GreaterOrEqual(t, testutil.ToFloat64(APIKeyRotationsTotal.WithLabelValues("B2B", "success")), rotationBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(APIKeyRotationsTotal.WithLabelValues("B2C", "failure")), rotationFailureBefore+1)
	require.Equal(t, 7.0, testutil.ToFloat64(APIKeysActive.WithLabelValues("active", "B2B")))
	require.GreaterOrEqual(t, testutil.ToFloat64(PortalConfigCacheHits.WithLabelValues("agent", "true")), cacheHitBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(PortalConfigCacheHits.WithLabelValues("agent", "false")), cacheMissBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(PasswordValidationFailures.WithLabelValues("system", "too_short")), passwordBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(SessionsCreated.WithLabelValues("system", "JWT")), sessionCreateBefore+1)
	require.Equal(t, 5.0, testutil.ToFloat64(ActiveSessions.WithLabelValues("system", "JWT")))
	require.GreaterOrEqual(t, testutil.ToFloat64(MFAChallengesIssued.WithLabelValues("system", "TOTP")), mfaChallengeBefore+1)
	require.GreaterOrEqual(t, testutil.ToFloat64(MFAVerificationResults.WithLabelValues("system", "EMAIL", "failure")), mfaVerifyBefore+1)

	require.Equal(t, rotationDurationBefore+2, histogramSampleCount(t, APIKeyRotationDuration))
	require.Equal(t, portalLoadBefore+1, histogramSampleCount(t, PortalConfigLoadDuration.WithLabelValues("agent").(prometheus.Metric)))
	require.Equal(t, sessionLatencyBefore+1, histogramSampleCount(t, SessionValidationLatency.WithLabelValues("JWT", "true").(prometheus.Metric)))
}
