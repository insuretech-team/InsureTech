package main

import (
	"context"
	"errors"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestResolveServiceAddr(t *testing.T) {
	services := map[string]serviceaddr.Service{}

	fraud := serviceaddr.Service{Name: "fraud"}
	fraud.Ports.Grpc = 50221
	services["fraud"] = fraud

	require.Equal(t, "explicit:7000", resolveServiceAddr(" explicit:7000 ", services, "fraud"))
	t.Setenv("ENVIRONMENT", "development")
	require.Equal(t, "localhost:50221", resolveServiceAddr("", services, "fraud"))
	t.Setenv("ENVIRONMENT", "production")
	require.Equal(t, "fraud:50221", resolveServiceAddr("", services, "fraud"))
	require.Equal(t, "", resolveServiceAddr("", services, "missing"))
}

func TestLoggingAndRecoveryInterceptors(t *testing.T) {
	logging := loggingInterceptor()
	resp, err := logging(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/fraud.v1.FraudService/CheckFraud"}, func(context.Context, any) (any, error) {
		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)

	recovery := recoveryInterceptor()
	_, err = recovery(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/fraud.v1.FraudService/CheckFraud"}, func(context.Context, any) (any, error) {
		panic("boom")
	})
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	resp, err = recovery(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/fraud.v1.FraudService/CheckFraud"}, func(context.Context, any) (any, error) {
		return "safe", errors.New("handler error")
	})
	require.Equal(t, "safe", resp)
	require.EqualError(t, err, "handler error")
}
