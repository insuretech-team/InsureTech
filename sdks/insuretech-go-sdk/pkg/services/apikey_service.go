package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// ApikeyService handles apikey-related API calls
type ApikeyService struct {
	Client Client
}

// GenerateApiKey Generate new API key for insurer/partner
func (s *ApikeyService) GenerateApiKey(ctx context.Context, req *models.ApiKeyGenerationRequest) (*models.ApiKeyGenerationResponse, error) {
	path := "/v1/api-keys"
	var result models.ApiKeyGenerationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListApiKeys List API keys for owner
func (s *ApikeyService) ListApiKeys(ctx context.Context) (*models.APIKeysListingResponse, error) {
	path := "/v1/api-keys"
	var result models.APIKeysListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetApiKey Get API key details
func (s *ApikeyService) GetApiKey(ctx context.Context, apiKeyId string) (*models.ApiKeyRetrievalResponse, error) {
	path := "/v1/api-keys/{api_key_id}"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	var result models.ApiKeyRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeApiKey Revoke API key
func (s *ApikeyService) RevokeApiKey(ctx context.Context, apiKeyId string) error {
	path := "/v1/api-keys/{api_key_id}"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// RotateApiKey Rotate API key
func (s *ApikeyService) RotateApiKey(ctx context.Context, apiKeyId string, req *models.APIKeyRotationRequest) (*models.APIKeyRotationResponse, error) {
	path := "/v1/api-keys/{api_key_id}"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	var result models.APIKeyRotationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUsageStats Get usage statistics
func (s *ApikeyService) GetUsageStats(ctx context.Context, apiKeyId string) (*models.UsageStatsRetrievalResponse, error) {
	path := "/v1/api-keys/{api_key_id}/usage"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	var result models.UsageStatsRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

