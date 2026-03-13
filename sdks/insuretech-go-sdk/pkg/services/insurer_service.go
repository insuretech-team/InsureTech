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
func (s *InsurerService) GetInsurerProduct(ctx context.Context, insurerProductId string) (*models.InsurerInsurerProductRetrievalResponse, error) {
	path := "/v1/insurer-products/{insurer_product_id}"
	path = strings.ReplaceAll(path, "{insurer_product_id}", insurerProductId)
	var result models.InsurerInsurerProductRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateInsurerProduct Update insurer product
func (s *InsurerService) UpdateInsurerProduct(ctx context.Context, insurerProductId string, req *models.InsurerInsurerProductUpdateRequest) (*models.InsurerInsurerProductUpdateResponse, error) {
	path := "/v1/insurer-products/{insurer_product_id}"
	path = strings.ReplaceAll(path, "{insurer_product_id}", insurerProductId)
	var result models.InsurerInsurerProductUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddInsurerProduct Add insurer product
func (s *InsurerService) AddInsurerProduct(ctx context.Context, insurerId string, req *models.AddInsurerProductRequest) (*models.AddInsurerProductResponse, error) {
	path := "/v1/insurers/{insurer_id}/products"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	var result models.AddInsurerProductResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListInsurerProducts List insurer products
func (s *InsurerService) ListInsurerProducts(ctx context.Context, insurerId string) (*models.InsurerInsurerProductsListingResponse, error) {
	path := "/v1/insurers/{insurer_id}/products"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	var result models.InsurerInsurerProductsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateInsurer Create insurer
func (s *InsurerService) CreateInsurer(ctx context.Context, req *models.InsurerInsurerCreationRequest) (*models.InsurerInsurerCreationResponse, error) {
	path := "/v1/insurers"
	var result models.InsurerInsurerCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListInsurers List insurers
func (s *InsurerService) ListInsurers(ctx context.Context) (*models.InsurerInsurersListingResponse, error) {
	path := "/v1/insurers"
	var result models.InsurerInsurersListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateInsurerConfig Update insurer config
func (s *InsurerService) UpdateInsurerConfig(ctx context.Context, insurerId string, req *models.InsurerInsurerConfigUpdateRequest) (*models.InsurerInsurerConfigUpdateResponse, error) {
	path := "/v1/insurers/{insurer_id}/config"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	var result models.InsurerInsurerConfigUpdateResponse
	err := s.Client.DoRequest(ctx, "PUT", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInsurer Get insurer details
func (s *InsurerService) GetInsurer(ctx context.Context, insurerId string) (*models.InsurerInsurerRetrievalResponse, error) {
	path := "/v1/insurers/{insurer_id}"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	var result models.InsurerInsurerRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateInsurer Update insurer
func (s *InsurerService) UpdateInsurer(ctx context.Context, insurerId string, req *models.InsurerInsurerUpdateRequest) (*models.InsurerInsurerUpdateResponse, error) {
	path := "/v1/insurers/{insurer_id}"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	var result models.InsurerInsurerUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

