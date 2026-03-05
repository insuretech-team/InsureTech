package service

import (
	"testing"
)

func TestApiKeyScopeValidator_ValidateScope(t *testing.T) {
	validator := NewApiKeyScopeValidator()

	tests := []struct {
		name           string
		scopes         []string
		object         string
		action         string
		expectedAllowed bool
		expectedReason string
	}{
		{
			name:           "exact match",
			scopes:         []string{"svc:policy/read", "svc:claim/write"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
		{
			name:           "wildcard all",
			scopes:         []string{"*"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
		{
			name:           "service wildcard match",
			scopes:         []string{"svc:policy/*"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
		{
			name:           "service wildcard no match",
			scopes:         []string{"svc:policy/*"},
			object:         "svc:claim/read",
			action:         "GET",
			expectedAllowed: false,
			expectedReason: "API key does not have required scope: svc:claim/read",
		},
		{
			name:           "action wildcard match",
			scopes:         []string{"svc:*/read"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
		{
			name:           "action wildcard match different service",
			scopes:         []string{"svc:*/read"},
			object:         "svc:claim/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
		{
			name:           "action wildcard no match",
			scopes:         []string{"svc:*/read"},
			object:         "svc:policy/write",
			action:         "POST",
			expectedAllowed: false,
			expectedReason: "API key does not have required scope: svc:policy/write",
		},
		{
			name:           "resource level match",
			scopes:         []string{"svc:policy"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
		{
			name:           "no scopes",
			scopes:         []string{},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: false,
			expectedReason: "API key has no scopes defined",
		},
		{
			name:           "no match",
			scopes:         []string{"svc:claim/read"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: false,
			expectedReason: "API key does not have required scope: svc:policy/read",
		},
		{
			name:           "multiple scopes with match",
			scopes:         []string{"svc:claim/read", "svc:policy/read", "svc:payment/write"},
			object:         "svc:policy/read",
			action:         "GET",
			expectedAllowed: true,
			expectedReason: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := validator.ValidateScope(tt.scopes, tt.object, tt.action)

			if allowed != tt.expectedAllowed {
				t.Errorf("ValidateScope() allowed = %v, want %v", allowed, tt.expectedAllowed)
			}

			if reason != tt.expectedReason {
				t.Errorf("ValidateScope() reason = %v, want %v", reason, tt.expectedReason)
			}
		})
	}
}

func TestApiKeyScopeValidator_ParseScopesFromAttributes(t *testing.T) {
	validator := NewApiKeyScopeValidator()

	tests := []struct {
		name           string
		attributes     map[string]string
		expectedScopes []string
	}{
		{
			name: "single scope",
			attributes: map[string]string{
				"api_key_scopes": "svc:policy/read",
			},
			expectedScopes: []string{"svc:policy/read"},
		},
		{
			name: "multiple scopes",
			attributes: map[string]string{
				"api_key_scopes": "svc:policy/read,svc:claim/write,svc:payment/*",
			},
			expectedScopes: []string{"svc:policy/read", "svc:claim/write", "svc:payment/*"},
		},
		{
			name: "scopes with spaces",
			attributes: map[string]string{
				"api_key_scopes": "svc:policy/read , svc:claim/write , svc:payment/*",
			},
			expectedScopes: []string{"svc:policy/read", "svc:claim/write", "svc:payment/*"},
		},
		{
			name:           "no api_key_scopes attribute",
			attributes:     map[string]string{"other": "value"},
			expectedScopes: nil,
		},
		{
			name:           "nil attributes",
			attributes:     nil,
			expectedScopes: nil,
		},
		{
			name: "empty scopes",
			attributes: map[string]string{
				"api_key_scopes": "",
			},
			expectedScopes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopes := validator.ParseScopesFromAttributes(tt.attributes)

			if len(scopes) != len(tt.expectedScopes) {
				t.Errorf("ParseScopesFromAttributes() returned %d scopes, want %d", len(scopes), len(tt.expectedScopes))
				return
			}

			for i, scope := range scopes {
				if scope != tt.expectedScopes[i] {
					t.Errorf("ParseScopesFromAttributes()[%d] = %v, want %v", i, scope, tt.expectedScopes[i])
				}
			}
		})
	}
}

func TestGetDefaultScopesForOwnerType(t *testing.T) {
	tests := []struct {
		name      string
		ownerType string
		wantCount int
	}{
		{
			name:      "insurer scopes",
			ownerType: "API_KEY_OWNER_TYPE_INSURER",
			wantCount: 6, // InsurerScopes has 6 default scopes
		},
		{
			name:      "partner scopes",
			ownerType: "API_KEY_OWNER_TYPE_PARTNER",
			wantCount: 5, // PartnerScopes has 5 default scopes
		},
		{
			name:      "internal scopes",
			ownerType: "API_KEY_OWNER_TYPE_INTERNAL",
			wantCount: 1, // InternalScopes has ["*"]
		},
		{
			name:      "unknown owner type",
			ownerType: "UNKNOWN",
			wantCount: 0, // Should return empty slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopes := GetDefaultScopesForOwnerType(tt.ownerType)

			if len(scopes) != tt.wantCount {
				t.Errorf("GetDefaultScopesForOwnerType() returned %d scopes, want %d", len(scopes), tt.wantCount)
			}

			// Verify internal gets wildcard
			if tt.ownerType == "API_KEY_OWNER_TYPE_INTERNAL" && len(scopes) > 0 {
				if scopes[0] != "*" {
					t.Errorf("Internal scopes should be [\"*\"], got %v", scopes)
				}
			}
		})
	}
}
