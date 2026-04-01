package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// ProductService handles product-related API calls
type ProductService struct {
	Client Client
}

// ListProducts List all active products
func (s *ProductService) ListProducts(ctx context.Context) error {
	path := "/v1/products"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateProduct Create product (admin)
func (s *ProductService) CreateProduct(ctx context.Context, req *models.ProductsProductCreationRequest) error {
	path := "/v1/products"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetProduct Get product details
func (s *ProductService) GetProduct(ctx context.Context, productId string) error {
	path := "/v1/products/{product_id}"
	path = strings.ReplaceAll(path, "{product_id}", productId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateProduct Update product (admin)
func (s *ProductService) UpdateProduct(ctx context.Context, productId string, req *models.ProductsProductUpdateRequest) error {
	path := "/v1/products/{product_id}"
	path = strings.ReplaceAll(path, "{product_id}", productId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// ActivateProduct Activate product
func (s *ProductService) ActivateProduct(ctx context.Context, productId string, req *models.ProductActivationRequest) error {
	path := "/v1/products/{product_id}:activate"
	path = strings.ReplaceAll(path, "{product_id}", productId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// CalculatePremium Calculate premium
func (s *ProductService) CalculatePremium(ctx context.Context, productId string, req *models.ProductsPremiumCalculationRequest) error {
	path := "/v1/products/{product_id}:calculate-premium"
	path = strings.ReplaceAll(path, "{product_id}", productId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DeactivateProduct Deactivate product
func (s *ProductService) DeactivateProduct(ctx context.Context, productId string, req *models.ProductDeactivationRequest) error {
	path := "/v1/products/{product_id}:deactivate"
	path = strings.ReplaceAll(path, "{product_id}", productId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DiscontinueProduct Discontinue product
func (s *ProductService) DiscontinueProduct(ctx context.Context, productId string, req *models.DiscontinueProductRequest) error {
	path := "/v1/products/{product_id}:discontinue"
	path = strings.ReplaceAll(path, "{product_id}", productId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SearchProducts Search products
func (s *ProductService) SearchProducts(ctx context.Context) error {
	path := "/v1/products:search"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

