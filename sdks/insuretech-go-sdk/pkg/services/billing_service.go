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

// ListInvoices List invoices with optional filters
func (s *BillingService) ListInvoices(ctx context.Context) error {
	path := "/v1/invoices"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateInvoice Create a new invoice for an order (B2C) or purchase order (B2B)
func (s *BillingService) CreateInvoice(ctx context.Context, req *models.InvoiceCreationRequest) error {
	path := "/v1/invoices"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetInvoice Get a single invoice by ID
func (s *BillingService) GetInvoice(ctx context.Context, invoiceId string) error {
	path := "/v1/invoices/{invoice_id}"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetInvoicePDF Get invoice PDF (pre-signed download URL or file ID)
func (s *BillingService) GetInvoicePDF(ctx context.Context, invoiceId string) error {
	path := "/v1/invoices/{invoice_id}/pdf"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CancelInvoice Cancel an invoice (only allowed before PAID)
func (s *BillingService) CancelInvoice(ctx context.Context, invoiceId string, req *models.InvoiceCancellationRequest) error {
	path := "/v1/invoices/{invoice_id}:cancel"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GenerateInvoicePDF Trigger async invoice PDF generation
func (s *BillingService) GenerateInvoicePDF(ctx context.Context, invoiceId string, req *models.InvoicePDFGenerationRequest) error {
	path := "/v1/invoices/{invoice_id}:generate-pdf"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// IssueInvoice Issue an invoice — transitions from DRAFT → ISSUED and sends to customer/org
func (s *BillingService) IssueInvoice(ctx context.Context, invoiceId string, req *models.InvoiceIssuanceRequest) error {
	path := "/v1/invoices/{invoice_id}:issue"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// MarkInvoicePaid Mark invoice as paid — called by payment-service after payment confirmation
func (s *BillingService) MarkInvoicePaid(ctx context.Context, invoiceId string, req *models.MarkInvoicePaidRequest) error {
	path := "/v1/invoices/{invoice_id}:mark-paid"
	path = strings.ReplaceAll(path, "{invoice_id}", invoiceId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetInvoiceByOrderId Get invoice by order ID — used by orders-service to link invoice after creation
func (s *BillingService) GetInvoiceByOrderId(ctx context.Context, orderId string) error {
	path := "/v1/orders/{order_id}/invoice"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

