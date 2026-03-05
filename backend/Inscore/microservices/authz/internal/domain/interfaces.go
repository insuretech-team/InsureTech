package domain

import (
	"context"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EnforcerIface abstracts the Casbin enforcer so it can be mocked in tests.
type EnforcerIface interface {
	// Enforce checks: does subject have permission to perform action on object in domain?
	// sub = "user:<uuid>", dom = "portal:tenant", obj = "svc:policy/create", act = "POST"
	// Returns (allowed, matchedRule, error).
	Enforce(ctx context.Context, sub, dom, obj, act string) (bool, string, error)

	// AddPolicy adds a p-type rule to the enforcer and DB.
	AddPolicy(sub, dom, obj, act, effect string) error

	// RemovePolicy removes a p-type rule.
	RemovePolicy(sub, dom, obj, act string) error

	// AddRoleForUserInDomain assigns role to user in domain (g-type rule).
	// sub = "user:<uuid>", role = "role:<name>", dom = "portal:tenant"
	AddRoleForUserInDomain(sub, role, dom string) error

	// DeleteRoleForUserInDomain removes role from user in domain.
	DeleteRoleForUserInDomain(sub, role, dom string) error

	// GetRolesForUserInDomain returns all roles for a user in a domain.
	GetRolesForUserInDomain(sub, dom string) ([]string, error)

	// GetPermissionsForUserInDomain returns effective permissions (all matched p-lines).
	GetPermissionsForUserInDomain(sub, dom string) ([][]string, error)

	// InvalidateCache forces a full policy reload from DB.
	InvalidateCache() error
}

// RoleRepository manages roles in authz_schema.roles.
type RoleRepository interface {
	Create(ctx context.Context, role *authzentityv1.Role) (*authzentityv1.Role, error)
	GetByID(ctx context.Context, roleID string) (*authzentityv1.Role, error)
	GetByNameAndPortal(ctx context.Context, name string, portal authzentityv1.Portal) (*authzentityv1.Role, error)
	Update(ctx context.Context, role *authzentityv1.Role) (*authzentityv1.Role, error)
	SoftDelete(ctx context.Context, roleID string) error
	List(ctx context.Context, portal authzentityv1.Portal, activeOnly bool, limit, offset int) ([]*authzentityv1.Role, error)
}

// UserRoleRepository manages user_roles in authz_schema.
type UserRoleRepository interface {
	Assign(ctx context.Context, ur *authzentityv1.UserRole) (*authzentityv1.UserRole, error)
	Remove(ctx context.Context, userID, roleID, domain string) error
	ListByUser(ctx context.Context, userID, domain string) ([]*authzentityv1.UserRole, error)
}

// PolicyRuleRepository manages policy_rules in authz_schema.
// Writes to this repo also sync to casbin_rules via the enforcer.
type PolicyRuleRepository interface {
	Create(ctx context.Context, pr *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error)
	Update(ctx context.Context, pr *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error)
	SoftDelete(ctx context.Context, policyID string) error
	List(ctx context.Context, domain string, activeOnly bool, limit, offset int) ([]*authzentityv1.PolicyRule, error)
}

// PortalConfigRepository manages portal_configs.
type PortalConfigRepository interface {
	GetByPortal(ctx context.Context, portal authzentityv1.Portal) (*authzentityv1.PortalConfig, error)
	Upsert(ctx context.Context, pc *authzentityv1.PortalConfig) (*authzentityv1.PortalConfig, error)
}

// AuditRepository writes access decision audit rows.
type AuditRepository interface {
	Create(ctx context.Context, audit *authzentityv1.AccessDecisionAudit) error
	List(ctx context.Context, req *authzservicev1.ListAccessDecisionAuditsRequest) ([]*authzentityv1.AccessDecisionAudit, int64, error)
}

// TokenConfigRepository manages JWT signing key configuration.
type TokenConfigRepository interface {
	GetActive(ctx context.Context) (*authzentityv1.TokenConfig, error)
	List(ctx context.Context) ([]*authzentityv1.TokenConfig, error)
	Create(ctx context.Context, cfg *authzentityv1.TokenConfig) (*authzentityv1.TokenConfig, error)
}

// PortalConfigRepositoryExt extends PortalConfigRepository with List.
type PortalConfigRepositoryExt interface {
	PortalConfigRepository
	List(ctx context.Context) ([]*authzentityv1.PortalConfig, error)
}

// SeedResult holds counts from a portal seeding operation.
type SeedResult struct {
	RolesSeeded    int
	PoliciesSeeded int
	RulesSeeded    int
}

// ── Input types for authz_repository.go ──────────────────────────────────────

// ListCasbinRulesInput filters for listing raw Casbin rules.
type ListCasbinRulesInput struct {
	PType    string // "p" or "g"
	Domain   string // "portal:tenant_id"
	Subject  string // "user:<uuid>" or "role:<name>"
	Resource string // "svc:policy/create"
}

// ListRolesInput filters for listing roles.
type ListRolesInput struct {
	Portal   authzentityv1.Portal
	TenantID string
	PageSize int32
}

// AssignUserRoleInput holds parameters for assigning a role to a user.
type AssignUserRoleInput struct {
	UserID   string
	RoleID   string
	TenantID string
	Portal   authzentityv1.Portal
}

// ListPoliciesInput filters for listing policy rules.
type ListPoliciesInput struct {
	Portal   authzentityv1.Portal
	TenantID string
	RoleID   string
	PageSize int32
}

// ListAccessAuditInput filters for listing access decision audit rows.
type ListAccessAuditInput struct {
	SubjectID string
	TenantID  string
	Portal    authzentityv1.Portal
	Allowed   *bool
	FromTime  *timestamppb.Timestamp
	ToTime    *timestamppb.Timestamp
	PageSize  int32
}
