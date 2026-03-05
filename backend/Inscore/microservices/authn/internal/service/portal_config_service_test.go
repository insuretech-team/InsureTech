package service

import (
	"context"
	"testing"
	"time"
)

func TestPortalConfigCache_GetDefaultConfig(t *testing.T) {
	cache := &PortalConfigCache{
		configs: make(map[string]*PortalConfig),
		ttl:     5 * time.Minute,
	}

	tests := []struct {
		name                 string
		portal               string
		expectedMFARequired  bool
		expectedMinLength    int32
		expectedSessionTTL   int32
	}{
		{
			name:                "customer portal",
			portal:              "customer",
			expectedMFARequired: false,
			expectedMinLength:   8,
			expectedSessionTTL:  3600,
		},
		{
			name:                "agent portal",
			portal:              "agent",
			expectedMFARequired: true,
			expectedMinLength:   10,
			expectedSessionTTL:  7200,
		},
		{
			name:                "business portal",
			portal:              "business",
			expectedMFARequired: false,
			expectedMinLength:   8,
			expectedSessionTTL:  7200,
		},
		{
			name:                "system portal",
			portal:              "system",
			expectedMFARequired: true,
			expectedMinLength:   12,
			expectedSessionTTL:  14400,
		},
		{
			name:                "unknown portal",
			portal:              "unknown",
			expectedMFARequired: false,
			expectedMinLength:   8,
			expectedSessionTTL:  3600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := cache.getDefaultConfig(tt.portal)

			if config.PortalName != tt.portal {
				t.Errorf("PortalName = %v, want %v", config.PortalName, tt.portal)
			}

			if config.MFARequired != tt.expectedMFARequired {
				t.Errorf("MFARequired = %v, want %v", config.MFARequired, tt.expectedMFARequired)
			}

			if config.PasswordMinLength != tt.expectedMinLength {
				t.Errorf("PasswordMinLength = %v, want %v", config.PasswordMinLength, tt.expectedMinLength)
			}

			if config.SessionTTLSeconds != tt.expectedSessionTTL {
				t.Errorf("SessionTTLSeconds = %v, want %v", config.SessionTTLSeconds, tt.expectedSessionTTL)
			}
		})
	}
}

func TestPortalConfigCache_ValidatePassword(t *testing.T) {
	cache := &PortalConfigCache{
		configs: make(map[string]*PortalConfig),
		ttl:     5 * time.Minute,
	}

	ctx := context.Background()

	tests := []struct {
		name      string
		portal    string
		password  string
		wantError bool
	}{
		{
			name:      "valid customer password",
			portal:    "customer",
			password:  "Password123",
			wantError: false,
		},
		{
			name:      "customer password too short",
			portal:    "customer",
			password:  "Pass1",
			wantError: true,
		},
		{
			name:      "customer password no uppercase",
			portal:    "customer",
			password:  "password123",
			wantError: true,
		},
		{
			name:      "customer password no lowercase",
			portal:    "customer",
			password:  "PASSWORD123",
			wantError: true,
		},
		{
			name:      "customer password no digit",
			portal:    "customer",
			password:  "PasswordABC",
			wantError: true,
		},
		{
			name:      "valid agent password with symbol",
			portal:    "agent",
			password:  "P@ssw0rd123!",
			wantError: false,
		},
		{
			name:      "agent password no symbol",
			portal:    "agent",
			password:  "Password123",
			wantError: true,
		},
		{
			name:      "agent password too short",
			portal:    "agent",
			password:  "P@ss123!",
			wantError: true,
		},
		{
			name:      "valid system password",
			portal:    "system",
			password:  "Syst3m!P@ssw0rd",
			wantError: false,
		},
		{
			name:      "system password too short",
			portal:    "system",
			password:  "Sys!P@ss1",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cache.ValidatePassword(ctx, tt.portal, tt.password)

			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePassword() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPortalConfigCache_Invalidate(t *testing.T) {
	cache := &PortalConfigCache{
		configs: make(map[string]*PortalConfig),
		ttl:     5 * time.Minute,
	}

	// Add some configs
	cache.configs["customer"] = &PortalConfig{PortalName: "customer"}
	cache.configs["agent"] = &PortalConfig{PortalName: "agent"}

	if len(cache.configs) != 2 {
		t.Fatalf("Expected 2 configs, got %d", len(cache.configs))
	}

	// Invalidate one
	cache.Invalidate("customer")

	if len(cache.configs) != 1 {
		t.Errorf("Expected 1 config after invalidation, got %d", len(cache.configs))
	}

	if _, exists := cache.configs["customer"]; exists {
		t.Error("Customer config should be invalidated")
	}

	if _, exists := cache.configs["agent"]; !exists {
		t.Error("Agent config should still exist")
	}
}

func TestPortalConfigCache_InvalidateAll(t *testing.T) {
	cache := &PortalConfigCache{
		configs: make(map[string]*PortalConfig),
		ttl:     5 * time.Minute,
	}

	// Add some configs
	cache.configs["customer"] = &PortalConfig{PortalName: "customer"}
	cache.configs["agent"] = &PortalConfig{PortalName: "agent"}
	cache.configs["business"] = &PortalConfig{PortalName: "business"}

	if len(cache.configs) != 3 {
		t.Fatalf("Expected 3 configs, got %d", len(cache.configs))
	}

	// Invalidate all
	cache.InvalidateAll()

	if len(cache.configs) != 0 {
		t.Errorf("Expected 0 configs after InvalidateAll, got %d", len(cache.configs))
	}
}

func TestPortalConfigCache_GetSessionTTL(t *testing.T) {
	cache := &PortalConfigCache{
		configs: make(map[string]*PortalConfig),
		ttl:     5 * time.Minute,
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		portal      string
		expectedTTL time.Duration
	}{
		{
			name:        "customer portal",
			portal:      "customer",
			expectedTTL: 3600 * time.Second, // 1 hour
		},
		{
			name:        "agent portal",
			portal:      "agent",
			expectedTTL: 7200 * time.Second, // 2 hours
		},
		{
			name:        "system portal",
			portal:      "system",
			expectedTTL: 14400 * time.Second, // 4 hours
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttl, err := cache.GetSessionTTL(ctx, tt.portal)

			if err != nil {
				t.Errorf("GetSessionTTL() error = %v", err)
				return
			}

			if ttl != tt.expectedTTL {
				t.Errorf("GetSessionTTL() = %v, want %v", ttl, tt.expectedTTL)
			}
		})
	}
}

func TestPortalConfigCache_IsMFARequired(t *testing.T) {
	cache := &PortalConfigCache{
		configs: make(map[string]*PortalConfig),
		ttl:     5 * time.Minute,
	}

	ctx := context.Background()

	tests := []struct {
		name             string
		portal           string
		expectedRequired bool
	}{
		{
			name:             "customer portal - MFA not required",
			portal:           "customer",
			expectedRequired: false,
		},
		{
			name:             "agent portal - MFA required",
			portal:           "agent",
			expectedRequired: true,
		},
		{
			name:             "system portal - MFA required",
			portal:           "system",
			expectedRequired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			required, err := cache.IsMFARequired(ctx, tt.portal)

			if err != nil {
				t.Errorf("IsMFARequired() error = %v", err)
				return
			}

			if required != tt.expectedRequired {
				t.Errorf("IsMFARequired() = %v, want %v", required, tt.expectedRequired)
			}
		})
	}
}
