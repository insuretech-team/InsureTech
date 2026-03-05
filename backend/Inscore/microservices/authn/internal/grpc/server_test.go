package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultServerConfig(t *testing.T) {
	cfg := DefaultServerConfig()
	require.NotNil(t, cfg)
	require.Equal(t, "0.0.0.0", cfg.Host)
	require.Equal(t, "50053", cfg.Port)
}

func TestNewServer_StartInvalidPort_StopAndHealthCheck(t *testing.T) {
	s, err := NewServer(&Config{Port: "invalid"}, nil)
	require.NoError(t, err)
	require.NotNil(t, s)

	err = s.Start()
	require.Error(t, err)

	s.Stop()
	require.Error(t, s.HealthCheck(context.Background()))
}

func TestServer_HealthCheck_ContextCanceled(t *testing.T) {
	s, err := NewServer(&Config{Port: "0"}, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// DB is nil, so this still returns the explicit nil DB error.
	require.Error(t, s.HealthCheck(ctx))

	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server stop timed out")
	}
}
