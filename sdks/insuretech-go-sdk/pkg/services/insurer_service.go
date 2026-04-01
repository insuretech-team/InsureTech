package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// InsurerService handles insurer-related API calls
type InsurerService struct {
	Client Client
}

// GetInsurerProduct Get insurer product
func (s *InsurerService) GetInsurerProduct(ctx context.Context, insurerProductId string) error {
	path := "/v1/insurer-products/{insurer_product_id}"
	path = strings.ReplaceAll(path, "{insurer_product_id}", insurerProductId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateInsurerProduct Update insurer product
func (s *InsurerService) UpdateInsurerProduct(ctx context.Context, insurerProductId string, req *models.InsurerInsurerProductUpdateRequest) error {
	path := "/v1/insurer-products/{insurer_product_id}"
	path = strings.ReplaceAll(path, "{insurer_product_id}", insurerProductId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// ListInsurers List insurers
func (s *InsurerService) ListInsurers(ctx context.Context) error {
	path := "/v1/insurers"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateInsurer Create insurer
func (s *InsurerService) CreateInsurer(ctx context.Context, req *models.InsurerInsurerCreationRequest) error {
	path := "/v1/insurers"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetInsurer Get insurer details
func (s *InsurerService) GetInsurer(ctx context.Context, insurerId string) error {
	path := "/v1/insurers/{insurer_id}"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateInsurer Update insurer
func (s *InsurerService) UpdateInsurer(ctx context.Context, insurerId string, req *models.InsurerInsurerUpdateRequest) error {
	path := "/v1/insurers/{insurer_id}"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// UpdateInsurerConfig Update insurer config
func (s *InsurerService) UpdateInsurerConfig(ctx context.Context, insurerId string, req *models.InsurerInsurerConfigUpdateRequest) error {
	path := "/v1/insurers/{insurer_id}/config"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	return s.Client.DoRequest(ctx, "PUT", path, req, nil)
}

// ListInsurerProducts List insurer products
func (s *InsurerService) ListInsurerProducts(ctx context.Context, insurerId string) error {
	path := "/v1/insurers/{insurer_id}/products"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// AddInsurerProduct Add insurer product
func (s *InsurerService) AddInsurerProduct(ctx context.Context, insurerId string, req *models.AddInsurerProductRequest) error {
	path := "/v1/insurers/{insurer_id}/products"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

