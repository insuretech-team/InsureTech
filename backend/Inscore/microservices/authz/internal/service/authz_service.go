package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/cache"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/domain"
	authzevents "github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/metrics"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthZService implements all business logic for the AuthZ microservice.
// It orchestrates: Casbin enforcer + repositories + audit + cache.
type AuthZService struct {
	enforcer          domain.EnforcerIface
	roleRepo          domain.RoleRepository
	userRoleRepo      domain.UserRoleRepository
	policyRepo        domain.PolicyRuleRepository
	portalRepo        domain.PortalConfigRepository
	auditRepo         domain.AuditRepository
	auditAllDecisions bool
	publisher         *authzevents.Publisher
	permCache         *cache.PermissionCache // Permission cache (optional)
}

func New(
	enforcer domain.EnforcerIface,
	roleRepo domain.RoleRepository,
	userRoleRepo domain.UserRoleRepository,
	policyRepo domain.PolicyRuleRepository,
	portalRepo domain.PortalConfigRepository,
	auditRepo domain.AuditRepository,
	auditAllDecisions bool,
	publisher *authzevents.Publisher,
) *AuthZService {
	return &AuthZService{
		enforcer:          enforcer,
		roleRepo:          roleRepo,
		userRoleRepo:      userRoleRepo,
		policyRepo:        policyRepo,
		portalRepo:        portalRepo,
		auditRepo:         auditRepo,
		auditAllDecisions: auditAllDecisions,
		publisher:         publisher,
		permCache:         nil, // Initialized separately via SetPermissionCache
	}
}

// SetPermissionCache sets the permission cache (optional).
// Call after creating the service to enable caching.
func (s *AuthZService) SetPermissionCache(permCache *cache.PermissionCache) {
	s.permCache = permCache
}

// ── CheckAccess — core enforcement ──────────────────────────────────────────

// CheckAccess performs a Casbin PERM enforce call and audits the decision.
// Deny-by-default: no matching rule → DENY.
// subject = "user:<user_id>", domain = "portal:tenant_id"
// 
// API Key Scope Integration:
// If the request context contains API key scopes (via attributes["api_key_scopes"]),
// the scopes are validated BEFORE the Casbin check. This ensures API keys are restricted
// to their defined scopes even if the user has broader permissions.
func (s *AuthZService) CheckAccess(ctx context.Context, req *authzservicev1.CheckAccessRequest) (*authzservicev1.CheckAccessResponse, error) {
	subject := "user:" + req.UserId
	start := time.Now()
	
	// 🔑 API Key Scope Validation (if present)
	// This check happens BEFORE Casbin to ensure API keys are properly restricted
	if req.Context != nil && req.Context.Attributes != nil {
		validator := NewApiKeyScopeValidator()
		scopes := validator.ParseScopesFromAttributes(req.Context.Attributes)
		
		if len(scopes) > 0 {
			// API key is being used - validate scopes first
			scopeStart := time.Now()
			allowed, reason := validator.ValidateScope(scopes, req.Object, req.Action)
			scopeLatencyMs := float64(time.Since(scopeStart).Milliseconds())
			
			// Record scope validation metrics
			metrics.RecordAPIScopeValidation(req.Domain, allowed, scopeLatencyMs)
			
			if !allowed {
				// Record the denial reason
				metrics.RecordAPIScopeDenial(reason)
				
				latencyMs := float64(time.Since(start).Milliseconds())
				metrics.RecordDecision(req.Domain, false, latencyMs)
				
				// Audit the denial
				if s.auditAllDecisions {
					ipAddr, userAgent, sessionID, _, _ := extractContext(req.Context)
					_ = s.auditRepo.Create(ctx, &authzentityv1.AccessDecisionAudit{
						AuditId:     uuid.New().String(),
						UserId:      req.UserId,
						SessionId:   sessionID,
						Domain:      req.Domain,
						Subject:     subject,
						Object:      req.Object,
						Action:      req.Action,
						Decision:    authzentityv1.PolicyEffect_POLICY_EFFECT_DENY,
						MatchedRule: "api_key_scope_check",
						IpAddress:   ipAddr,
						UserAgent:   userAgent,
						DecidedAt:   timestamppb.Now(),
					})
				}
				
				return &authzservicev1.CheckAccessResponse{
					Allowed:     false,
					Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_DENY,
					Reason:      reason,
					MatchedRule: "api_key_scope_check",
				}, nil
			}
			// Scope check passed - continue to Casbin check
		}
	}
	
	// 🚀 Check cache first
	if s.permCache != nil {
		if cached, found := s.permCache.Get(ctx, req.UserId, req.Domain, req.Object, req.Action); found {
			latencyMs := float64(time.Since(start).Milliseconds())
			metrics.RecordDecision(req.Domain, cached, latencyMs)
			metrics.RecordCacheHit(true)
			
			effect := authzentityv1.PolicyEffect_POLICY_EFFECT_DENY
			if cached {
				effect = authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW
			}
			
			return &authzservicev1.CheckAccessResponse{
				Allowed:     cached,
				Effect:      effect,
				Reason:      "",
				MatchedRule: "(cached)",
			}, nil
		}
		metrics.RecordCacheHit(false)
	}
	
	// Cache miss - query Casbin
	allowed, matchedRule, err := s.enforcer.Enforce(ctx, subject, req.Domain, req.Object, req.Action)
	latencyMs := float64(time.Since(start).Milliseconds())
	if err != nil {
		return nil, errors.New("enforce error: " + err.Error())
	}

	effect := authzentityv1.PolicyEffect_POLICY_EFFECT_DENY
	reason := "no matching policy (deny-by-default)"
	if allowed {
		effect = authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW
		reason = ""
	}
	
	// Store in cache
	if s.permCache != nil {
		_ = s.permCache.Set(ctx, req.UserId, req.Domain, req.Object, req.Action, allowed)
	}

	// Record decision metrics
	metrics.RecordDecision(req.Domain, allowed, latencyMs)

	// 🔍 Audit ────────────────────────────────────────────────────────────────
	if !allowed || s.auditAllDecisions {
		ipAddr, userAgent, sessionID, _, _ := extractContext(req.Context)
		_ = s.auditRepo.Create(ctx, &authzentityv1.AccessDecisionAudit{
			AuditId:     uuid.New().String(),
			UserId:      req.UserId,
			SessionId:   sessionID,
			Domain:      req.Domain,
			Subject:     subject,
			Object:      req.Object,
			Action:      req.Action,
			Decision:    effect,
			MatchedRule: matchedRule,
			IpAddress:   ipAddr,
			UserAgent:   userAgent,
			DecidedAt:   timestamppb.Now(),
		})
		// 📡 Publish DENY event → SIEM
		if !allowed && s.publisher != nil {
			_ = s.publisher.PublishAccessDenied(ctx, req.UserId, req.Domain, req.Object, req.Action, sessionID, ipAddr)
		}
	}

	return &authzservicev1.CheckAccessResponse{
		Allowed:     allowed,
		Effect:      effect,
		MatchedRule: matchedRule,
		Reason:      reason,
	}, nil
}

// BatchCheckAccess checks multiple (obj, act) tuples for the same user+domain.
func (s *AuthZService) BatchCheckAccess(ctx context.Context, req *authzservicev1.BatchCheckAccessRequest) (*authzservicev1.BatchCheckAccessResponse, error) {
	subject := "user:" + req.UserId
	results := make([]*authzservicev1.AccessCheckResult, 0, len(req.Checks))
	for _, check := range req.Checks {
		allowed, _, err := s.enforcer.Enforce(ctx, subject, req.Domain, check.Object, check.Action)
		reason := ""
		if err != nil {
			reason = err.Error()
			allowed = false
		} else if !allowed {
			reason = "no matching policy"
		}
		results = append(results, &authzservicev1.AccessCheckResult{
			Object:  check.Object,
			Action:  check.Action,
			Allowed: allowed,
			Reason:  reason,
		})
	}
	return &authzservicev1.BatchCheckAccessResponse{Results: results}, nil
}

// ── Role Management ──────────────────────────────────────────────────────────

func (s *AuthZService) CreateRole(ctx context.Context, req *authzservicev1.CreateRoleRequest) (*authzservicev1.CreateRoleResponse, error) {
	role := &authzentityv1.Role{
		RoleId:      uuid.New().String(),
		Name:        req.Name,
		Portal:      req.Portal,
		Description: req.Description,
		IsSystem:    false,
		IsActive:    true,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}
	if _, err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, errors.New("create role: " + err.Error())
	}
	if s.publisher != nil {
		_ = s.publisher.PublishRoleCreated(ctx, role)
	}
	return &authzservicev1.CreateRoleResponse{Role: role}, nil
}

func (s *AuthZService) GetRole(ctx context.Context, req *authzservicev1.GetRoleRequest) (*authzservicev1.GetRoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, req.RoleId)
	if err != nil {
		return nil, errors.New("get role: " + err.Error())
	}
	return &authzservicev1.GetRoleResponse{Role: role}, nil
}

func (s *AuthZService) UpdateRole(ctx context.Context, req *authzservicev1.UpdateRoleRequest) (*authzservicev1.UpdateRoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, req.RoleId)
	if err != nil {
		return nil, errors.New("role not found: " + err.Error())
	}
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	role.UpdatedAt = timestamppb.Now()
	if _, err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, errors.New("update role: " + err.Error())
	}
	return &authzservicev1.UpdateRoleResponse{Role: role}, nil
}

func (s *AuthZService) DeleteRole(ctx context.Context, req *authzservicev1.DeleteRoleRequest) (*authzservicev1.DeleteRoleResponse, error) {
	if err := s.roleRepo.SoftDelete(ctx, req.RoleId); err != nil {
		return nil, errors.New("delete role: " + err.Error())
	}
	return &authzservicev1.DeleteRoleResponse{Message: "role deleted"}, nil
}

func (s *AuthZService) ListRoles(ctx context.Context, req *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 50
	}
	roles, err := s.roleRepo.List(ctx, req.Portal, req.ActiveOnly, pageSize, 0 /* page_token cursor not yet impl */)
	if err != nil {
		return nil, errors.New("list roles: " + err.Error())
	}
	return &authzservicev1.ListRolesResponse{Roles: roles}, nil
}

// ── User-Role Assignment ─────────────────────────────────────────────────────

func (s *AuthZService) AssignRole(ctx context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, req.RoleId)
	if err != nil {
		return nil, errors.New("role not found: " + err.Error())
	}

	domain := req.Domain // "portal:tenant_id"
	subject := "user:" + req.UserId
	roleName := "role:" + role.Name

	// 1. Add g-type rule to Casbin enforcer (in-memory + DB via gorm-adapter)
	if err := s.enforcer.AddRoleForUserInDomain(subject, roleName, domain); err != nil {
		return nil, errors.New("casbin assign role: " + err.Error())
	}

	// 2. Persist user_roles row for audit/query
	var expiresAt *timestamppb.Timestamp
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt
	}
	ur := &authzentityv1.UserRole{
		UserRoleId: uuid.New().String(),
		UserId:     req.UserId,
		RoleId:     req.RoleId,
		Domain:     domain,
		AssignedBy: req.AssignedBy,
		AssignedAt: timestamppb.Now(),
		ExpiresAt:  expiresAt,
	}
	if _, err := s.userRoleRepo.Assign(ctx, ur); err != nil {
		return nil, errors.New("persist user_role: " + err.Error())
	}

	if s.publisher != nil {
		_ = s.publisher.PublishRoleAssigned(ctx, req.UserId, req.RoleId, role.Name, domain, req.AssignedBy)
	}
	return &authzservicev1.AssignRoleResponse{
		UserRole: ur,
	}, nil
}

func (s *AuthZService) RemoveRole(ctx context.Context, req *authzservicev1.RemoveRoleRequest) (*authzservicev1.RemoveRoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, req.RoleId)
	if err != nil {
		return nil, errors.New("role not found: " + err.Error())
	}
	subject := "user:" + req.UserId
	roleName := "role:" + role.Name

	if err := s.enforcer.DeleteRoleForUserInDomain(subject, roleName, req.Domain); err != nil {
		return nil, errors.New("casbin remove role: " + err.Error())
	}
	if err := s.userRoleRepo.Remove(ctx, req.UserId, req.RoleId, req.Domain); err != nil {
		return nil, errors.New("remove user_role row: " + err.Error())
	}
	return &authzservicev1.RemoveRoleResponse{Message: "role removed"}, nil
}

func (s *AuthZService) ListUserRoles(ctx context.Context, req *authzservicev1.ListUserRolesRequest) (*authzservicev1.ListUserRolesResponse, error) {
	urs, err := s.userRoleRepo.ListByUser(ctx, req.UserId, req.Domain)
	if err != nil {
		return nil, errors.New("list user roles: " + err.Error())
	}
	return &authzservicev1.ListUserRolesResponse{UserRoles: urs}, nil
}

func (s *AuthZService) GetUserPermissions(ctx context.Context, req *authzservicev1.GetUserPermissionsRequest) (*authzservicev1.GetUserPermissionsResponse, error) {
	subject := "user:" + req.UserId
	perms, err := s.enforcer.GetPermissionsForUserInDomain(subject, req.Domain)
	if err != nil {
		return nil, errors.New("get permissions: " + err.Error())
	}
	// Also resolve role-inherited permissions via g-rules
	roles, _ := s.enforcer.GetRolesForUserInDomain(subject, req.Domain)
	effective := make([]*authzservicev1.EffectivePermission, 0, len(perms))
	// Direct user permissions
	for _, p := range perms {
		if len(p) >= 4 {
			effective = append(effective, &authzservicev1.EffectivePermission{
				Object:  p[2],
				Action:  p[3],
				ViaRole: "",
			})
		}
	}
	// Role-inherited permissions
	for _, role := range roles {
		rolePerms, _ := s.enforcer.GetPermissionsForUserInDomain(role, req.Domain)
		for _, p := range rolePerms {
			if len(p) >= 4 {
				effective = append(effective, &authzservicev1.EffectivePermission{
					Object:  p[2],
					Action:  p[3],
					ViaRole: role,
				})
			}
		}
	}
	return &authzservicev1.GetUserPermissionsResponse{Permissions: effective}, nil
}

// ── Policy Rule Management ────────────────────────────────────────────────────

func (s *AuthZService) CreatePolicyRule(ctx context.Context, req *authzservicev1.CreatePolicyRuleRequest) (*authzservicev1.CreatePolicyRuleResponse, error) {
	effect := authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW
	if req.Effect == authzentityv1.PolicyEffect_POLICY_EFFECT_DENY {
		effect = authzentityv1.PolicyEffect_POLICY_EFFECT_DENY
	}
	effectStr := "allow"
	if effect == authzentityv1.PolicyEffect_POLICY_EFFECT_DENY {
		effectStr = "deny"
	}

	// 1. Add p-type rule to Casbin
	if err := s.enforcer.AddPolicy(req.Subject, req.Domain, req.Object, req.Action, effectStr); err != nil {
		return nil, errors.New("casbin add policy: " + err.Error())
	}

	// 2. Persist human-readable policy_rules row
	pr := &authzentityv1.PolicyRule{
		PolicyId:    uuid.New().String(),
		Subject:     req.Subject,
		Domain:      req.Domain,
		Object:      req.Object,
		Action:      req.Action,
		Effect:      effect,
		Condition:   req.Condition,
		Description: req.Description,
		IsActive:    true,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}
	if _, err := s.policyRepo.Create(ctx, pr); err != nil {
		return nil, errors.New("persist policy_rule: " + err.Error())
	}
	if s.publisher != nil {
		_ = s.publisher.PublishPolicyRuleCreated(ctx, pr)
	}
	return &authzservicev1.CreatePolicyRuleResponse{Policy: pr}, nil
}

func (s *AuthZService) UpdatePolicyRule(ctx context.Context, req *authzservicev1.UpdatePolicyRuleRequest) (*authzservicev1.UpdatePolicyRuleResponse, error) {
	// Fetch existing rule
	rules, err := s.policyRepo.List(ctx, "", true, 1, 0)
	if err != nil || len(rules) == 0 {
		return nil, errors.New("policy rule " + req.PolicyId + " not found")
	}
	// Find the specific rule by ID via a targeted list
	existing, err := s.policyRepo.List(ctx, "", false, 200, 0)
	if err != nil {
		return nil, errors.New("fetch policy rules: " + err.Error())
	}
	var pr *authzentityv1.PolicyRule
	for _, r := range existing {
		if r.PolicyId == req.PolicyId {
			pr = r
			break
		}
	}
	if pr == nil {
		return nil, errors.New("policy rule " + req.PolicyId + " not found")
	}

	// Remove old Casbin p-rule
	oldEffect := "allow"
	if pr.Effect == authzentityv1.PolicyEffect_POLICY_EFFECT_DENY {
		oldEffect = "deny"
	}
	_ = s.enforcer.RemovePolicy(pr.Subject, pr.Domain, pr.Object, pr.Action)
	_ = oldEffect // was used for the remove call above

	// Apply updates
	if req.Action != "" {
		pr.Action = req.Action
	}
	if req.Effect != authzentityv1.PolicyEffect_POLICY_EFFECT_UNSPECIFIED {
		pr.Effect = req.Effect
	}
	if req.Condition != "" {
		pr.Condition = req.Condition
	}
	if req.Description != "" {
		pr.Description = req.Description
	}
	pr.IsActive = req.IsActive
	pr.UpdatedAt = timestamppb.Now()

	// Persist update
	if _, err := s.policyRepo.Update(ctx, pr); err != nil {
		return nil, errors.New("update policy rule: " + err.Error())
	}

	// Re-add Casbin p-rule with new values
	newEffect := "allow"
	if pr.Effect == authzentityv1.PolicyEffect_POLICY_EFFECT_DENY {
		newEffect = "deny"
	}
	if pr.IsActive {
		if err := s.enforcer.AddPolicy(pr.Subject, pr.Domain, pr.Object, pr.Action, newEffect); err != nil {
			return nil, errors.New("re-add casbin policy: " + err.Error())
		}
	}

	return &authzservicev1.UpdatePolicyRuleResponse{Policy: pr}, nil
}

func (s *AuthZService) DeletePolicyRule(ctx context.Context, req *authzservicev1.DeletePolicyRuleRequest) (*authzservicev1.DeletePolicyRuleResponse, error) {
	if err := s.policyRepo.SoftDelete(ctx, req.PolicyId); err != nil {
		return nil, errors.New("soft delete policy_rule: " + err.Error())
	}
	// Note: Casbin rule is NOT auto-removed here; call InvalidatePolicyCache to reload.
	return &authzservicev1.DeletePolicyRuleResponse{Message: "policy rule deleted"}, nil
}

func (s *AuthZService) ListPolicyRules(ctx context.Context, req *authzservicev1.ListPolicyRulesRequest) (*authzservicev1.ListPolicyRulesResponse, error) {
	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 50
	}
	rules, err := s.policyRepo.List(ctx, req.Domain, req.ActiveOnly, pageSize, 0 /* page_token cursor not yet impl */)
	if err != nil {
		return nil, errors.New("list policy rules: " + err.Error())
	}
	return &authzservicev1.ListPolicyRulesResponse{Policies: rules}, nil
}

// ── Portal Configuration ──────────────────────────────────────────────────────

func (s *AuthZService) GetPortalConfig(ctx context.Context, req *authzservicev1.GetPortalConfigRequest) (*authzservicev1.GetPortalConfigResponse, error) {
	pc, err := s.portalRepo.GetByPortal(ctx, req.Portal)
	if err != nil {
		return nil, errors.New("get portal config: " + err.Error())
	}
	return &authzservicev1.GetPortalConfigResponse{Config: pc}, nil
}

func (s *AuthZService) UpdatePortalConfig(ctx context.Context, req *authzservicev1.UpdatePortalConfigRequest) (*authzservicev1.UpdatePortalConfigResponse, error) {
	// UpdatePortalConfigRequest has flat fields (no .Config sub-message).
	// Build a PortalConfig entity from the flat request fields.
	pc := &authzentityv1.PortalConfig{
		Portal:                  req.Portal,
		MfaRequired:             req.MfaRequired,
		MfaMethods:              req.MfaMethods,
		AccessTokenTtlSeconds:   req.AccessTokenTtlSeconds,
		RefreshTokenTtlSeconds:  req.RefreshTokenTtlSeconds,
		SessionTtlSeconds:       req.SessionTtlSeconds,
		IdleTimeoutSeconds:      req.IdleTimeoutSeconds,
		AllowConcurrentSessions: req.AllowConcurrentSessions,
		MaxConcurrentSessions:   req.MaxConcurrentSessions,
		UpdatedBy:               req.UpdatedBy,
		UpdatedAt:               timestamppb.Now(),
	}
	if _, err := s.portalRepo.Upsert(ctx, pc); err != nil {
		return nil, errors.New("upsert portal config: " + err.Error())
	}
	if s.publisher != nil {
		_ = s.publisher.PublishPortalConfigUpdated(ctx, req.Portal, pc.MfaRequired,
			pc.AccessTokenTtlSeconds, pc.RefreshTokenTtlSeconds, pc.SessionTtlSeconds,
			pc.IdleTimeoutSeconds, req.UpdatedBy)
	}
	return &authzservicev1.UpdatePortalConfigResponse{Config: pc}, nil
}

// ── Audit ─────────────────────────────────────────────────────────────────────

func (s *AuthZService) ListAccessDecisionAudits(ctx context.Context, req *authzservicev1.ListAccessDecisionAuditsRequest) (*authzservicev1.ListAccessDecisionAuditsResponse, error) {
	audits, total, err := s.auditRepo.List(ctx, req)
	if err != nil {
		return nil, errors.New("list audits: " + err.Error())
	}
	return &authzservicev1.ListAccessDecisionAuditsResponse{
		Audits:     audits,
		TotalCount: int32(total),
	}, nil
}

// ── Cache Invalidation ────────────────────────────────────────────────────────

func (s *AuthZService) InvalidatePolicyCache(ctx context.Context, req *authzservicev1.InvalidatePolicyCacheRequest) (*authzservicev1.InvalidatePolicyCacheResponse, error) {
	if err := s.enforcer.InvalidateCache(); err != nil {
		return nil, errors.New("invalidate cache: " + err.Error())
	}
	if s.publisher != nil {
		_ = s.publisher.PublishPolicyCacheInvalidated(ctx, req.Domain, "")
	}
	return &authzservicev1.InvalidatePolicyCacheResponse{
		Invalidated: true,
	}, nil
}

// ── JWKS ──────────────────────────────────────────────────────────────────────

// GetJWKS returns the RS256 public key set for JWT verification.
// The public key PEM is read from the path configured in authn service.
// AuthZ service serves it here so gateway can cache and serve /.well-known/jwks.json.
func (s *AuthZService) GetJWKS(ctx context.Context, req *authzservicev1.GetJWKSRequest) (*authzservicev1.GetJWKSResponse, error) {
	// The JWKS public key path is supplied via AUTHZ_JWKS_PUBLIC_KEY_PATH env var.
	// This is the same RS256 public key generated alongside the authn private key.
	pubKeyPath := getEnvOrDefault("AUTHZ_JWKS_PUBLIC_KEY_PATH", "/secrets/jwt_rsa_public.pem")
	kidValue := getEnvOrDefault("JWT_KEY_ID", "insuretech-2025-01")

	pemBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, errors.New("read public key: " + err.Error())
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM block in public key file")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse public key: " + err.Error())
	}
	rsaKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	nBytes := rsaKey.N.Bytes()
	eBytes := make([]byte, 4)
	eBytes[0] = byte(rsaKey.E >> 24)
	eBytes[1] = byte(rsaKey.E >> 16)
	eBytes[2] = byte(rsaKey.E >> 8)
	eBytes[3] = byte(rsaKey.E)
	// Trim leading zeros from exponent
	i := 0
	for i < len(eBytes)-1 && eBytes[i] == 0 {
		i++
	}
	eBytes = eBytes[i:]

	return &authzservicev1.GetJWKSResponse{
		Keys: []*authzservicev1.JWK{{
			Kty: "RSA",
			Use: "sig",
			Alg: "RS256",
			Kid: kidValue,
			N:   base64.RawURLEncoding.EncodeToString(nBytes),
			E:   base64.RawURLEncoding.EncodeToString(eBytes),
		}},
	}, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func extractContext(ctx *authzservicev1.AccessContext) (ipAddr, userAgent, sessionID, tokenID, deviceID string) {
	if ctx == nil {
		return
	}
	return ctx.IpAddress, ctx.UserAgent, ctx.SessionId, ctx.TokenId, ctx.DeviceId
}
