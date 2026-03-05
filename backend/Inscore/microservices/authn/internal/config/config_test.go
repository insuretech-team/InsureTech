package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadAndHelpers(t *testing.T) {
	t.Setenv("DB_PASSWORD", "pw")
	t.Setenv("SSLWIRELESS_API_BASE", "https://sms.example.com")
	t.Setenv("SSLWIRELESS_SID", "sid")
	t.Setenv("SSLWIRELESS_API_KEY", "api")
	t.Setenv("JWT_PRIVATE_KEY_PATH", "priv.pem")
	t.Setenv("JWT_PUBLIC_KEY_PATH", "pub.pem")
	t.Setenv("JWT_KEY_ID", "kid")
	t.Setenv("OTP_EXPIRY", "7m")
	t.Setenv("EMAIL_TLS", "false")
	t.Setenv("CSV", "a,b")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "pw", cfg.Database.Password)
	require.Equal(t, 7*time.Minute, cfg.Security.OTPExpiry)
	require.False(t, cfg.Email.TLS)

	require.Equal(t, 10, getEnvAsInt("MISSING_INT", 10))
	require.Equal(t, true, getEnvAsBool("MISSING_BOOL", true))
	require.Equal(t, time.Minute, getEnvAsDuration("MISSING_DUR", time.Minute))
	require.Equal(t, []string{"a", "b"}, getEnvAsSlice("CSV", []string{"x"}))
}

func TestValidate_FailureCases(t *testing.T) {
	c := &Config{}
	require.Error(t, c.Validate())

	c.JWT.PrivateKeyPath = "priv"
	c.JWT.PublicKeyPath = "pub"
	c.JWT.KeyID = "kid"
	c.SMS.MaskingEnabled = true
	c.SMS.APIBase = "https://sms"
	c.SMS.SID = "sid"
	c.SMS.APIKey = "api"
	require.Error(t, c.Validate()) // DB password missing

	c.Database.Password = "pw"
	require.NoError(t, c.Validate())
}

func TestLoad_AuthnPortsFromServicesYAML_NoEnvOverride(t *testing.T) {
	t.Setenv("DB_PASSWORD", "pw")
	t.Setenv("SSLWIRELESS_API_BASE", "")
	t.Setenv("SSLWIRELESS_MASKING_ENABLED", "false")
	t.Setenv("JWT_PRIVATE_KEY_PATH", "priv.pem")
	t.Setenv("JWT_PUBLIC_KEY_PATH", "pub.pem")
	t.Setenv("JWT_KEY_ID", "kid")

	// These must not override services.yaml values.
	t.Setenv("AUTHN_GRPC_PORT", "60123")
	t.Setenv("AUTHN_HTTP_PORT", "60124")
	t.Setenv("AUTHN_PORT", "60125")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, 50060, cfg.Server.GRPCPort)
	require.Equal(t, 50061, cfg.Server.HTTPPort)
}
