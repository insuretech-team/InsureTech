package middleware

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type fakeStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (f *fakeStream) Context() context.Context { return f.ctx }

func makeJWT(t *testing.T, priv *rsa.PrivateKey, method jwt.SigningMethod, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	s, err := token.SignedString(priv)
	require.NoError(t, err)
	return s
}

func TestJWTInterceptor_UnaryAndStream(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	i := NewJWTInterceptor(&priv.PublicKey, []string{"/skip"})

	claims := jwt.MapClaims{
		"sub":    "u1",
		"portal": "system",
		"email":  "u1@example.com",
		"roles":  []string{"admin", "auditor"},
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	raw := makeJWT(t, priv, jwt.SigningMethodRS256, claims)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+raw))

	resp, err := i.UnaryInterceptor()(ctx, "req", &grpc.UnaryServerInfo{FullMethod: "/m"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		c := GetClaims(ctx)
		require.NotNil(t, c)
		require.Equal(t, "u1", c.UserID)
		require.Equal(t, "system", c.PortalID)
		require.Equal(t, "u1@example.com", c.Email)
		require.Equal(t, []string{"admin", "auditor"}, c.Roles)
		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)

	err = i.StreamInterceptor()("srv", &fakeStream{ctx: ctx}, &grpc.StreamServerInfo{FullMethod: "/m"}, func(srv interface{}, ss grpc.ServerStream) error {
		require.NotNil(t, GetClaims(ss.Context()))
		return nil
	})
	require.NoError(t, err)

	_, err = i.UnaryInterceptor()(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/m"}, func(ctx context.Context, req interface{}) (interface{}, error) { return "ok", nil })
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestJWTInterceptor_SkipNoopAndParsePEM(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	i := NewJWTInterceptor(nil, []string{"/skip"})

	_, err = i.UnaryInterceptor()(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/skip"}, func(ctx context.Context, req interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, err)

	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	pemStr := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	parsed, err := ParseRSAPublicKeyFromPEM(pemStr)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	escaped := strings.ReplaceAll(strings.TrimSpace(pemStr), "\n", `\n`)
	quotedEscaped := `"` + escaped + `"`
	parsed, err = ParseRSAPublicKeyFromPEM(quotedEscaped)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	f, err := os.CreateTemp(t.TempDir(), "authz-pub-*.pem")
	require.NoError(t, err)
	_, err = f.WriteString(pemStr)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	parsed, err = ParseRSAPublicKeyFromPEM(f.Name())
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed, err = ParseRSAPublicKeyFromPEM("")
	require.NoError(t, err)
	require.Nil(t, parsed)
	_, err = ParseRSAPublicKeyFromPEM("bad")
	require.Error(t, err)
}

func TestJWTInterceptor_InternalRoleManagementMethodsCanSkipJWT(t *testing.T) {
	i := NewJWTInterceptor(newTestPublicKey(t), []string{
		"/insuretech.authz.services.v1.AuthZService/ListRoles",
		"/insuretech.authz.services.v1.AuthZService/AssignRole",
		"/insuretech.authz.services.v1.AuthZService/CreatePolicyRule",
	})

	skipped := []string{
		"/insuretech.authz.services.v1.AuthZService/ListRoles",
		"/insuretech.authz.services.v1.AuthZService/AssignRole",
		"/insuretech.authz.services.v1.AuthZService/CreatePolicyRule",
	}
	for _, method := range skipped {
		_, err := i.UnaryInterceptor()(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: method}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		})
		require.NoError(t, err, method)
	}
}

func TestJWTInterceptor_TrustedInternalServiceCanSkipJWT(t *testing.T) {
	i := NewJWTInterceptor(newTestPublicKey(t), nil)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-internal-service", "gateway"))
	_, err := i.UnaryInterceptor()(ctx, "req", &grpc.UnaryServerInfo{FullMethod: "/insuretech.authz.services.v1.AuthZService/AssignRole"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})
	require.NoError(t, err)

	unknownCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-internal-service", "unknown-service"))
	_, err = i.UnaryInterceptor()(unknownCtx, "req", &grpc.UnaryServerInfo{FullMethod: "/insuretech.authz.services.v1.AuthZService/AssignRole"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestJWTInterceptor_InvalidSigningMethodAndPeerIP(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	i := NewJWTInterceptor(&priv.PublicKey, nil)

	badToken := makeJWT(t, priv, jwt.SigningMethodRS256, jwt.MapClaims{"sub": "u1", "exp": time.Now().Add(time.Hour).Unix()})
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+badToken))
	_, err = i.extractClaims(ctx)
	require.NoError(t, err)

	require.Equal(t, "unknown", peerIP(context.Background()))
	ctxWithPeer := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}})
	require.Equal(t, "127.0.0.1", peerIP(ctxWithPeer))
}

func newTestPublicKey(t *testing.T) *rsa.PublicKey {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return &priv.PublicKey
}
