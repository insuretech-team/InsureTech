package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// TLSConfig holds TLS configuration
type TLSConfig struct {
	// Server certificate and key
	CertFile string
	KeyFile  string

	// Client CA certificate (for mTLS)
	ClientCAFile string

	// Minimum TLS version
	MinVersion uint16

	// Cipher suites (empty = use defaults)
	CipherSuites []uint16

	// Enable mutual TLS
	RequireClientCert bool
}

// DefaultTLSConfig returns secure default TLS configuration
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		MinVersion: tls.VersionTLS12, // TLS 1.2 minimum
		CipherSuites: []uint16{
			// TLS 1.3 cipher suites (preferred)
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,

			// TLS 1.2 cipher suites (fallback)
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		RequireClientCert: false,
	}
}

// LoadTLSConfig creates a tls.Config from TLSConfig
func LoadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	if cfg == nil {
		logger.Error("TLS config is nil")
		return nil, fmt.Errorf("TLS config is nil")
	}

	// Load server certificate
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		logger.Error("Failed to load server certificate", zap.Error(err))
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   cfg.MinVersion,
		CipherSuites: cfg.CipherSuites,
		// Prefer server cipher suites
		PreferServerCipherSuites: true,
	}

	// Load client CA for mTLS
	if cfg.ClientCAFile != "" {
		clientCAData, err := os.ReadFile(cfg.ClientCAFile)
		if err != nil {
			logger.Error("Failed to read client CA file", zap.Error(err))
			return nil, fmt.Errorf("failed to read client CA file: %w", err)
		}

		clientCAs := x509.NewCertPool()
		if !clientCAs.AppendCertsFromPEM(clientCAData) {
			logger.Error("Failed to parse client CA certificate")
			return nil, fmt.Errorf("failed to parse client CA certificate")
		}

		tlsConfig.ClientCAs = clientCAs

		if cfg.RequireClientCert {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		}
	}

	return tlsConfig, nil
}

// GenerateSelfSignedCert generates a self-signed certificate for development
// WARNING: Only use for development/testing - never in production!
func GenerateSelfSignedCert(certFile, keyFile string) error {
	// This is a placeholder - in real implementation, use crypto/x509 to generate
	// For now, we'll document how to generate using openssl
	msg := fmt.Sprintf("self-signed certificate generation not implemented - use openssl:\n"+
		"  openssl req -x509 -newkey rsa:4096 -keyout %s -out %s -days 365 -nodes -subj '/CN=localhost'",
		keyFile, certFile)
	logger.Error(msg)
	return fmt.Errorf("%s", msg)
}

// ValidateTLSConfig validates TLS configuration
func ValidateTLSConfig(cfg *TLSConfig) error {
	if cfg == nil {
		logger.Error("TLS config is nil")
		return fmt.Errorf("TLS config is nil")
	}

	// Check certificate files exist
	if _, err := os.Stat(cfg.CertFile); os.IsNotExist(err) {
		logger.Error("Certificate file does not exist", zap.String("file", cfg.CertFile))
		return fmt.Errorf("certificate file does not exist: %s", cfg.CertFile)
	}

	if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
		logger.Error("Key file does not exist", zap.String("file", cfg.KeyFile))
		return fmt.Errorf("key file does not exist: %s", cfg.KeyFile)
	}

	// Check client CA if mTLS is enabled
	if cfg.ClientCAFile != "" {
		if _, err := os.Stat(cfg.ClientCAFile); os.IsNotExist(err) {
			logger.Error("Client CA file does not exist", zap.String("file", cfg.ClientCAFile))
			return fmt.Errorf("client CA file does not exist: %s", cfg.ClientCAFile)
		}
	}

	// Validate TLS version
	if cfg.MinVersion < tls.VersionTLS12 {
		logger.Error("Minimum TLS version must be 1.2 or higher", zap.Uint16("version", cfg.MinVersion))
		return fmt.Errorf("minimum TLS version must be 1.2 or higher (got %d)", cfg.MinVersion)
	}

	return nil
}

// TLSConfigFromEnv creates TLS config from environment variables
func TLSConfigFromEnv() (*TLSConfig, error) {
	tlsEnabled := os.Getenv("TLS_ENABLED")
	if tlsEnabled != "true" {
		return nil, nil // TLS not enabled
	}

	cfg := DefaultTLSConfig()

	// Required fields
	cfg.CertFile = os.Getenv("TLS_CERT_FILE")
	cfg.KeyFile = os.Getenv("TLS_KEY_FILE")

	if cfg.CertFile == "" || cfg.KeyFile == "" {
		logger.Error("TLS_CERT_FILE and TLS_KEY_FILE must be set when TLS_ENABLED=true")
		return nil, fmt.Errorf("TLS_CERT_FILE and TLS_KEY_FILE must be set when TLS_ENABLED=true")
	}

	// Optional fields
	cfg.ClientCAFile = os.Getenv("TLS_CLIENT_CA_FILE")

	if os.Getenv("TLS_REQUIRE_CLIENT_CERT") == "true" {
		cfg.RequireClientCert = true
	}

	// Validate configuration
	if err := ValidateTLSConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
