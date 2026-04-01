package middleware

import (
	"context"
	"errors"
	"testing"

	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeAuthzClient struct {
	resp    *authzservicev1.CheckAccessResponse
	err     error
	lastReq *authzservicev1.CheckAccessRequest
}

func (f *fakeAuthzClient) CheckAccess(_ context.Context, req *authzservicev1.CheckAccessRequest, _ ...grpc.CallOption) (*authzservicev1.CheckAccessResponse, error) {
	f.lastReq = req
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func TestOrderAuthZInterceptorScenarios(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		method     string
		client     *fakeAuthzClient
		wantCode   codes.Code
		wantCalled bool
	}{
		{
			name:       "skip health",
			ctx:        context.Background(),
			method:     "/grpc.health.v1.Health/Check",
			client:     &fakeAuthzClient{},
			wantCode:   codes.OK,
			wantCalled: true,
		},
		{
			name:     "missing metadata",
			ctx:      context.Background(),
			method:   "/insuretech.orders.services.v1.OrderService/GetOrder",
			client:   &fakeAuthzClient{},
			wantCode: codes.Unauthenticated,
		},
		{
			name: "missing user",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-portal", "PORTAL_SYSTEM",
			)),
			method:   "/insuretech.orders.services.v1.OrderService/GetOrder",
			client:   &fakeAuthzClient{},
			wantCode: codes.Unauthenticated,
		},
		{
			name: "b2b missing org",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-user-id", "user-1",
				"x-portal", "PORTAL_B2B",
			)),
			method:   "/insuretech.orders.services.v1.OrderService/GetOrder",
			client:   &fakeAuthzClient{},
			wantCode: codes.PermissionDenied,
		},
		{
			name: "denied",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-user-id", "user-1",
				"x-portal", "PORTAL_SYSTEM",
			)),
			method:   "/insuretech.orders.services.v1.OrderService/CancelOrder",
			client:   &fakeAuthzClient{resp: &authzservicev1.CheckAccessResponse{Allowed: false}},
			wantCode: codes.PermissionDenied,
		},
		{
			name: "authz error fails open",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-user-id", "user-1",
				"x-portal", "PORTAL_SYSTEM",
			)),
			method:     "/insuretech.orders.services.v1.OrderService/GetOrder",
			client:     &fakeAuthzClient{err: errors.New("unavailable")},
			wantCode:   codes.OK,
			wantCalled: true,
		},
		{
			name: "allowed",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-user-id", "user-1",
				"x-portal", "PORTAL_B2B",
				"x-business-id", "org-1",
				"x-tenant-id", "tenant-1",
			)),
			method:     "/insuretech.orders.services.v1.OrderService/CreateOrder",
			client:     &fakeAuthzClient{resp: &authzservicev1.CheckAccessResponse{Allowed: true}},
			wantCode:   codes.OK,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := NewOrderAuthZInterceptor(tt.client).UnaryServerInterceptor()
			called := false
			_, err := interceptor(tt.ctx, "req", &grpc.UnaryServerInfo{FullMethod: tt.method}, func(ctx context.Context, req interface{}) (interface{}, error) {
				called = true
				return "ok", nil
			})
			if status.Code(err) != tt.wantCode {
				t.Fatalf("code = %v, want %v", status.Code(err), tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Fatalf("handler called = %v, want %v", called, tt.wantCalled)
			}
		})
	}
}

func TestOrderAuthZHelpers(t *testing.T) {
	resource, action := mapOrderMethodToResourceAction("/insuretech.orders.services.v1.OrderService/GetOrder")
	if resource != "svc:order/*" || action != "GET" {
		t.Fatalf("unexpected read mapping: %q %q", resource, action)
	}
	resource, action = mapOrderMethodToResourceAction("/insuretech.orders.services.v1.OrderService/CreateOrder")
	if resource != "svc:order/*" || action != "POST" {
		t.Fatalf("unexpected create mapping: %q %q", resource, action)
	}
	resource, action = mapOrderMethodToResourceAction("/insuretech.orders.services.v1.OrderService/CancelOrder")
	if resource != "svc:order/*" || action != "DELETE" {
		t.Fatalf("unexpected cancel mapping: %q %q", resource, action)
	}
	resource, action = mapOrderMethodToResourceAction("/insuretech.orders.services.v1.OrderService/Unknown")
	if resource != "" || action != "" {
		t.Fatalf("expected unknown method skip")
	}

	if got := resolveOrderAuthzDomain(metadata.Pairs("x-portal", "PORTAL_SYSTEM"), ""); got != "system:root" {
		t.Fatalf("system domain = %q", got)
	}
	if got := resolveOrderAuthzDomain(metadata.Pairs("x-portal", "PORTAL_B2B", "x-tenant-id", "tenant-1"), "org-1"); got != "b2b:org-1" {
		t.Fatalf("b2b domain = %q", got)
	}
	if got := resolveOrderAuthzDomain(metadata.Pairs("x-portal", "PORTAL_AGENT", "x-tenant-id", "tenant-1"), ""); got != "agent:tenant-1" {
		t.Fatalf("agent domain = %q", got)
	}
	if !isSkipMethod("/grpc.health.v1.Health/Check") || isSkipMethod("/orders") {
		t.Fatalf("unexpected skip method result")
	}
	if isOrderBootstrapMethod("/anything") {
		t.Fatalf("expected no bootstrap methods")
	}
}
