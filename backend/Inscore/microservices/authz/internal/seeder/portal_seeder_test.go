package seeder

import (
	"context"
	"errors"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/domain"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeRoleRepo struct {
	roles     map[string]*authzentityv1.Role
	createErr error
}

func (f *fakeRoleRepo) Create(_ context.Context, role *authzentityv1.Role) (*authzentityv1.Role, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	if f.roles == nil {
		f.roles = map[string]*authzentityv1.Role{}
	}
	k := role.Name + "|" + role.Portal.String()
	if role.RoleId == "" {
		role.RoleId = "rid-" + k
	}
	f.roles[k] = role
	return role, nil
}
func (f *fakeRoleRepo) GetByID(context.Context, string) (*authzentityv1.Role, error) {
	return nil, errors.New("not found")
}
func (f *fakeRoleRepo) GetByNameAndPortal(_ context.Context, name string, portal authzentityv1.Portal) (*authzentityv1.Role, error) {
	if f.roles == nil {
		return nil, errors.New("not found")
	}
	if r, ok := f.roles[name+"|"+portal.String()]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}
func (f *fakeRoleRepo) Update(context.Context, *authzentityv1.Role) (*authzentityv1.Role, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRoleRepo) SoftDelete(context.Context, string) error {
	return errors.New("not implemented")
}
func (f *fakeRoleRepo) List(context.Context, authzentityv1.Portal, bool, int, int) ([]*authzentityv1.Role, error) {
	return nil, nil
}

type fakePolicyRepo struct {
	createCount int
	dupByObject map[string]bool
	createErr   error
}

func (f *fakePolicyRepo) Create(_ context.Context, pr *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	if f.dupByObject != nil && f.dupByObject[pr.Object] {
		return nil, errors.New("duplicate")
	}
	f.createCount++
	return pr, nil
}
func (f *fakePolicyRepo) Update(context.Context, *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error) {
	return nil, errors.New("not implemented")
}
func (f *fakePolicyRepo) SoftDelete(context.Context, string) error {
	return errors.New("not implemented")
}
func (f *fakePolicyRepo) List(context.Context, string, bool, int, int) ([]*authzentityv1.PolicyRule, error) {
	return nil, nil
}

type fakeEnforcer struct {
	domain.EnforcerIface
	addPolicyCalls int
}

func (f *fakeEnforcer) AddPolicy(sub, dom, obj, act, effect string) error {
	f.addPolicyCalls++
	return nil
}

type fakePortalCfgRepo struct {
	upserts int
}

func (f *fakePortalCfgRepo) GetByPortal(context.Context, authzentityv1.Portal) (*authzentityv1.PortalConfig, error) {
	return nil, errors.New("not implemented")
}
func (f *fakePortalCfgRepo) Upsert(_ context.Context, pc *authzentityv1.PortalConfig) (*authzentityv1.PortalConfig, error) {
	f.upserts++
	return pc, nil
}

type fakeTokenCfgRepo struct {
	active    *authzentityv1.TokenConfig
	createErr error
	created   int
}

func (f *fakeTokenCfgRepo) GetActive(context.Context) (*authzentityv1.TokenConfig, error) {
	return f.active, nil
}
func (f *fakeTokenCfgRepo) List(context.Context) ([]*authzentityv1.TokenConfig, error) {
	return nil, nil
}
func (f *fakeTokenCfgRepo) Create(_ context.Context, cfg *authzentityv1.TokenConfig) (*authzentityv1.TokenConfig, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	f.created++
	f.active = cfg
	return cfg, nil
}

func TestSeedPortal_DryRunAndLive(t *testing.T) {
	s := New(&fakeRoleRepo{}, &fakePolicyRepo{}, &fakeEnforcer{}, &fakePortalCfgRepo{}, &fakeTokenCfgRepo{}, nil, zap.NewNop())
	ctx := context.Background()

	res, err := s.SeedPortal(ctx, authzentityv1.Portal_PORTAL_SYSTEM, GlobalTenantID, true)
	require.NoError(t, err)
	require.Greater(t, res.RolesSeeded, 0)
	require.Equal(t, 0, res.PoliciesSeeded)

	rr := &fakeRoleRepo{}
	pr := &fakePolicyRepo{}
	ef := &fakeEnforcer{}
	s = New(rr, pr, ef, &fakePortalCfgRepo{}, &fakeTokenCfgRepo{}, nil, zap.NewNop())
	res, err = s.SeedPortal(ctx, authzentityv1.Portal_PORTAL_AGENT, GlobalTenantID, false)
	require.NoError(t, err)
	require.Greater(t, res.RolesSeeded, 0)
	require.Greater(t, pr.createCount, 0)
	require.Greater(t, ef.addPolicyCalls, 0)

	res, err = s.SeedPortal(ctx, authzentityv1.Portal_PORTAL_UNSPECIFIED, GlobalTenantID, false)
	require.NoError(t, err)
	require.Equal(t, 0, res.RolesSeeded)
}

func TestSeedAllPortals_ConfigsAndToken(t *testing.T) {
	t.Setenv("JWT_KEY_ID", "kid-1")
	t.Setenv("JWT_PUBLIC_KEY_PEM", "pem")
	t.Setenv("JWT_PRIVATE_KEY_REF", "vault://key")

	pc := &fakePortalCfgRepo{}
	tc := &fakeTokenCfgRepo{}
	s := New(&fakeRoleRepo{}, &fakePolicyRepo{}, &fakeEnforcer{}, pc, tc, nil, zap.NewNop())
	require.NoError(t, s.SeedAllPortals(context.Background()))
	require.Equal(t, 6, pc.upserts)
	require.Equal(t, 1, tc.created)

	tc.active = &authzentityv1.TokenConfig{Kid: "kid-1", CreatedAt: timestamppb.Now()}
	require.NoError(t, s.SeedTokenConfig(context.Background()))
	require.Equal(t, 1, tc.created)
}

func TestSeedTokenConfig_MissingPublicKeyAndCreateError(t *testing.T) {
	t.Setenv("JWT_PUBLIC_KEY_PEM", "")
	tc := &fakeTokenCfgRepo{createErr: errors.New("create fail")}
	s := New(&fakeRoleRepo{}, &fakePolicyRepo{}, &fakeEnforcer{}, &fakePortalCfgRepo{}, tc, nil, zap.NewNop())
	require.NoError(t, s.SeedTokenConfig(context.Background()))

	t.Setenv("JWT_PUBLIC_KEY_PEM", "pem")
	require.Error(t, s.SeedTokenConfig(context.Background()))
}

func TestEffectFromString(t *testing.T) {
	require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_DENY, effectFromString("deny"))
	require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW, effectFromString("allow"))
}

func TestPortalDomainKey(t *testing.T) {
	require.Equal(t, "system:root", portalDomainKey(authzentityv1.Portal_PORTAL_SYSTEM, ""))
	require.Equal(t, "b2b:root", portalDomainKey(authzentityv1.Portal_PORTAL_B2B, "root"))
	require.Equal(t, "business:tenant-1", portalDomainKey(authzentityv1.Portal_PORTAL_BUSINESS, "tenant-1"))
}

func TestSeedPortal_FailsOnNonDuplicatePolicyError(t *testing.T) {
	ctx := context.Background()
	rr := &fakeRoleRepo{}
	pr := &fakePolicyRepo{createErr: errors.New("invalid input syntax for type uuid")}
	s := New(rr, pr, &fakeEnforcer{}, &fakePortalCfgRepo{}, &fakeTokenCfgRepo{}, nil, zap.NewNop())

	_, err := s.SeedPortal(ctx, authzentityv1.Portal_PORTAL_AGENT, GlobalTenantID, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create policy")
}
