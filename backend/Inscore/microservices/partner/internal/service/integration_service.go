package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetPartnerAPICredentials requests integration keys bound to a partner profile
// Usually calls AuthN under the hood
func (s *PartnerService) GetPartnerAPICredentials(ctx context.Context, req *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error) {
	if req.PartnerId == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidArgument)
	}
	if s.authnClient == nil {
		return nil, fmt.Errorf("%w: authn client is not configured", ErrUnavailable)
	}

	keySummary, err := s.getFirstActivePartnerKey(ctx, req.PartnerId)
	if err != nil {
		return nil, err
	}
	// Existing key found: raw secret is not recoverable from AuthN.
	if keySummary != nil {
		return &partnerservicev1.GetPartnerAPICredentialsResponse{
			ApiKey:    keySummary.KeyId,
			ApiSecret: "",
			ExpiresAt: keySummary.ExpiresAt,
		}, nil
	}

	created, err := s.authnClient.CreateAPIKey(ctx, &authnservicev1.CreateAPIKeyRequest{
		Name:               "partner-" + req.PartnerId,
		OwnerId:            req.PartnerId,
		OwnerType:          "SERVICE",
		Scopes:             defaultPartnerAPIScopes(),
		RateLimitPerMinute: 300,
	})
	if err != nil {
		logger.Errorf("CreateAPIKey failed for partner=%s: %v", req.PartnerId, err)
		return nil, fmt.Errorf("provision API credentials: %w", err)
	}
	return &partnerservicev1.GetPartnerAPICredentialsResponse{
		ApiKey:    created.KeyId,
		ApiSecret: created.RawKey,
		ExpiresAt: created.ExpiresAt,
	}, nil
}

// RotatePartnerAPIKey invalidates old credentials and returns a new set
func (s *PartnerService) RotatePartnerAPIKey(ctx context.Context, req *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error) {
	if strings.TrimSpace(req.PartnerId) == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidArgument)
	}
	if s.authnClient == nil {
		return nil, fmt.Errorf("%w: authn client is not configured", ErrUnavailable)
	}

	keySummary, err := s.getFirstActivePartnerKey(ctx, req.PartnerId)
	if err != nil {
		return nil, err
	}
	if keySummary == nil {
		created, createErr := s.authnClient.CreateAPIKey(ctx, &authnservicev1.CreateAPIKeyRequest{
			Name:               "partner-" + req.PartnerId,
			OwnerId:            req.PartnerId,
			OwnerType:          "SERVICE",
			Scopes:             defaultPartnerAPIScopes(),
			RateLimitPerMinute: 300,
		})
		if createErr != nil {
			return nil, fmt.Errorf("create partner api key before rotate: %w", createErr)
		}
		return &partnerservicev1.RotatePartnerAPIKeyResponse{
			NewApiKey:    created.KeyId,
			NewApiSecret: created.RawKey,
			ExpiresAt:    created.ExpiresAt,
		}, nil
	}

	rotated, err := s.authnClient.RotateAPIKey(ctx, &authnservicev1.RotateAPIKeyRequest{
		KeyId:            keySummary.KeyId,
		GracePeriodHours: 24,
	})
	if err != nil {
		logger.Errorf("RotateAPIKey failed for partner=%s key_id=%s: %v", req.PartnerId, keySummary.KeyId, err)
		return nil, fmt.Errorf("rotate partner api key: %w", err)
	}

	// Resolve new key expiry from listing if available.
	var expiresAt *timestamppb.Timestamp
	keysResp, listErr := s.authnClient.ListAPIKeys(ctx, &authnservicev1.ListAPIKeysRequest{
		OwnerId:    req.PartnerId,
		OwnerType:  "SERVICE",
		ActiveOnly: true,
		PageSize:   20,
	})
	if listErr == nil {
		for _, k := range keysResp.Keys {
			if k.KeyId == rotated.NewKeyId {
				expiresAt = k.ExpiresAt
				break
			}
		}
	}
	if expiresAt == nil {
		expiresAt = rotated.OldKeyExpiresAt
	}
	s.metrics.IncAPIKeyRotated()

	return &partnerservicev1.RotatePartnerAPIKeyResponse{
		NewApiKey:    rotated.NewKeyId,
		NewApiSecret: rotated.RawKey,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *PartnerService) getFirstActivePartnerKey(ctx context.Context, partnerID string) (*authnservicev1.APIKeySummary, error) {
	keysResp, err := s.authnClient.ListAPIKeys(ctx, &authnservicev1.ListAPIKeysRequest{
		OwnerId:    partnerID,
		OwnerType:  "SERVICE",
		ActiveOnly: true,
		PageSize:   10,
	})
	if err != nil {
		logger.Errorf("ListAPIKeys failed for partner=%s: %v", partnerID, err)
		return nil, fmt.Errorf("lookup partner API keys: %w", err)
	}
	for _, k := range keysResp.Keys {
		if strings.EqualFold(k.Status, "active") || strings.EqualFold(k.Status, "api_key_status_active") || k.Status == "" {
			return k, nil
		}
	}
	if len(keysResp.Keys) > 0 {
		return keysResp.Keys[0], nil
	}
	return nil, nil
}

func defaultPartnerAPIScopes() []string {
	return []string{
		"policy:read",
		"policy:write",
		"claim:read",
		"claim:write",
		"partner:read",
	}
}
