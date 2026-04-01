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

// ListTenants List tenants
func (s *TenantService) ListTenants(ctx context.Context) error {
	path := "/v1/tenants"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateTenant Create tenant
func (s *TenantService) CreateTenant(ctx context.Context, req *models.TenantCreationRequest) error {
	path := "/v1/tenants"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetTenant Get tenant
func (s *TenantService) GetTenant(ctx context.Context, tenantId string) error {
	path := "/v1/tenants/{tenant_id}"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateTenant Update tenant
func (s *TenantService) UpdateTenant(ctx context.Context, tenantId string, req *models.TenantUpdateRequest) error {
	path := "/v1/tenants/{tenant_id}"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// GetTenantConfig Get tenant config
func (s *TenantService) GetTenantConfig(ctx context.Context, tenantId string) error {
	path := "/v1/tenants/{tenant_id}/config"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateTenantConfig Update tenant config
func (s *TenantService) UpdateTenantConfig(ctx context.Context, tenantId string, req *models.TenantConfigUpdateRequest) error {
	path := "/v1/tenants/{tenant_id}/config"
	path = strings.ReplaceAll(path, "{tenant_id}", tenantId)
	return s.Client.DoRequest(ctx, "PUT", path, req, nil)
}

