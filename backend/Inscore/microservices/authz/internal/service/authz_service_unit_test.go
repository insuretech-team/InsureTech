package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/cache"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
)

type addPolicyCall struct {
	subject string
	domain  string
	object  string
	action  string
	effect  string
}

type fakeUnitEnforcer struct {
	roles             []string
	rolesErr          error
	addPolicyCalls    []addPolicyCall
	addPolicyErrByObj map[string]error
}

func (f *fakeUnitEnforcer) Enforce(context.Context, string, string, string, string) (bool, string, error) {
	return false, "", nil
}

func (f *fakeUnitEnforcer) AddPolicy(sub, dom, obj, act, effect string) error {
	f.addPolicyCalls = append(f.addPolicyCalls, addPolicyCall{
		subject: sub,
		domain:  dom,
		object:  obj,
		action:  act,
		effect:  effect,
	})
	if f.addPolicyErrByObj != nil {
		if err, ok := f.addPolicyErrByObj[obj]; ok {
			return err
		}
	}
	return nil
}

func (f *fakeUnitEnforcer) RemovePolicy(string, string, string, string) error { return nil }

func (f *fakeUnitEnforcer) AddRoleForUserInDomain(string, string, string) error { return nil }

func (f *fakeUnitEnforcer) DeleteRoleForUserInDomain(string, string, string) error { return nil }

func (f *fakeUnitEnforcer) GetRolesForUserInDomain(string, string) ([]string, error) {
	if f.rolesErr != nil {
		return nil, f.rolesErr
	}
	return f.roles, nil
}

func (f *fakeUnitEnforcer) GetPermissionsForUserInDomain(string, string) ([][]string, error) {
	return nil, nil
}

func (f *fakeUnitEnforcer) InvalidateCache() error { return nil }

type fakeUnitPolicyRepo struct {
	listPages      map[int][]*authzentityv1.PolicyRule
	listErrByOff   map[int]error
	listCalls      []int
	createCalls    []*authzentityv1.PolicyRule
	createErrByObj map[string]error
}

func (f *fakeUnitPolicyRepo) Create(_ context.Context, pr *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error) {
	if f.createErrByObj != nil {
		if err, ok := f.createErrByObj[pr.Object]; ok {
			return nil, err
		}
	}
	f.createCalls = append(f.createCalls, pr)
	return pr, nil
}

func (f *fakeUnitPolicyRepo) Update(context.Context, *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error) {
	return nil, nil
}

func (f *fakeUnitPolicyRepo) SoftDelete(context.Context, string) error {
	return nil
}

func (f *fakeUnitPolicyRepo) List(_ context.Context, _ string, _ bool, _ int, offset int) ([]*authzentityv1.PolicyRule, error) {
	f.listCalls = append(f.listCalls, offset)
	if f.listErrByOff != nil {
		if err, ok := f.listErrByOff[offset]; ok {
			return nil, err
		}
	}
	return f.listPages[offset], nil
}

func TestAuthZService_SetPermissionCache(t *testing.T) {
	svc := &AuthZService{}
	permCache := cache.NewPermissionCache(nil, time.Minute)

	svc.SetPermissionCache(permCache)

	require.Same(t, permCache, svc.permCache)
}

func TestAuthZService_ListPolicyRulesByDomain_PaginatesAndErrors(t *testing.T) {
	page := make([]*authzentityv1.PolicyRule, 500)
	for i := range page {
		page[i] = &authzentityv1.PolicyRule{PolicyId: "p"}
	}

	repo := &fakeUnitPolicyRepo{
		listPages: map[int][]*authzentityv1.PolicyRule{
			0:   page,
			500: {{PolicyId: "last"}},
		},
	}
	svc := &AuthZService{policyRepo: repo}

	policies, err := svc.listPolicyRulesByDomain(context.Background(), "b2b:root")
	require.NoError(t, err)
	require.Len(t, policies, 501)
	require.Equal(t, []int{0, 500}, repo.listCalls)

	_, err = (&AuthZService{
		policyRepo: &fakeUnitPolicyRepo{listErrByOff: map[int]error{0: errors.New("boom")}},
	}).listPolicyRulesByDomain(context.Background(), "b2b:root")
	require.Error(t, err)
}

func TestAuthZService_CheckB2BRootDomainFallback(t *testing.T) {
	ctx := context.Background()
	req := &authzservicev1.CheckAccessRequest{
		Domain: "b2b:tenant-1",
		Object: "svc:b2b/employees",
		Action: "GET",
	}

	t.Run("allow", func(t *testing.T) {
		svc := &AuthZService{
			enforcer: &fakeUnitEnforcer{roles: []string{"role:partner_admin"}},
			policyRepo: &fakeUnitPolicyRepo{
				listPages: map[int][]*authzentityv1.PolicyRule{
					0: {
						nil,
						{Subject: "role:other", Object: "svc:b2b/*", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW},
						{Subject: "role:partner_admin", Object: "svc:b2b/*", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW},
					},
				},
			},
		}

		allowed, matchedRule, err := svc.checkB2BRootDomainFallback(ctx, "user:u1", req)
		require.NoError(t, err)
		require.True(t, allowed)
		require.Contains(t, matchedRule, "root-fallback:")
	})

	t.Run("deny", func(t *testing.T) {
		svc := &AuthZService{
			enforcer: &fakeUnitEnforcer{roles: []string{"role:partner_admin"}},
			policyRepo: &fakeUnitPolicyRepo{
				listPages: map[int][]*authzentityv1.PolicyRule{
					0: {
						{Subject: "role:partner_admin", Object: "svc:b2b/*", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_DENY},
					},
				},
			},
		}

		allowed, matchedRule, err := svc.checkB2BRootDomainFallback(ctx, "user:u1", req)
		require.NoError(t, err)
		require.False(t, allowed)
		require.Equal(t, "root-fallback-deny", matchedRule)
	})

	t.Run("no roles", func(t *testing.T) {
		svc := &AuthZService{
			enforcer:   &fakeUnitEnforcer{},
			policyRepo: &fakeUnitPolicyRepo{},
		}

		allowed, matchedRule, err := svc.checkB2BRootDomainFallback(ctx, "user:u1", req)
		require.NoError(t, err)
		require.False(t, allowed)
		require.Empty(t, matchedRule)
	})

	t.Run("role lookup error", func(t *testing.T) {
		svc := &AuthZService{
			enforcer:   &fakeUnitEnforcer{rolesErr: errors.New("role lookup failed")},
			policyRepo: &fakeUnitPolicyRepo{},
		}

		_, _, err := svc.checkB2BRootDomainFallback(ctx, "user:u1", req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list domain roles")
	})

	t.Run("policy list error", func(t *testing.T) {
		svc := &AuthZService{
			enforcer: &fakeUnitEnforcer{roles: []string{"role:partner_admin"}},
			policyRepo: &fakeUnitPolicyRepo{
				listErrByOff: map[int]error{0: errors.New("list failed")},
			},
		}

		_, _, err := svc.checkB2BRootDomainFallback(ctx, "user:u1", req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list root policies")
	})
}

func TestAuthZService_EnsureScopedRolePolicies(t *testing.T) {
	ctx := context.Background()

	t.Run("skip non scoped roles", func(t *testing.T) {
		enforcer := &fakeUnitEnforcer{}
		repo := &fakeUnitPolicyRepo{}
		svc := &AuthZService{enforcer: enforcer, policyRepo: repo}

		require.NoError(t, svc.ensureScopedRolePolicies(ctx, nil, "b2b:tenant-1"))
		require.NoError(t, svc.ensureScopedRolePolicies(ctx, &authzentityv1.Role{Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM}, "b2b:tenant-1"))
		require.NoError(t, svc.ensureScopedRolePolicies(ctx, &authzentityv1.Role{Name: "partner_admin", Portal: authzentityv1.Portal_PORTAL_B2B}, "b2b:root"))
		require.Empty(t, enforcer.addPolicyCalls)
		require.Empty(t, repo.createCalls)
	})

	t.Run("copies matching root policies and tolerates duplicates", func(t *testing.T) {
		enforcer := &fakeUnitEnforcer{
			addPolicyErrByObj: map[string]error{
				"svc:b2b/departments": errors.New("rule already exists"),
			},
		}
		repo := &fakeUnitPolicyRepo{
			listPages: map[int][]*authzentityv1.PolicyRule{
				0: {
					nil,
					{Subject: "role:other", Object: "svc:b2b/*", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW},
					{Subject: "role:partner_admin", Object: "svc:b2b/departments", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW},
					{Subject: "role:partner_admin", Object: "svc:b2b/orders", Action: "POST", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_DENY, CreatedBy: "creator"},
				},
			},
			createErrByObj: map[string]error{
				"svc:b2b/orders": errors.New("duplicate"),
			},
		}
		svc := &AuthZService{enforcer: enforcer, policyRepo: repo}

		err := svc.ensureScopedRolePolicies(ctx, &authzentityv1.Role{
			Name:   "partner_admin",
			Portal: authzentityv1.Portal_PORTAL_B2B,
		}, "b2b:tenant-1")
		require.NoError(t, err)
		require.Len(t, enforcer.addPolicyCalls, 2)
		require.Equal(t, "allow", enforcer.addPolicyCalls[0].effect)
		require.Equal(t, "deny", enforcer.addPolicyCalls[1].effect)
		require.Len(t, repo.createCalls, 1)
		require.Equal(t, "system", repo.createCalls[0].CreatedBy)
		require.Equal(t, "b2b:tenant-1", repo.createCalls[0].Domain)
	})

	t.Run("propagates add policy errors", func(t *testing.T) {
		svc := &AuthZService{
			enforcer: &fakeUnitEnforcer{
				addPolicyErrByObj: map[string]error{"svc:b2b/departments": errors.New("boom")},
			},
			policyRepo: &fakeUnitPolicyRepo{
				listPages: map[int][]*authzentityv1.PolicyRule{
					0: {
						{Subject: "role:partner_admin", Object: "svc:b2b/departments", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW},
					},
				},
			},
		}

		err := svc.ensureScopedRolePolicies(ctx, &authzentityv1.Role{
			Name:   "partner_admin",
			Portal: authzentityv1.Portal_PORTAL_B2B,
		}, "b2b:tenant-1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "copy scoped casbin policy")
	})

	t.Run("propagates create errors", func(t *testing.T) {
		svc := &AuthZService{
			enforcer: &fakeUnitEnforcer{},
			policyRepo: &fakeUnitPolicyRepo{
				listPages: map[int][]*authzentityv1.PolicyRule{
					0: {
						{Subject: "role:partner_admin", Object: "svc:b2b/departments", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW},
					},
				},
				createErrByObj: map[string]error{"svc:b2b/departments": errors.New("persist failed")},
			},
		}

		err := svc.ensureScopedRolePolicies(ctx, &authzentityv1.Role{
			Name:   "partner_admin",
			Portal: authzentityv1.Portal_PORTAL_B2B,
		}, "b2b:tenant-1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "persist scoped policy")
	})
}

func TestAuthZService_HelperPredicates(t *testing.T) {
	require.True(t, isAlreadyExistsErr(errors.New("duplicate key")))
	require.True(t, isAlreadyExistsErr(errors.New("rule already exists")))
	require.False(t, isAlreadyExistsErr(errors.New("boom")))
	require.Equal(t, "value", nonEmpty("", "  ", "value", "fallback"))
	require.Equal(t, "", nonEmpty("", "  "))
}
