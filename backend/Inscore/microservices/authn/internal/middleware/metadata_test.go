package middleware

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestMetadataExtractor_ExtractAll(t *testing.T) {
	m := NewMetadataExtractor()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "10.0.0.1, 10.0.0.2",
		"user-agent", "ua-test",
		"x-device-id", "dev-1",
		"cookie", "session_token=s123; theme=light",
		"x-csrf-token", "csrf-1",
		"authorization", "Bearer tok-1",
	))

	out := m.ExtractAll(ctx)
	require.Equal(t, "10.0.0.1", out.IPAddress)
	require.Equal(t, "ua-test", out.UserAgent)
	require.Equal(t, "dev-1", out.DeviceID)
	require.Equal(t, "s123", out.SessionToken)
	require.Equal(t, "csrf-1", out.CSRFToken)
	require.Equal(t, "tok-1", out.Authorization)
}

func TestMetadataExtractor_Fallbacks(t *testing.T) {
	m := NewMetadataExtractor()
	base := context.Background()
	ctx := peer.NewContext(base, &peer.Peer{Addr: &net.IPAddr{IP: net.ParseIP("127.0.0.1")}})

	require.Equal(t, "127.0.0.1", m.ExtractIPAddress(ctx))
	require.Equal(t, "unknown", m.ExtractUserAgent(ctx))
	require.Equal(t, "", m.ExtractDeviceID(ctx))
	require.Equal(t, "", m.ExtractSessionToken(ctx))
	require.Equal(t, "", m.ExtractCSRFToken(ctx))
	require.Equal(t, "", m.ExtractAuthorizationToken(ctx))
	require.Equal(t, "", parseCookie("k1=v1; k2=v2", "missing"))
}
