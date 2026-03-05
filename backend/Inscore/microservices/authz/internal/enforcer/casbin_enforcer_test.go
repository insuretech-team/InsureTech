package enforcer

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/stretchr/testify/require"
)

func newMemoryEnforcer(t *testing.T) *CasbinEnforcer {
	t.Helper()
	m, err := model.NewModelFromString(builtinModel)
	require.NoError(t, err)
	e, err := casbin.NewSyncedEnforcer(m)
	require.NoError(t, err)
	return &CasbinEnforcer{enforcer: e}
}

func TestCasbinEnforcer_EnforceAndPolicyOps(t *testing.T) {
	ce := newMemoryEnforcer(t)
	ctx := context.Background()

	allowed, rule, err := ce.Enforce(ctx, "user:u1", "system:root", "svc:user/get", "GET")
	require.NoError(t, err)
	require.False(t, allowed)
	require.Empty(t, rule)

	require.NoError(t, ce.AddRoleForUserInDomain("user:u1", "role:admin", "system:root"))
	require.NoError(t, ce.AddPolicy("role:admin", "system:root", "svc:user/*", "GET", "allow"))
	allowed, rule, err = ce.Enforce(ctx, "user:u1", "system:root", "svc:user/get", "GET")
	require.NoError(t, err)
	require.True(t, allowed)
	require.Empty(t, rule)

	require.NoError(t, ce.AddPolicy("user:u1", "system:root", "svc:profile/*", "GET", "allow"))
	allowed, rule, err = ce.Enforce(ctx, "user:u1", "system:root", "svc:profile/get", "GET")
	require.NoError(t, err)
	require.True(t, allowed)
	require.NotEmpty(t, rule)

	roles, err := ce.GetRolesForUserInDomain("user:u1", "system:root")
	require.NoError(t, err)
	require.Contains(t, roles, "role:admin")

	perms, err := ce.GetPermissionsForUserInDomain("role:admin", "system:root")
	require.NoError(t, err)
	require.NotEmpty(t, perms)

	require.Error(t, ce.AddPolicy("role:admin", "system:root", "svc:user/*", "GET", "allow"))
	require.NoError(t, ce.RemovePolicy("role:admin", "system:root", "svc:user/*", "GET"))
	require.NoError(t, ce.DeleteRoleForUserInDomain("user:u1", "role:admin", "system:root"))
	require.Error(t, ce.DeleteRoleForUserInDomain("user:u1", "role:admin", "system:root"))
}

func TestCasbinEnforcer_NewErrorAndReloadPaths(t *testing.T) {
	require.Panics(t, func() {
		_, _ = New(nil, "")
	})

	ce := newMemoryEnforcer(t)
	require.NotPanics(t, func() {
		ce.StartAutoReload(1)
	})
	require.Panics(t, func() {
		_ = ce.InvalidateCache()
	})
}

func TestCasbinEnforcer_B2BPathPattern(t *testing.T) {
	ce := newMemoryEnforcer(t)
	ctx := context.Background()

	require.NoError(t, ce.AddRoleForUserInDomain("user:sys-1", "role:support", "system:root"))
	require.NoError(t, ce.AddPolicy("role:support", "system:root", "svc:b2b/*", "GET", "allow"))

	allowed, _, err := ce.Enforce(ctx, "user:sys-1", "system:root", "svc:b2b/b2b/employees", "GET")
	require.NoError(t, err)
	require.True(t, allowed)
}
