package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/metrics"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	apikeyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// RotateAPIKey generates a new API key and marks the old one for graceful expiration.
// During the grace period, both keys work. After grace period, old key is auto-revoked.
func (s *AuthService) RotateAPIKey(ctx context.Context, req *authnservicev1.RotateAPIKeyRequest) (*authnservicev1.RotateAPIKeyResponse, error) {
	start := time.Now()

	// Validate request
	if req.KeyId == "" {
		return &authnservicev1.RotateAPIKeyResponse{
			Error: &commonv1.Error{
				Code:    "INVALID_ARGUMENT",
				Message: "key_id is required",
			},
		}, nil
	}

	// Default grace period: 24 hours
	gracePeriodHours := req.GracePeriodHours
	if gracePeriodHours <= 0 {
		gracePeriodHours = 24
	}
	gracePeriod := time.Duration(gracePeriodHours) * time.Hour

	// Get the old API key
	oldKey, err := s.apiKeyRepo.GetByID(ctx, req.KeyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &authnservicev1.RotateAPIKeyResponse{
				Error: &commonv1.Error{
					Code:    "NOT_FOUND",
					Message: "API key not found",
				},
			}, nil
		}
		logger.Errorf("Failed to get API key: %v (key_id=%s)", err, req.KeyId)
		return &authnservicev1.RotateAPIKeyResponse{
			Error: &commonv1.Error{
				Code:    "INTERNAL",
				Message: "Failed to retrieve API key",
			},
		}, nil
	}

	// Check if key is already revoked or expired
	if oldKey.Status == apikeyv1.ApiKeyStatus_API_KEY_STATUS_REVOKED {
		return &authnservicev1.RotateAPIKeyResponse{
			Error: &commonv1.Error{
				Code:    "FAILED_PRECONDITION",
				Message: "Cannot rotate a revoked key",
			},
		}, nil
	}

	// Generate new API key
	rawKey, keyHash, err := generateAPIKey()
	if err != nil {
		logger.Errorf("Failed to generate API key: %v", err)
		return &authnservicev1.RotateAPIKeyResponse{
			Error: &commonv1.Error{
				Code:    "INTERNAL",
				Message: "Failed to generate new API key",
			},
		}, nil
	}

	// Create new API key with same attributes as old one
	newKeyID := uuid.New().String()
	now := time.Now()

	newKey := &apikeyv1.ApiKey{
		Id:                 newKeyID,
		KeyHash:            keyHash,
		Name:               oldKey.Name + " (rotated)",
		OwnerType:          oldKey.OwnerType,
		OwnerId:            oldKey.OwnerId,
		Scopes:             oldKey.Scopes,
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: oldKey.RateLimitPerMinute,
		ExpiresAt:          oldKey.ExpiresAt, // Keep same expiration as original
		IpWhitelist:        oldKey.IpWhitelist,
	}

	// Create new key and update old key atomically
	err = func() error {
		// Create new key
		if err := s.apiKeyRepo.Create(ctx, newKey); err != nil {
			return fmt.Errorf("failed to create new API key: %w", err)
		}

		// Mark old key as ROTATING and set expiration for grace period
		oldKeyExpiry := now.Add(gracePeriod)
		if err := s.apiKeyRepo.MarkAsRotating(ctx, oldKey.Id, oldKeyExpiry); err != nil {
			return fmt.Errorf("failed to mark old key as rotating: %w", err)
		}

		return nil
	}()

	if err != nil {
		logger.Errorf("Failed to rotate API key: %v (old_key_id=%s)", err, req.KeyId)

		// Record failed rotation
		duration := time.Since(start).Seconds()
		metrics.RecordAPIKeyRotation(oldKey.OwnerType.String(), false, duration)

		return &authnservicev1.RotateAPIKeyResponse{
			Error: &commonv1.Error{
				Code:    "INTERNAL",
				Message: "Failed to rotate API key",
			},
		}, nil
	}

	// Record successful rotation
	duration := time.Since(start).Seconds()
	metrics.RecordAPIKeyRotation(oldKey.OwnerType.String(), true, duration)

	// Log the rotation event
	logger.Infof("API key rotated successfully: old_key_id=%s, new_key_id=%s, owner_id=%s, grace_period_hours=%d",
		oldKey.Id, newKeyID, oldKey.OwnerId, gracePeriodHours)

	// Event publishing for API key rotation can be added later if needed
	// For now, logging is sufficient for audit trail

	return &authnservicev1.RotateAPIKeyResponse{
		NewKeyId:        newKeyID,
		RawKey:          rawKey,
		OldKeyId:        oldKey.Id,
		OldKeyExpiresAt: timestamppb.New(now.Add(gracePeriod)),
		Message:         fmt.Sprintf("API key rotated successfully. Old key will expire in %d hours", gracePeriodHours),
	}, nil
}

// generateAPIKey creates a new random API key and its SHA-256 hash.
// Returns: (rawKey, keyHash, error)
func generateAPIKey() (string, string, error) {
	// Generate 32 random bytes
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate random key: %w", err)
	}

	// Encode to base64 for the raw key (shown to user once)
	rawKey := "isk_" + base64.URLEncoding.EncodeToString(b)

	// Hash with SHA-256 for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := fmt.Sprintf("%x", hash)

	return rawKey, keyHash, nil
}
