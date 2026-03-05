package enforcer

// casbin_enforcer.go — wraps the Casbin SyncedEnforcer with the PERM model.
//
// PERM model (built-in, no external .conf needed):
//   [request_definition]  r = sub, dom, obj, act
//   [policy_definition]   p = sub, dom, obj, act, eft
//   [role_definition]     g = _, _, _          (user → role, domain-scoped)
//   [policy_effect]       e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
//   [matchers]            m = g(r.sub, p.sub, r.dom) && r.dom == p.dom &&
//                             keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
//
// Domain format:  "portal:tenant_id"  e.g. "system:root", "agent:tenant-abc"
// Subject format: "user:<uuid>"       e.g. "user:550e8400-..."
// Role format:    "role:<name>"       e.g. "role:underwriter"
// Object format:  "svc:<svc>/<res>"   e.g. "svc:policy/create"
// Action format:  HTTP verb or *      e.g. "POST", "GET", "*"
// Effect:         "allow" | "deny"

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"gorm.io/gorm"
)

const builtinModel = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act, eft

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`

// CasbinEnforcer wraps casbin.SyncedEnforcer and implements domain.EnforcerIface.
type CasbinEnforcer struct {
	enforcer *casbin.SyncedEnforcer
	mu       sync.RWMutex
}

// New creates a CasbinEnforcer backed by the gorm-adapter (casbin_rules table in authz_schema).
// modelPath: path to a .conf file; empty = use built-in PERM model above.
func New(db *gorm.DB, modelPath string) (*CasbinEnforcer, error) {
	// gorm-adapter: table = casbin_rules, schema = authz_schema, auto-create table.
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &entityv1.CasbinRule{}, "casbin_rules")
	if err != nil {
		return nil, errors.New("casbin gorm-adapter init: " + err.Error())
	}

	var m model.Model
	if modelPath != "" {
		m, err = model.NewModelFromFile(modelPath)
		if err != nil {
			return nil, errors.New("casbin model from file " + modelPath + ": " + err.Error())
		}
	} else {
		m, err = model.NewModelFromString(builtinModel)
		if err != nil {
			return nil, errors.New("casbin built-in model: " + err.Error())
		}
	}

	e, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, errors.New("casbin enforcer init: " + err.Error())
	}

	// Load policies from DB on startup.
	if err := e.LoadPolicy(); err != nil {
		return nil, errors.New("casbin load policy: " + err.Error())
	}

	return &CasbinEnforcer{enforcer: e}, nil
}

// StartAutoReload starts background policy reloading every interval.
func (e *CasbinEnforcer) StartAutoReload(intervalSeconds int) {
	e.enforcer.StartAutoLoadPolicy(time.Duration(intervalSeconds) * time.Second)
}

// Enforce implements domain.EnforcerIface.
// Returns (allowed, matchedRule, error). Deny-by-default: no matching rule → deny.
func (e *CasbinEnforcer) Enforce(_ context.Context, sub, dom, obj, act string) (bool, string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	allowed, err := e.enforcer.Enforce(sub, dom, obj, act)
	if err != nil {
		return false, "", errors.New("casbin enforce: " + err.Error())
	}
	matchedRule := ""
	if allowed {
		// Retrieve the matched rule for audit logging.
		rules, _ := e.enforcer.GetFilteredPolicy(0, sub)
		if len(rules) > 0 {
			matchedRule = "p: " + sub + ", " + dom + ", " + obj + ", " + act
		}
	}
	return allowed, matchedRule, nil
}

// AddPolicy implements domain.EnforcerIface.
func (e *CasbinEnforcer) AddPolicy(sub, dom, obj, act, effect string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	ok, err := e.enforcer.AddPolicy(sub, dom, obj, act, effect)
	if err != nil {
		return errors.New("casbin AddPolicy: " + err.Error())
	}
	if !ok {
		return errors.New("casbin AddPolicy: rule already exists")
	}
	return nil
}

// RemovePolicy implements domain.EnforcerIface.
func (e *CasbinEnforcer) RemovePolicy(sub, dom, obj, act string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, err := e.enforcer.RemoveFilteredPolicy(0, sub, dom, obj, act)
	return err
}

// AddRoleForUserInDomain implements domain.EnforcerIface.
func (e *CasbinEnforcer) AddRoleForUserInDomain(sub, role, dom string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	ok, err := e.enforcer.AddRoleForUserInDomain(sub, role, dom)
	if err != nil {
		return errors.New("casbin AddRoleForUserInDomain: " + err.Error())
	}
	if !ok {
		return errors.New("casbin AddRoleForUserInDomain: assignment already exists")
	}
	return nil
}

// DeleteRoleForUserInDomain implements domain.EnforcerIface.
func (e *CasbinEnforcer) DeleteRoleForUserInDomain(sub, role, dom string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	ok, err := e.enforcer.DeleteRoleForUserInDomain(sub, role, dom)
	if err != nil {
		return errors.New("casbin DeleteRoleForUserInDomain: " + err.Error())
	}
	if !ok {
		return errors.New("casbin DeleteRoleForUserInDomain: role assignment does not exist")
	}
	return nil
}

// GetRolesForUserInDomain implements domain.EnforcerIface.
func (e *CasbinEnforcer) GetRolesForUserInDomain(sub, dom string) ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.GetRolesForUserInDomain(sub, dom), nil
}

// GetPermissionsForUserInDomain implements domain.EnforcerIface.
func (e *CasbinEnforcer) GetPermissionsForUserInDomain(sub, dom string) ([][]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.GetPermissionsForUserInDomain(sub, dom), nil
}

// InvalidateCache forces a full policy reload from DB.
func (e *CasbinEnforcer) InvalidateCache() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enforcer.LoadPolicy()
}
