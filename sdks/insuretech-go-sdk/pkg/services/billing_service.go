package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// BillingService handles billing-related API calls
type BillingService struct {
	Client Client
}

// GetInvoiceByOrderId Get invoice by order ID — used by orders-service to link invoice after creation
func (s *BillingService) GetInvoiceByOrderId(ctx context.Context, orderId string) (*models.InvoiceByOrderIdRetrievalResponse, error) {
	path := "/v1/orders/{order_id}/invoice"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	var result models.InvoiceByOrderIdRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GenerateInvoicePDF Trigger async invoice PDF generation
func (s *BillingService) GenerateInvoicePDF(ctx context.Context, invoiceId string, req *models.InvoicePDFGenerationRequest) (*models.InvoicePDFGenerationResponse, error) {
	path := "/v1/invoices/{invoice_id}:generate-pdf"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	var result models.InvoicePDFGenerationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateInvoice Create a new invoice for an order (B2C) or purchase order (B2B)
func (s *BillingService) CreateInvoice(ctx context.Context, req *models.InvoiceCreationRequest) (*models.InvoiceCreationResponse, error) {
	path := "/v1/invoices"
	var result models.InvoiceCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListInvoices List invoices with optional filters
func (s *BillingService) ListInvoices(ctx context.Context) (*models.InvoicesListingResponse, error) {
	path := "/v1/invoices"
	var result models.InvoicesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IssueInvoice Issue an invoice — transitions from DRAFT → ISSUED and sends to customer/org
func (s *BillingService) IssueInvoice(ctx context.Context, invoiceId string, req *models.InvoiceIssuanceRequest) (*models.InvoiceIssuanceResponse, error) {
	path := "/v1/invoices/{invoice_id}:issue"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	var result models.InvoiceIssuanceResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelInvoice Cancel an invoice (only allowed before PAID)
func (s *BillingService) CancelInvoice(ctx context.Context, invoiceId string, req *models.InvoiceCancellationRequest) (*models.InvoiceCancellationResponse, error) {
	path := "/v1/invoices/{invoice_id}:cancel"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	var result models.InvoiceCancellationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInvoicePDF Get invoice PDF (pre-signed download URL or file ID)
func (s *BillingService) GetInvoicePDF(ctx context.Context, invoiceId string) (*models.InvoicePDFRetrievalResponse, error) {
	path := "/v1/invoices/{invoice_id}/pdf"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	var result models.InvoicePDFRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInvoice Get a single invoice by ID
func (s *BillingService) GetInvoice(ctx context.Context, invoiceId string) (*models.InvoiceRetrievalResponse, error) {
	path := "/v1/invoices/{invoice_id}"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	var result models.InvoiceRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarkInvoicePaid Mark invoice as paid — called by payment-service after payment confirmation
func (s *BillingService) MarkInvoicePaid(ctx context.Context, invoiceId string, req *models.MarkInvoicePaidRequest) (*models.MarkInvoicePaidResponse, error) {
	path := "/v1/invoices/{invoice_id}:mark-paid"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	var result models.MarkInvoicePaidResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

