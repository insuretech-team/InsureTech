// Code generated supplement - DO NOT EDIT manually.
// Adds GORM TableName() methods to proto-generated authz entity structs.
package entityv1

func (*Role) TableName() string                { return "authz_roles" }
func (*CasbinRule) TableName() string          { return "authz_casbin_rules" }
func (*UserRole) TableName() string            { return "authz_user_roles" }
func (*PolicyRule) TableName() string          { return "authz_policy_rules" }
func (*PortalConfig) TableName() string        { return "authz_portal_configs" }
func (*TokenConfig) TableName() string         { return "authz_token_configs" }
func (*AccessDecisionAudit) TableName() string { return "authz_access_decision_audits" }
