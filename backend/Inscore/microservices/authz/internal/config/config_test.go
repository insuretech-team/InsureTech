package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad_ValidateAndParsers(t *testing.T) {
	t.Setenv("DB_PASSWORD", "")
	_, err := Load()
	require.Error(t, err)

	t.Setenv("DB_PASSWORD", "pw")
	t.Setenv("AUTHZ_GRPC_PORT", "5123")
	t.Setenv("AUTHZ_POLICY_CACHE_TTL", "90s")
	t.Setenv("CASBIN_AUDIT_ALL", "true")
	t.Setenv("AUTHZ_JWT_PUBLIC_KEY_PEM", "pem")
	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, 5123, cfg.Server.GRPCPort)
	require.Equal(t, 90*time.Second, cfg.Redis.PolicyCacheTTL)
	require.True(t, cfg.Casbin.AuditAllDecisions)
	require.Equal(t, "pem", cfg.Auth.PublicKeyPEM)
	require.True(t, cfg.Casbin.DenyByDefault)
}

func TestGetEnvHelpers_Defaults(t *testing.T) {
	t.Setenv("X_INT", "bad")
	t.Setenv("X_BOOL", "bad")
	t.Setenv("X_DUR", "bad")
	require.Equal(t, "d", getEnv("X_STR", "d"))
	require.Equal(t, 7, getEnvAsInt("X_INT", 7))
	require.Equal(t, true, getEnvAsBool("X_BOOL", true))
	require.Equal(t, 3*time.Second, getEnvAsDuration("X_DUR", 3*time.Second))
}

