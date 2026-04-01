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

// ListApiKeys List API keys for owner
func (s *ApikeyService) ListApiKeys(ctx context.Context) error {
	path := "/v1/api-keys"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GenerateApiKey Generate new API key for insurer/partner
func (s *ApikeyService) GenerateApiKey(ctx context.Context, req *models.ApiKeyGenerationRequest) error {
	path := "/v1/api-keys"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetApiKey Get API key details
func (s *ApikeyService) GetApiKey(ctx context.Context, apiKeyId string) error {
	path := "/v1/api-keys/{api_key_id}"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RevokeApiKey Revoke API key
func (s *ApikeyService) RevokeApiKey(ctx context.Context, apiKeyId string) error {
	path := "/v1/api-keys/{api_key_id}"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetUsageStats Get usage statistics
func (s *ApikeyService) GetUsageStats(ctx context.Context, apiKeyId string) error {
	path := "/v1/api-keys/{api_key_id}/usage"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RotateApiKey Rotate API key
func (s *ApikeyService) RotateApiKey(ctx context.Context, apiKeyId string, req *models.APIKeyRotationRequest) error {
	path := "/v1/api-keys/{api_key_id}:rotate"
	path = strings.ReplaceAll(path, "{api_key_id}", apiKeyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ValidateApiKey Validate API key and check scopes
func (s *ApikeyService) ValidateApiKey(ctx context.Context, req *models.ApiKeyValidationRequest) error {
	path := "/v1/api-keys:validate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

