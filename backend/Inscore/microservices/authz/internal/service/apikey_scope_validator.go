package service

import (
	"fmt"
	"strings"
)

// ApiKeyScopeValidator validates API key scopes against requested permissions
type ApiKeyScopeValidator struct{}

// NewApiKeyScopeValidator creates a new API key scope validator
func NewApiKeyScopeValidator() *ApiKeyScopeValidator {
	return &ApiKeyScopeValidator{}
}

// ValidateScope checks if the API key has the required scope for the requested action
// Scopes format: "svc:policy/read", "svc:policy/*", "svc:*/read", "*"
// Returns (allowed, reason)
func (v *ApiKeyScopeValidator) ValidateScope(scopes []string, object, action string) (bool, string) {
	if len(scopes) == 0 {
		return false, "API key has no scopes defined"
	}
	
	// Build the required permission from object and action
	// object: "svc:policy/create" -> "svc:policy/create"
	// action: "POST" -> we already have the object with action embedded
	requiredPermission := object
	
	// Check each scope
	for _, scope := range scopes {
		if v.matchesScope(scope, requiredPermission, action) {
			return true, ""
		}
	}
	
	return false, fmt.Sprintf("API key does not have required scope: %s", requiredPermission)
}

// matchesScope checks if a scope pattern matches the required permission
func (v *ApiKeyScopeValidator) matchesScope(scope, requiredPermission, action string) bool {
	// Wildcard: "*" allows everything
	if scope == "*" {
		return true
	}
	
	// Exact match
	if scope == requiredPermission {
		return true
	}
	
	// Wildcard patterns
	// "svc:policy/*" matches "svc:policy/read", "svc:policy/create", etc.
	if strings.HasSuffix(scope, "/*") {
		prefix := strings.TrimSuffix(scope, "/*")
		if strings.HasPrefix(requiredPermission, prefix+"/") {
			return true
		}
	}
	
	// "svc:*/read" matches "svc:policy/read", "svc:claim/read", etc.
	if strings.Contains(scope, "*/") {
		parts := strings.Split(scope, "*/")
		if len(parts) == 2 {
			prefix := parts[0]
			suffix := parts[1]
			if strings.HasPrefix(requiredPermission, prefix) && strings.HasSuffix(requiredPermission, "/"+suffix) {
				return true
			}
		}
	}
	
	// Resource-level wildcard: "svc:policy" matches any action on policy
	if !strings.Contains(scope, "/") && strings.HasPrefix(requiredPermission, scope+"/") {
		return true
	}
	
	return false
}

// ParseScopesFromAttributes extracts API key scopes from AccessContext attributes
func (v *ApiKeyScopeValidator) ParseScopesFromAttributes(attributes map[string]string) []string {
	if attributes == nil {
		return nil
	}
	
	scopesStr, exists := attributes["api_key_scopes"]
	if !exists || scopesStr == "" {
		return nil
	}
	
	// Scopes are comma-separated in the attribute
	scopes := strings.Split(scopesStr, ",")
	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		trimmed := strings.TrimSpace(scope)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// Common API key scope patterns for different owner types
var (
	// INSURER scopes - full access to their tenant's data
	InsurerScopes = []string{
		"svc:policy/*",
		"svc:claim/*",
		"svc:payment/*",
		"svc:customer/read",
		"svc:agent/read",
		"svc:report/*",
	}
	
	// PARTNER scopes - limited to partner-specific operations
	PartnerScopes = []string{
		"svc:policy/read",
		"svc:policy/create",
		"svc:customer/read",
		"svc:customer/create",
		"svc:commission/read",
	}
	
	// INTERNAL scopes - service-to-service communication
	InternalScopes = []string{
		"*", // Internal services get full access
	}
)

// GetDefaultScopesForOwnerType returns default scopes based on API key owner type
func GetDefaultScopesForOwnerType(ownerType string) []string {
	switch ownerType {
	case "API_KEY_OWNER_TYPE_INSURER":
		return InsurerScopes
	case "API_KEY_OWNER_TYPE_PARTNER":
		return PartnerScopes
	case "API_KEY_OWNER_TYPE_INTERNAL":
		return InternalScopes
	default:
		return []string{} // No default scopes for unknown types
	}
}
