package grpc

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ── Recovery interceptor ────────────────────────────────────────────────────

func TestRecoveryUnaryInterceptor_NoPanic(t *testing.T) {
	interceptor := recoveryUnaryInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}
	resp, err := interceptor(context.Background(), nil, info, handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRecoveryUnaryInterceptor_Panic(t *testing.T) {
	interceptor := recoveryUnaryInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	handler := func(ctx context.Context, req any) (any, error) {
		panic("something went wrong")
	}
	_, err := interceptor(context.Background(), nil, info, handler)
	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T", err)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("expected codes.Internal, got %v", st.Code())
	}
}

// ── RequestID interceptor ───────────────────────────────────────────────────

func TestRequestIDUnaryInterceptor_GeneratesID(t *testing.T) {
	interceptor := requestIDUnaryInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

	var capturedCtx context.Context
	handler := func(ctx context.Context, req any) (any, error) {
		capturedCtx = ctx
		return nil, nil
	}

	_, _ = interceptor(context.Background(), nil, info, handler)

	rid := getRequestID(capturedCtx)
	if rid == "" {
		t.Fatal("expected a request ID to be set in context, got empty string")
	}
}

func TestRequestIDUnaryInterceptor_PropagatesExistingID(t *testing.T) {
	interceptor := requestIDUnaryInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

	existingRID := "test-request-id-123"
	md := metadata.Pairs("x-request-id", existingRID)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	var capturedCtx context.Context
	handler := func(ctx context.Context, req any) (any, error) {
		capturedCtx = ctx
		return nil, nil
	}

	_, _ = interceptor(ctx, nil, info, handler)

	rid := getRequestID(capturedCtx)
	if rid != existingRID {
		t.Fatalf("expected request ID %q, got %q", existingRID, rid)
	}
}

// ── Logging interceptor ─────────────────────────────────────────────────────

func TestLoggingUnaryInterceptor_Success(t *testing.T) {
	interceptor := loggingUnaryInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	handler := func(ctx context.Context, req any) (any, error) {
		return "result", nil
	}
	resp, err := interceptor(context.Background(), nil, info, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "result" {
		t.Fatalf("expected 'result', got %v", resp)
	}
}

func TestLoggingUnaryInterceptor_Error(t *testing.T) {
	interceptor := loggingUnaryInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	_, err := interceptor(context.Background(), nil, info, handler)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── withRequestID helper ────────────────────────────────────────────────────

func TestWithRequestID_EmptyMetadata(t *testing.T) {
	ctx, rid := withRequestID(context.Background())
	if rid == "" {
		t.Fatal("expected a generated request ID")
	}
	if getRequestID(ctx) != rid {
		t.Fatal("request ID not stored in context")
	}
}

func TestWithRequestID_FromCorrelationID(t *testing.T) {
	md := metadata.Pairs("x-correlation-id", "corr-456")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	_, rid := withRequestID(ctx)
	if rid != "corr-456" {
		t.Fatalf("expected 'corr-456', got %q", rid)
	}
}

// ── Stream recovery interceptor ─────────────────────────────────────────────

func TestRecoveryStreamInterceptor_NoPanic(t *testing.T) {
	interceptor := recoveryStreamInterceptor()
	info := &grpc.StreamServerInfo{FullMethod: "/test.Service/Stream"}
	handler := func(srv any, ss grpc.ServerStream) error {
		return nil
	}
	err := interceptor(nil, &mockServerStream{ctx: context.Background()}, info, handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRecoveryStreamInterceptor_Panic(t *testing.T) {
	interceptor := recoveryStreamInterceptor()
	info := &grpc.StreamServerInfo{FullMethod: "/test.Service/Stream"}
	handler := func(srv any, ss grpc.ServerStream) error {
		panic("stream panic")
	}
	err := interceptor(nil, &mockServerStream{ctx: context.Background()}, info, handler)
	if err == nil {
		t.Fatal("expected error from panic recovery")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.Internal {
		t.Fatalf("expected codes.Internal, got %v", st.Code())
	}
}

// ── Mock helpers ────────────────────────────────────────────────────────────

type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context { return m.ctx }
func (m *mockServerStream) SendMsg(msg any) error    { return nil }
func (m *mockServerStream) RecvMsg(msg any) error    { return errors.New("EOF") }
