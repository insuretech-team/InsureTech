package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// TenantService handles tenant-related API calls
type TenantService struct {
	Client Client
}

// GetTenantConfig Get tenant config
func (s *TenantService) GetTenantConfig(ctx context.Context, tenantId string) (*models.TenantConfigRetrievalResponse, error) {
	path := "/v1/tenants/{tenant_id}/config"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	var result models.TenantConfigRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateTenantConfig Update tenant config
func (s *TenantService) UpdateTenantConfig(ctx context.Context, tenantId string, req *models.TenantConfigUpdateRequest) (*models.TenantConfigUpdateResponse, error) {
	path := "/v1/tenants/{tenant_id}/config"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	var result models.TenantConfigUpdateResponse
	err := s.Client.DoRequest(ctx, "PUT", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTenant Get tenant
func (s *TenantService) GetTenant(ctx context.Context, tenantId string) (*models.TenantRetrievalResponse, error) {
	path := "/v1/tenants/{tenant_id}"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	var result models.TenantRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateTenant Update tenant
func (s *TenantService) UpdateTenant(ctx context.Context, tenantId string, req *models.TenantUpdateRequest) (*models.TenantUpdateResponse, error) {
	path := "/v1/tenants/{tenant_id}"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	var result models.TenantUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateTenant Create tenant
func (s *TenantService) CreateTenant(ctx context.Context, req *models.TenantCreationRequest) (*models.TenantCreationResponse, error) {
	path := "/v1/tenants"
	var result models.TenantCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListTenants List tenants
func (s *TenantService) ListTenants(ctx context.Context) (*models.TenantsListingResponse, error) {
	path := "/v1/tenants"
	var result models.TenantsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

