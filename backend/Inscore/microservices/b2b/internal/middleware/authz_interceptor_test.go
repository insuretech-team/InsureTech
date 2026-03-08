package middleware

import (
	"context"
	"testing"

	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ── Test doubles ─────────────────────────────────────────────────────────────

type fakeAuthZClient struct{}

func (f *fakeAuthZClient) CheckAccess(_ context.Context, _ *authzservicev1.CheckAccessRequest) (*authzservicev1.CheckAccessResponse, error) {
	return &authzservicev1.CheckAccessResponse{Allowed: true}, nil
}

type captureAuthZClient struct {
	last    *authzservicev1.CheckAccessRequest
	allowed bool
}

func (c *captureAuthZClient) CheckAccess(_ context.Context, req *authzservicev1.CheckAccessRequest) (*authzservicev1.CheckAccessResponse, error) {
	c.last = req
	return &authzservicev1.CheckAccessResponse{Allowed: c.allowed}, nil
}

func newCapture(allowed bool) *captureAuthZClient { return &captureAuthZClient{allowed: allowed} }

func okHandler(_ context.Context, _ interface{}) (interface{}, error) { return "ok", nil }

func makeCtx(pairs ...string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(pairs...))
}

func runInterceptor(t *testing.T, client AuthZClient, method string, ctx context.Context) (interface{}, error) {
	t.Helper()
	interceptor := NewAuthZInterceptor(client).UnaryServerInterceptor()
	return interceptor(ctx, struct{}{}, &grpc.UnaryServerInfo{FullMethod: method}, okHandler)
}

// ── Original tests (updated for correct behaviour) ────────────────────────────

// ListOrganisations called by a NON-system user with NO org context must now be
// PermissionDenied — the old bypass was the bug we fixed.
func TestUnaryServerInterceptor_ListOrganisations_NonSystemNoOrg_Denied(t *testing.T) {
	ctx := makeCtx("x-user-id", "user-1", "x-business-id", "")
	_, err := runInterceptor(t, &fakeAuthZClient{},
		"/insuretech.b2b.services.v1.B2BService/ListOrganisations", ctx)
	require.Error(t, err)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
}

// System portal user with NO org ID should succeed (domain = system:root).
func TestUnaryServerInterceptor_ListOrgMembers_UsesSystemRootForSystemUsers(t *testing.T) {
	client := newCapture(true)
	ctx := makeCtx(
		"x-user-id", "user-1",
		"x-business-id", "org-1",
		"x-portal", "system",
		"x-tenant-id", "tenant-1",
	)
	_, err := runInterceptor(t, client,
		"/insuretech.b2b.services.v1.B2BService/ListOrgMembers", ctx)
	require.NoError(t, err)
	require.NotNil(t, client.last)
	require.Equal(t, "system:root", client.last.Domain)
	require.Equal(t, "svc:b2b/*", client.last.Object)
	require.Equal(t, "GET", client.last.Action)
}

// B2B portal user with org ID → domain must be b2b:{org_id}.
func TestUnaryServerInterceptor_AssignOrgAdmin_UsesB2BDomainForPartnerUsers(t *testing.T) {
	client := newCapture(true)
	ctx := makeCtx("x-user-id", "user-2", "x-business-id", "org-2", "x-portal", "b2b")
	_, err := runInterceptor(t, client,
		"/insuretech.b2b.services.v1.B2BService/AssignOrgAdmin", ctx)
	require.NoError(t, err)
	require.NotNil(t, client.last)
	require.Equal(t, "b2b:org-2", client.last.Domain)
	require.Equal(t, "svc:b2b/*", client.last.Object)
	require.Equal(t, "POST", client.last.Action)
}

// ── New tests ─────────────────────────────────────────────────────────────────

// Missing metadata → Unauthenticated.
func TestUnaryServerInterceptor_MissingMetadata_Unauthenticated(t *testing.T) {
	interceptor := NewAuthZInterceptor(&fakeAuthZClient{}).UnaryServerInterceptor()
	_, err := interceptor(context.Background(), struct{}{},
		&grpc.UnaryServerInfo{FullMethod: "/svc/Any"}, okHandler)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

// Missing user_id → Unauthenticated.
func TestUnaryServerInterceptor_MissingUserID_Unauthenticated(t *testing.T) {
	ctx := makeCtx("x-business-id", "org-1")
	_, err := runInterceptor(t, &fakeAuthZClient{}, "/svc/Any", ctx)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

// ResolveMyOrganisation always passes through — no org needed, no Casbin check.
func TestUnaryServerInterceptor_ResolveMyOrganisation_AlwaysPasses(t *testing.T) {
	client := newCapture(false) // Casbin would deny — but should never be called.
	ctx := makeCtx("x-user-id", "user-1") // no portal, no org
	resp, err := runInterceptor(t, client,
		"/insuretech.b2b.services.v1.B2BService/ResolveMyOrganisation", ctx)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.Nil(t, client.last, "Casbin should NOT be called for ResolveMyOrganisation")
}

// Super admin (PORTAL_SYSTEM, no org) calling org management → system:root, allowed.
func TestUnaryServerInterceptor_SuperAdmin_NoOrg_OrgManagement_Allowed(t *testing.T) {
	orgMgmtMethods := []struct{ method, action string }{
		{"/insuretech.b2b.services.v1.B2BService/ListOrganisations", "GET"},
		{"/insuretech.b2b.services.v1.B2BService/CreateOrganisation", "POST"},
		{"/insuretech.b2b.services.v1.B2BService/GetOrganisation", "GET"},
		{"/insuretech.b2b.services.v1.B2BService/UpdateOrganisation", "PATCH"},
		{"/insuretech.b2b.services.v1.B2BService/DeleteOrganisation", "DELETE"},
		{"/insuretech.b2b.services.v1.B2BService/AssignOrgAdmin", "POST"},
		{"/insuretech.b2b.services.v1.B2BService/AddOrgMember", "POST"},
		{"/insuretech.b2b.services.v1.B2BService/RemoveOrgMember", "DELETE"},
	}
	for _, tc := range orgMgmtMethods {
		tc := tc
		t.Run(tc.method, func(t *testing.T) {
			client := newCapture(true)
			ctx := makeCtx("x-user-id", "super-1", "x-portal", "PORTAL_SYSTEM")
			_, err := runInterceptor(t, client, tc.method, ctx)
			require.NoError(t, err)
			require.NotNil(t, client.last)
			require.Equal(t, "system:root", client.last.Domain)
			require.Equal(t, "svc:b2b/*", client.last.Object)
			require.Equal(t, tc.action, client.last.Action)
		})
	}
}

// Super admin denied by Casbin → PermissionDenied propagated.
func TestUnaryServerInterceptor_SuperAdmin_CasbinDenies_PermissionDenied(t *testing.T) {
	client := newCapture(false)
	ctx := makeCtx("x-user-id", "super-1", "x-portal", "PORTAL_SYSTEM")
	_, err := runInterceptor(t, client,
		"/insuretech.b2b.services.v1.B2BService/CreateOrganisation", ctx)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
}

// B2B admin with org ID — all b2b methods map to b2b:{org_id} domain.
func TestUnaryServerInterceptor_B2BAdmin_WithOrg_CorrectDomain(t *testing.T) {
	methods := []struct{ method, action string }{
		{"/insuretech.b2b.services.v1.B2BService/ListDepartments", "GET"},
		{"/insuretech.b2b.services.v1.B2BService/CreateDepartment", "POST"},
		{"/insuretech.b2b.services.v1.B2BService/UpdateDepartment", "PATCH"},
		{"/insuretech.b2b.services.v1.B2BService/DeleteDepartment", "DELETE"},
		{"/insuretech.b2b.services.v1.B2BService/ListEmployees", "GET"},
		{"/insuretech.b2b.services.v1.B2BService/AddEmployee", "POST"},
	}
	for _, tc := range methods {
		tc := tc
		t.Run(tc.method, func(t *testing.T) {
			client := newCapture(true)
			ctx := makeCtx("x-user-id", "b2b-admin-1", "x-business-id", "org-abc", "x-portal", "PORTAL_B2B")
			_, err := runInterceptor(t, client, tc.method, ctx)
			require.NoError(t, err)
			require.NotNil(t, client.last)
			require.Equal(t, "b2b:org-abc", client.last.Domain)
			require.Equal(t, "svc:b2b/*", client.last.Object)
			require.Equal(t, tc.action, client.last.Action)
		})
	}
}

// Non-system user without org context → PermissionDenied for all protected methods.
func TestUnaryServerInterceptor_NonSystem_NoOrg_Protected_Denied(t *testing.T) {
	protectedMethods := []string{
		"/insuretech.b2b.services.v1.B2BService/ListDepartments",
		"/insuretech.b2b.services.v1.B2BService/CreateDepartment",
		"/insuretech.b2b.services.v1.B2BService/ListOrganisations",
		"/insuretech.b2b.services.v1.B2BService/AssignOrgAdmin",
		"/insuretech.b2b.services.v1.B2BService/ListEmployees",
	}
	for _, method := range protectedMethods {
		method := method
		t.Run(method, func(t *testing.T) {
			ctx := makeCtx("x-user-id", "user-x", "x-portal", "PORTAL_B2B")
			_, err := runInterceptor(t, &fakeAuthZClient{}, method, ctx)
			require.Equal(t, codes.PermissionDenied, status.Code(err))
		})
	}
}

// System portal with PORTAL_ prefix in header value is normalised correctly.
func TestUnaryServerInterceptor_SystemPortalPrefix_Normalised(t *testing.T) {
	client := newCapture(true)
	ctx := makeCtx("x-user-id", "super-1", "x-portal", "PORTAL_SYSTEM")
	_, err := runInterceptor(t, client,
		"/insuretech.b2b.services.v1.B2BService/ListOrganisations", ctx)
	require.NoError(t, err)
	require.NotNil(t, client.last)
	require.Equal(t, "system:root", client.last.Domain)
}

// Unknown method falls through without Casbin check (returns ok, no error).
func TestUnaryServerInterceptor_UnknownMethod_PassesThrough(t *testing.T) {
	client := newCapture(false) // if called, would deny
	ctx := makeCtx("x-user-id", "user-1", "x-business-id", "org-1", "x-portal", "b2b")
	resp, err := runInterceptor(t, client, "/insuretech.b2b.services.v1.B2BService/UnknownRPC", ctx)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.Nil(t, client.last, "Casbin should NOT be called for unknown methods")
}
