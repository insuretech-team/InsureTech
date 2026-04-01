package grpcclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envTLSEnabled   = "INTERNAL_GRPC_TLS_ENABLED"
	envMTLSEnabled  = "INTERNAL_GRPC_MTLS_ENABLED"
	envCAFile       = "INTERNAL_GRPC_CA_FILE"
	envClientCert   = "INTERNAL_GRPC_CLIENT_CERT_FILE"
	envClientKey    = "INTERNAL_GRPC_CLIENT_KEY_FILE"
	envServerName   = "INTERNAL_GRPC_SERVER_NAME"
	envInsecureSkip = "INTERNAL_GRPC_INSECURE_SKIP_VERIFY"
)

func NewClient(target string, extraOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	baseOpts, err := DefaultDialOptions()
	if err != nil {
		return nil, err
	}
	return grpc.NewClient(target, append(baseOpts, extraOpts...)...)
}

func DefaultDialOptions() ([]grpc.DialOption, error) {
	creds, err := TransportCredentialsFromEnv()
	if err != nil {
		return nil, err
	}
	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"pick_first"}`),
	}, nil
}

func TransportCredentialsFromEnv() (credentials.TransportCredentials, error) {
	if !envBool(envTLSEnabled) {
		return FallbackInsecureCredentials(), nil
	}

	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		ServerName:         strings.TrimSpace(os.Getenv(envServerName)),
		InsecureSkipVerify: envBool(envInsecureSkip),
	}

	if caFile := strings.TrimSpace(os.Getenv(envCAFile)); caFile != "" {
		pemBytes, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read internal grpc ca file: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pemBytes) {
			return nil, fmt.Errorf("parse internal grpc ca file: %s", caFile)
		}
		tlsConfig.RootCAs = pool
	}

	if envBool(envMTLSEnabled) {
		certFile := strings.TrimSpace(os.Getenv(envClientCert))
		keyFile := strings.TrimSpace(os.Getenv(envClientKey))
		if certFile == "" || keyFile == "" {
			return nil, fmt.Errorf("%s and %s are required when %s=true", envClientCert, envClientKey, envMTLSEnabled)
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("load internal grpc client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return credentials.NewTLS(tlsConfig), nil
}

func FallbackInsecureCredentials() credentials.TransportCredentials {
	return insecure.NewCredentials()
}

func envBool(key string) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	return value == "1" || value == "true" || value == "yes"
}
