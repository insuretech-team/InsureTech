package service

import (
	"testing"
)

func TestGenerateAPIKey(t *testing.T) {
	// Test multiple generations to ensure randomness
	keys := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		rawKey, keyHash, err := generateAPIKey()
		
		if err != nil {
			t.Fatalf("generateAPIKey() failed: %v", err)
		}
		
		// Check prefix
		if len(rawKey) < 4 || rawKey[:4] != "isk_" {
			t.Errorf("generateAPIKey() rawKey should start with 'isk_', got: %s", rawKey[:4])
		}
		
		// Check length (32 bytes base64 encoded + prefix should be ~48+ chars)
		if len(rawKey) < 40 {
			t.Errorf("generateAPIKey() rawKey too short: %d chars", len(rawKey))
		}
		
		// Check hash length (SHA-256 = 64 hex chars)
		if len(keyHash) != 64 {
			t.Errorf("generateAPIKey() keyHash should be 64 chars, got: %d", len(keyHash))
		}
		
		// Check uniqueness
		if keys[rawKey] {
			t.Errorf("generateAPIKey() generated duplicate key")
		}
		keys[rawKey] = true
		
		// Verify hash is deterministic
		rawKey2, keyHash2, err := generateAPIKey()
		if err != nil {
			t.Fatalf("generateAPIKey() failed on second call: %v", err)
		}
		
		// Different keys should have different hashes
		if rawKey == rawKey2 {
			t.Error("generateAPIKey() should generate unique keys")
		}
		if keyHash == keyHash2 && rawKey != rawKey2 {
			t.Error("generateAPIKey() different keys should have different hashes")
		}
	}
	
	// Should have generated 100 unique keys
	if len(keys) != 100 {
		t.Errorf("Expected 100 unique keys, got %d", len(keys))
	}
}

func TestGenerateAPIKey_Format(t *testing.T) {
	rawKey, keyHash, err := generateAPIKey()
	
	if err != nil {
		t.Fatalf("generateAPIKey() error = %v", err)
	}
	
	// Test rawKey format
	t.Run("rawKey has correct prefix", func(t *testing.T) {
		if rawKey[:4] != "isk_" {
			t.Errorf("rawKey prefix = %s, want isk_", rawKey[:4])
		}
	})
	
	// Test keyHash format (should be hex)
	t.Run("keyHash is valid hex", func(t *testing.T) {
		for _, ch := range keyHash {
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f')) {
				t.Errorf("keyHash contains invalid hex char: %c", ch)
			}
		}
	})
}

// Note: Full integration tests for RotateAPIKey would require:
// - Mock repositories
// - Mock event publishers
// - Database setup
// These should be in integration test suite, not unit tests.
// Here we just test the helper function.
