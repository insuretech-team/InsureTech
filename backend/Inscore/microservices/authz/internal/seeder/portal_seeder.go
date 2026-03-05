package seeder

// portal_seeder.go — seeds default roles, policy rules, and Casbin g/p rules
// for every portal on first startup or when explicitly triggered via SeedPortalDefaults RPC.
//
// Portal → Default Roles mapping (aligned with authz_implementation.md):
//
//  PORTAL_SYSTEM:    super_admin | admin | support | auditor | readonly
//  PORTAL_BUSINESS:  owner | admin | finance | hr | readonly
//  PORTAL_B2B:       partner_admin | partner_user | api_client
//  PORTAL_AGENT:     senior_agent | agent | agent_trainee
//  PORTAL_REGULATOR: regulator_admin | inspector | auditor
//  PORTAL_CUSTOMER:  customer
//
// Domain key format: "PORTAL_SYSTEM:global"
// Subject format:    "user:<uuid>"  →  inherits  "role:super_admin" in domain
// Object format:     "svc:<service>/<resource>"  e.g. "svc:policy/*"
// Action format:     HTTP verb or "*"

import (
	"context"
	"errors"
	"os"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/domain"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
)

// GlobalTenantID is the tenant suffix used by gateway/runtime when no tenant is provided.
const GlobalTenantID = "root"

// portalRole describes a role and its seeded policy rules.
type portalRole struct {
	Name        string
	DisplayName string
	IsSystem    bool
	Policies    []seedPolicy
}

// seedPolicy is a single (sub, dom, obj, act, eft) Casbin p-rule.
type seedPolicy struct {
	Object string // "svc:policy/*"
	Action string // "POST" | "GET" | "*"
	Effect string // "allow" | "deny"
}

// portalSeedMap defines the default roles + policies seeded for each portal.
var portalSeedMap = map[authzentityv1.Portal][]portalRole{

	// ── PORTAL_SYSTEM ─────────────────────────────────────────────────────────
	authzentityv1.Portal_PORTAL_SYSTEM: {
		{Name: "super_admin", DisplayName: "Super Administrator", IsSystem: true, Policies: []seedPolicy{
			{Object: "svc:*", Action: "*", Effect: "allow"},
		}},
		{Name: "admin", DisplayName: "System Administrator", IsSystem: true, Policies: []seedPolicy{
			{Object: "svc:user/*", Action: "*", Effect: "allow"},
			{Object: "svc:role/*", Action: "*", Effect: "allow"},
			{Object: "svc:policy/*", Action: "*", Effect: "allow"},
			{Object: "svc:storage/*", Action: "*", Effect: "allow"},
			{Object: "svc:audit/*", Action: "GET", Effect: "allow"},
		}},
		{Name: "support", DisplayName: "Support Agent", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:user/*", Action: "GET", Effect: "allow"},
			{Object: "svc:claim/*", Action: "GET", Effect: "allow"},
			{Object: "svc:session/*", Action: "DELETE", Effect: "allow"},
			{Object: "svc:b2b/*", Action: "GET", Effect: "allow"},
			{Object: "svc:b2b/purchase-orders", Action: "POST", Effect: "allow"},
		}},
		{Name: "auditor", DisplayName: "System Auditor", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:audit/*", Action: "GET", Effect: "allow"},
			{Object: "svc:*", Action: "GET", Effect: "allow"},
		}},
		{Name: "readonly", DisplayName: "Read Only", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:*", Action: "GET", Effect: "allow"},
		}},
	},

	// ── PORTAL_BUSINESS ───────────────────────────────────────────────────────
	authzentityv1.Portal_PORTAL_BUSINESS: {
		{Name: "owner", DisplayName: "Business Owner", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:*", Action: "*", Effect: "allow"},
		}},
		{Name: "admin", DisplayName: "Business Administrator", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:user/*", Action: "*", Effect: "allow"},
			{Object: "svc:policy/*", Action: "*", Effect: "allow"},
			{Object: "svc:claim/*", Action: "*", Effect: "allow"},
			{Object: "svc:document/*", Action: "*", Effect: "allow"},
			{Object: "svc:storage/*", Action: "*", Effect: "allow"},
		}},
		{Name: "finance", DisplayName: "Finance Manager", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:invoice/*", Action: "*", Effect: "allow"},
			{Object: "svc:payment/*", Action: "*", Effect: "allow"},
			{Object: "svc:report/*", Action: "GET", Effect: "allow"},
		}},
		{Name: "hr", DisplayName: "HR Manager", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:employee/*", Action: "*", Effect: "allow"},
			{Object: "svc:enrollment/*", Action: "*", Effect: "allow"},
		}},
		{Name: "readonly", DisplayName: "Read Only", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:*", Action: "GET", Effect: "allow"},
		}},
	},

	// ── PORTAL_B2B ────────────────────────────────────────────────────────────
	authzentityv1.Portal_PORTAL_B2B: {
		{Name: "partner_admin", DisplayName: "Partner Administrator", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:b2b/*", Action: "*", Effect: "allow"},
			{Object: "svc:partner/*", Action: "*", Effect: "allow"},
			{Object: "svc:policy/*", Action: "*", Effect: "allow"},
			{Object: "svc:user/*", Action: "*", Effect: "allow"},
		}},
		{Name: "partner_user", DisplayName: "Partner User", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:b2b/*", Action: "GET", Effect: "allow"},
			{Object: "svc:b2b/purchase-orders", Action: "POST", Effect: "allow"},
			{Object: "svc:policy/*", Action: "GET", Effect: "allow"},
			{Object: "svc:claim/*", Action: "POST", Effect: "allow"},
			{Object: "svc:claim/*", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/upload", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/upload-batch", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/upload-url", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/finalize", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/get", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/download-url", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/download-url", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/update", Action: "PATCH", Effect: "allow"},
			{Object: "svc:storage/delete", Action: "DELETE", Effect: "allow"},
		}},
		{Name: "api_client", DisplayName: "API Client (Machine)", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:policy/quote", Action: "POST", Effect: "allow"},
			{Object: "svc:policy/bind", Action: "POST", Effect: "allow"},
			{Object: "svc:claim/submit", Action: "POST", Effect: "allow"},
		}},
	},

	// ── PORTAL_AGENT ──────────────────────────────────────────────────────────
	authzentityv1.Portal_PORTAL_AGENT: {
		{Name: "senior_agent", DisplayName: "Senior Agent", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:policy/*", Action: "*", Effect: "allow"},
			{Object: "svc:claim/*", Action: "*", Effect: "allow"},
			{Object: "svc:customer/*", Action: "*", Effect: "allow"},
			{Object: "svc:document/*", Action: "*", Effect: "allow"},
			{Object: "svc:commission/*", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/*", Action: "*", Effect: "allow"},
		}},
		{Name: "agent", DisplayName: "Insurance Agent", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:policy/*", Action: "GET", Effect: "allow"},
			{Object: "svc:policy/quote", Action: "POST", Effect: "allow"},
			{Object: "svc:claim/*", Action: "GET", Effect: "allow"},
			{Object: "svc:customer/*", Action: "GET", Effect: "allow"},
			{Object: "svc:document/*", Action: "*", Effect: "allow"},
			{Object: "svc:storage/upload", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/upload-batch", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/upload-url", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/finalize", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/get", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/download-url", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/download-url", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/update", Action: "PATCH", Effect: "allow"},
			{Object: "svc:storage/delete", Action: "DELETE", Effect: "allow"},
		}},
		{Name: "agent_trainee", DisplayName: "Agent Trainee", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:policy/*", Action: "GET", Effect: "allow"},
			{Object: "svc:customer/*", Action: "GET", Effect: "allow"},
		}},
	},

	// ── PORTAL_REGULATOR ──────────────────────────────────────────────────────
	authzentityv1.Portal_PORTAL_REGULATOR: {
		{Name: "regulator_admin", DisplayName: "Regulatory Administrator", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:audit/*", Action: "*", Effect: "allow"},
			{Object: "svc:report/*", Action: "*", Effect: "allow"},
			{Object: "svc:compliance/*", Action: "*", Effect: "allow"},
		}},
		{Name: "inspector", DisplayName: "Regulatory Inspector", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:audit/*", Action: "GET", Effect: "allow"},
			{Object: "svc:report/*", Action: "GET", Effect: "allow"},
			{Object: "svc:compliance/*", Action: "GET", Effect: "allow"},
		}},
		{Name: "auditor", DisplayName: "Regulatory Auditor", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:audit/*", Action: "GET", Effect: "allow"},
			{Object: "svc:report/*", Action: "GET", Effect: "allow"},
		}},
	},

	// ── PORTAL_CUSTOMER ───────────────────────────────────────────────────────
	authzentityv1.Portal_PORTAL_B2C: {
		{Name: "customer", DisplayName: "Customer", IsSystem: false, Policies: []seedPolicy{
			{Object: "svc:policy/my/*", Action: "GET", Effect: "allow"},
			{Object: "svc:claim/my/*", Action: "*", Effect: "allow"},
			{Object: "svc:profile/*", Action: "*", Effect: "allow"},
			{Object: "svc:document/my/*", Action: "*", Effect: "allow"},
			{Object: "svc:document/*", Action: "*", Effect: "allow"},
			{Object: "svc:storage/upload", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/upload-batch", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/upload-url", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/finalize", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/get", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/download-url", Action: "GET", Effect: "allow"},
			{Object: "svc:storage/download-url", Action: "POST", Effect: "allow"},
			{Object: "svc:storage/update", Action: "PATCH", Effect: "allow"},
			{Object: "svc:storage/delete", Action: "DELETE", Effect: "allow"},
		}},
	},
}

// ── Seeder ────────────────────────────────────────────────────────────────────

// PortalSeeder implements domain.PortalSeeder.
type PortalSeeder struct {
	roleRepo         domain.RoleRepository
	policyRepo       domain.PolicyRuleRepository
	enforcer         domain.EnforcerIface
	portalConfigRepo domain.PortalConfigRepository
	tokenConfigRepo  domain.TokenConfigRepository
	db               *gorm.DB
	logger           *zap.Logger
}

func New(
	roleRepo domain.RoleRepository,
	policyRepo domain.PolicyRuleRepository,
	enforcer domain.EnforcerIface,
	portalConfigRepo domain.PortalConfigRepository,
	tokenConfigRepo domain.TokenConfigRepository,
	db *gorm.DB,
	logger *zap.Logger,
) *PortalSeeder {
	if logger == nil {
		// appLogger has AddCallerSkip(+1) for wrapper helpers; compensate when
		// using the underlying zap logger directly inside this layer.
		logger = appLogger.GetLogger().WithOptions(zap.AddCallerSkip(-1))
	}
	return &PortalSeeder{
		roleRepo:         roleRepo,
		policyRepo:       policyRepo,
		enforcer:         enforcer,
		portalConfigRepo: portalConfigRepo,
		tokenConfigRepo:  tokenConfigRepo,
		db:               db,
		logger:           logger,
	}
}

// SeedAllPortals seeds every portal with GlobalTenantID.
func (s *PortalSeeder) SeedAllPortals(ctx context.Context) error {
	portals := []authzentityv1.Portal{
		authzentityv1.Portal_PORTAL_SYSTEM,
		authzentityv1.Portal_PORTAL_BUSINESS,
		authzentityv1.Portal_PORTAL_B2B,
		authzentityv1.Portal_PORTAL_AGENT,
		authzentityv1.Portal_PORTAL_REGULATOR,
		authzentityv1.Portal_PORTAL_B2C,
	}
	var errs []error
	for _, portal := range portals {
		if _, err := s.SeedPortal(ctx, portal, GlobalTenantID, false); err != nil {
			errs = append(errs, errors.New("seed "+portal.String()+": "+err.Error()))
			s.logger.Error("portal seed failed", zap.String("portal", portal.String()), zap.Error(err))
		} else {
			s.logger.Info("portal seeded", zap.String("portal", portal.String()))
		}
	}
	if err := s.SeedPortalConfigs(ctx); err != nil {
		s.logger.Warn("portal_configs seeding had errors (non-fatal)", zap.Error(err))
	}
	if err := s.SeedTokenConfig(ctx); err != nil {
		s.logger.Warn("token_config seeding had errors (non-fatal)", zap.Error(err))
	}
	if err := s.SeedDefaultSystemRoleBindings(ctx); err != nil {
		s.logger.Warn("system role binding seeding had errors (non-fatal)", zap.Error(err))
	}
	if len(errs) > 0 {
		return errors.New("seed errors occurred (check logs for details)")
	}
	return nil
}

type authnSystemUserRow struct {
	UserID string `gorm:"column:user_id"`
}

// SeedDefaultSystemRoleBindings ensures existing SYSTEM_USER accounts have at
// least a support role in the system:root domain.
//
// Why this exists:
//   - authz already seeds roles/policies, but pre-existing users (for example
//     ADMIN user seeded by authn) may not have a g-rule assignment yet.
//   - this bootstrap is idempotent and safe on every authz startup.
func (s *PortalSeeder) SeedDefaultSystemRoleBindings(ctx context.Context) error {
	if s.db == nil || s.enforcer == nil {
		return nil
	}

	var users []authnSystemUserRow
	if err := s.db.WithContext(ctx).
		Table("authn_schema.users").
		Select("user_id").
		Where("deleted_at IS NULL").
		Where("user_type IN ?", []string{"USER_TYPE_SYSTEM_USER", "SYSTEM_USER", "4"}).
		Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	assigned := 0
	for _, user := range users {
		userID := strings.TrimSpace(user.UserID)
		if userID == "" {
			continue
		}
		err := s.enforcer.AddRoleForUserInDomain("user:"+userID, "role:support", "system:root")
		if err != nil {
			if isLikelyAlreadyExistsErr(err) {
				continue
			}
			s.logger.Warn("system role binding failed",
				zap.String("user_id", userID),
				zap.Error(err),
			)
			continue
		}
		assigned++
	}

	s.logger.Info("system role binding seed complete",
		zap.Int("system_users", len(users)),
		zap.Int("assigned", assigned),
	)
	return nil
}

// SeedPortal seeds a single portal+tenant with default roles and policies.
// If dryRun=true it returns counts without writing anything.
func (s *PortalSeeder) SeedPortal(ctx context.Context, portal authzentityv1.Portal, tenantID string, dryRun bool) (*domain.SeedResult, error) {
	roles, ok := portalSeedMap[portal]
	if !ok {
		return &domain.SeedResult{}, nil
	}

	domainKey := portalDomainKey(portal, tenantID)
	result := &domain.SeedResult{}
	existingPolicyKeys := map[string]struct{}{}
	if !dryRun {
		var err error
		existingPolicyKeys, err = s.loadExistingPolicyKeys(ctx, domainKey)
		if err != nil {
			return nil, errors.New("load existing policy keys: " + err.Error())
		}
	}

	for _, rd := range roles {
		// 1. Create or skip role
		existingRole, err := s.roleRepo.GetByNameAndPortal(ctx, rd.Name, portal)
		if err != nil {
			if !isNotFoundErr(err) {
				return nil, errors.New("lookup role " + rd.Name + ": " + err.Error())
			}
			// Role doesn't exist — create it
			if !dryRun {
				role := &authzentityv1.Role{
					Name:        rd.Name,
					Portal:      portal,
					Description: rd.DisplayName,
					IsSystem:    rd.IsSystem,
					IsActive:    true,
					CreatedBy:   "",
				}
				createdRole, err := s.roleRepo.Create(ctx, role)
				if err != nil {
					return nil, errors.New("create role " + rd.Name + ": " + err.Error())
				}
				existingRole = createdRole
			}
			result.RolesSeeded++
		}

		if existingRole == nil {
			continue
		}

		roleSubject := "role:" + rd.Name

		// 2. Seed policy rules + Casbin p-rules
		for _, sp := range rd.Policies {
			policyKey := buildPolicySeedKey(roleSubject, domainKey, sp.Object, sp.Action, sp.Effect)
			alreadyPresent := false
			if _, ok := existingPolicyKeys[policyKey]; ok {
				alreadyPresent = true
			}

			if !dryRun {
				if !alreadyPresent {
					// Persist to policy_rules table
					pr := &authzentityv1.PolicyRule{
						Subject:     roleSubject,
						Domain:      domainKey,
						Object:      sp.Object,
						Action:      sp.Action,
						Effect:      effectFromString(sp.Effect),
						Description: "seeded default for " + portal.String() + "/" + rd.Name,
						IsActive:    true,
						CreatedBy:   "",
					}
					if _, err := s.policyRepo.Create(ctx, pr); err != nil {
						if isLikelyAlreadyExistsErr(err) {
							s.logger.Debug("policy already exists (skip)", zap.String("role", rd.Name), zap.String("obj", sp.Object), zap.String("action", sp.Action))
							existingPolicyKeys[policyKey] = struct{}{}
							alreadyPresent = true
						} else {
							return nil, errors.New("create policy for role " + rd.Name + " (" + sp.Object + " " + sp.Action + "): " + err.Error())
						}
					} else {
						existingPolicyKeys[policyKey] = struct{}{}
					}
				}

				// Add Casbin p-rule: (role:name, domain, obj, act, eft)
				if s.enforcer != nil && !alreadyPresent {
					if addErr := s.enforcer.AddPolicy(roleSubject, domainKey, sp.Object, sp.Action, sp.Effect); addErr != nil {
						if isLikelyAlreadyExistsErr(addErr) {
							s.logger.Debug("casbin policy already exists (skip)", zap.String("role", rd.Name), zap.String("obj", sp.Object), zap.String("action", sp.Action))
						} else {
							return nil, errors.New("add casbin policy for role " + rd.Name + " (" + sp.Object + " " + sp.Action + "): " + addErr.Error())
						}
					}
				}
				if alreadyPresent {
					continue
				}
			}
			result.PoliciesSeeded++
			result.RulesSeeded++
		}
	}

	return result, nil
}

// SeedPortalConfigs upserts a default PortalConfig row for each portal.
// Safe to run multiple times (idempotent via UPSERT ON CONFLICT portal).
func (s *PortalSeeder) SeedPortalConfigs(ctx context.Context) error {
	portals := []authzentityv1.Portal{
		authzentityv1.Portal_PORTAL_SYSTEM,
		authzentityv1.Portal_PORTAL_BUSINESS,
		authzentityv1.Portal_PORTAL_B2B,
		authzentityv1.Portal_PORTAL_AGENT,
		authzentityv1.Portal_PORTAL_REGULATOR,
		authzentityv1.Portal_PORTAL_B2C,
	}
	for _, portal := range portals {
		mfaRequired := false
		mfaMethods := []string{}
		switch portal {
		case authzentityv1.Portal_PORTAL_SYSTEM:
			mfaRequired = true
			mfaMethods = []string{"totp"}
		case authzentityv1.Portal_PORTAL_BUSINESS:
			mfaRequired = true
			mfaMethods = []string{"email_otp"}
		case authzentityv1.Portal_PORTAL_B2B:
			mfaRequired = true
			mfaMethods = []string{"totp", "email_otp"}
		case authzentityv1.Portal_PORTAL_AGENT:
			mfaRequired = false
			mfaMethods = []string{"sms_otp"}
		case authzentityv1.Portal_PORTAL_REGULATOR:
			mfaRequired = true
			mfaMethods = []string{"totp"}
		case authzentityv1.Portal_PORTAL_B2C:
			mfaRequired = false
			mfaMethods = []string{"sms_otp"}
		}
		cfg := &authzentityv1.PortalConfig{
			Portal:                  portal,
			MfaRequired:             mfaRequired,
			MfaMethods:              mfaMethods,
			AccessTokenTtlSeconds:   900,    // 15 min
			RefreshTokenTtlSeconds:  604800, // 7 days
			SessionTtlSeconds:       28800,  // 8 hours
			IdleTimeoutSeconds:      1800,   // 30 min
			AllowConcurrentSessions: true,
			MaxConcurrentSessions:   5,
			UpdatedBy:               "seeder",
		}
		if _, err := s.portalConfigRepo.Upsert(ctx, cfg); err != nil {
			s.logger.Error("portal_config seed failed", zap.String("portal", portal.String()), zap.Error(err))
		} else {
			s.logger.Info("portal_config seeded", zap.String("portal", portal.String()))
		}
	}
	return nil
}

// SeedTokenConfig seeds the active JWT signing key into authz_schema.token_configs.
// Reads public key from JWT_PUBLIC_KEY_PEM env var and kid from JWT_KEY_ID (default: insuretech-2025-01).
// Safe to run multiple times (skips if kid already exists).
func (s *PortalSeeder) SeedTokenConfig(ctx context.Context) error {
	kid := os.Getenv("JWT_KEY_ID")
	if kid == "" {
		kid = "insuretech-2025-01"
	}
	publicKeyPEM := os.Getenv("JWT_PUBLIC_KEY_PEM")
	if publicKeyPEM == "" {
		s.logger.Warn("JWT_PUBLIC_KEY_PEM not set — skipping TokenConfig seed")
		return nil
	}
	privateKeyRef := os.Getenv("JWT_PRIVATE_KEY_REF")
	if privateKeyRef == "" {
		privateKeyRef = "secret/authz/jwt-signing-key"
	}

	// Check if already exists
	existing, _ := s.tokenConfigRepo.GetActive(ctx)
	if existing != nil && existing.Kid == kid {
		s.logger.Info("token_config already seeded", zap.String("kid", kid))
		return nil
	}

	cfg := &authzentityv1.TokenConfig{
		Kid:           kid,
		Algorithm:     "RS256",
		PublicKeyPem:  publicKeyPEM,
		PrivateKeyRef: privateKeyRef,
		IsActive:      true,
		CreatedAt:     timestamppb.Now(),
	}
	if _, err := s.tokenConfigRepo.Create(ctx, cfg); err != nil {
		s.logger.Error("token_config seed failed", zap.String("kid", kid), zap.Error(err))
		return err
	}
	s.logger.Info("token_config seeded", zap.String("kid", kid))
	return nil
}

func effectFromString(s string) authzentityv1.PolicyEffect {
	if s == "deny" {
		return authzentityv1.PolicyEffect_POLICY_EFFECT_DENY
	}
	return authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW
}

func normalizeSeedEffect(effect string) string {
	if strings.EqualFold(strings.TrimSpace(effect), "deny") {
		return "deny"
	}
	return "allow"
}

func buildPolicySeedKey(subject, domain, object, action, effect string) string {
	return strings.ToLower(strings.TrimSpace(subject)) + "|" +
		strings.ToLower(strings.TrimSpace(domain)) + "|" +
		strings.ToLower(strings.TrimSpace(object)) + "|" +
		strings.ToLower(strings.TrimSpace(action)) + "|" +
		normalizeSeedEffect(effect)
}

func (s *PortalSeeder) loadExistingPolicyKeys(ctx context.Context, domainKey string) (map[string]struct{}, error) {
	keys := make(map[string]struct{})
	const pageSize = 500
	offset := 0
	for {
		rows, err := s.policyRepo.List(ctx, domainKey, false, pageSize, offset)
		if err != nil {
			return nil, err
		}
		for _, pr := range rows {
			if pr == nil {
				continue
			}
			effect := "allow"
			if pr.Effect == authzentityv1.PolicyEffect_POLICY_EFFECT_DENY {
				effect = "deny"
			}
			keys[buildPolicySeedKey(pr.Subject, pr.Domain, pr.Object, pr.Action, effect)] = struct{}{}
		}
		if len(rows) < pageSize {
			break
		}
		offset += pageSize
	}
	return keys, nil
}

func isLikelyAlreadyExistsErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "violates unique")
}

func isNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "no rows")
}

func portalDomainKey(portal authzentityv1.Portal, tenantID string) string {
	prefix := strings.ToLower(strings.TrimPrefix(portal.String(), "PORTAL_"))
	if prefix == "unspecified" {
		prefix = "b2c"
	}
	if tenantID == "" {
		tenantID = GlobalTenantID
	}
	return prefix + ":" + tenantID
}
