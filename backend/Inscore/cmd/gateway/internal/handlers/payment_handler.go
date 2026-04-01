package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	paymentv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
)

// PaymentHandler proxies payment CRUD requests to the payment gRPC service.
// BUG-006 FIX: Replaces PoliSyncHandler (HTTP proxy) which cannot reach the
// gRPC-only payment service (no HTTP companion server on :50191).
type PaymentHandler struct {
	client paymentv1.PaymentServiceClient
}

// NewPaymentHandler creates a PaymentHandler from a gRPC connection to the payment service.
func NewPaymentHandler(conn *grpc.ClientConn) *PaymentHandler {
	return &PaymentHandler{client: paymentv1.NewPaymentServiceClient(conn)}
}

// List handles GET /v1/payments
func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListPayments(ctx, &paymentv1.ListPaymentsRequest{
			UserId:   q.Get("user_id"),
			PolicyId: q.Get("policy_id"),
		})
	})
}

// Get handles GET /v1/payments/{payment_id}
func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPayment(ctx, &paymentv1.GetPaymentRequest{PaymentId: paymentID})
	})
}

// Initiate handles POST /v1/payments
// BUG FIX: Clients may send flat "amount": 10000 (number) + "currency": "BDT"
// but the proto expects "amount": {"amount": 10000, "currency": "BDT"} (Money object).
// Normalize the body before proto unmarshalling.
func (h *PaymentHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		body = normalizePaymentBody(body)
		var req paymentv1.InitiatePaymentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.InitiatePayment(ctx, &req)
	})
}

// normalizePaymentBody converts flat amount/currency fields to a Money object.
// Accepts: {"amount": 10000, "currency": "BDT", "user_id": "...", ...}
// Outputs: {"amount": {"amount": 10000, "currency": "BDT"}, "user_id": "...", ...}
// Also maps client-friendly field names to proto field names.
func normalizePaymentBody(body []byte) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return body
	}
	changed := false

	// If "amount" is a number (not already an object), wrap it in Money format
	if amtRaw, ok := m["amount"]; ok {
		switch amt := amtRaw.(type) {
		case float64:
			// Flat number: {"amount": 10000} → {"amount": {"amount": 10000, "currency": "BDT"}}
			currency := "BDT"
			if c, ok := m["currency"].(string); ok && c != "" {
				currency = c
				delete(m, "currency") // remove top-level currency
			}
			m["amount"] = map[string]interface{}{
				"amount":   int64(amt),
				"currency": currency,
			}
			changed = true
		case json.Number:
			v, _ := amt.Float64()
			currency := "BDT"
			if c, ok := m["currency"].(string); ok && c != "" {
				currency = c
				delete(m, "currency")
			}
			m["amount"] = map[string]interface{}{
				"amount":   int64(v),
				"currency": currency,
			}
			changed = true
		}
	}

	// Map "provider" → "payment_method" if payment_method not set
	if provider, ok := m["provider"]; ok {
		if _, hasMethod := m["payment_method"]; !hasMethod {
			if provStr, ok := provider.(string); ok {
				// Map provider to payment_method enum
				methodMap := map[string]string{
					"sslcommerz": "PAYMENT_METHOD_CARD",
					"card":       "PAYMENT_METHOD_CARD",
					"mobile":     "PAYMENT_METHOD_MOBILE_BANKING",
					"bank":       "PAYMENT_METHOD_BANK_TRANSFER",
					"cash":       "PAYMENT_METHOD_CASH",
				}
				if method, found := methodMap[provStr]; found {
					m["payment_method"] = method
					changed = true
				}
			}
		}
	}

	// Map return_url/cancel_url → callback_url
	if returnURL, ok := m["return_url"]; ok {
		if _, hasCB := m["callback_url"]; !hasCB {
			m["callback_url"] = returnURL
			changed = true
		}
		delete(m, "return_url")
		delete(m, "cancel_url")
	}

	// Auto-generate idempotency_key if not provided (required by payment service)
	if _, hasKey := m["idempotency_key"]; !hasKey {
		// Use a deterministic key from user_id + timestamp
		userID, _ := m["user_id"].(string)
		m["idempotency_key"] = fmt.Sprintf("gw-%s-%d", userID, time.Now().UnixMilli())
		changed = true
	}

	if !changed {
		return body
	}
	out, err := json.Marshal(m)
	if err != nil {
		return body
	}
	return out
}


// Verify handles POST /v1/payments/{payment_id}/verify
func (h *PaymentHandler) Verify(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req paymentv1.VerifyPaymentRequest
		_ = protoUnmarshal(body, &req)
		req.PaymentId = paymentID
		return h.client.VerifyPayment(ctx, &req)
	})
}

// InitiateRefund handles POST /v1/payments/{payment_id}/refunds
func (h *PaymentHandler) InitiateRefund(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req paymentv1.InitiateRefundRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PaymentId = paymentID
		return h.client.InitiateRefund(ctx, &req)
	})
}

// GetRefund handles GET /v1/refunds/{refund_id}
func (h *PaymentHandler) GetRefund(w http.ResponseWriter, r *http.Request) {
	refundID := r.PathValue("refund_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetRefundStatus(ctx, &paymentv1.GetRefundStatusRequest{RefundId: refundID})
	})
}

// ListMethods handles GET /v1/users/{user_id}/payment-methods
func (h *PaymentHandler) ListMethods(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListPaymentMethods(ctx, &paymentv1.ListPaymentMethodsRequest{UserId: userID})
	})
}

// AddMethod handles POST /v1/users/{user_id}/payment-methods
func (h *PaymentHandler) AddMethod(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req paymentv1.AddPaymentMethodRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		req.UserId = userID
		return h.client.AddPaymentMethod(ctx, &req)
	})
}

// SubmitProof handles POST /v1/payments/{payment_id}/submit-proof
func (h *PaymentHandler) SubmitProof(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req paymentv1.SubmitManualPaymentProofRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PaymentId = paymentID
		return h.client.SubmitManualPaymentProof(ctx, &req)
	})
}

// Review handles POST /v1/payments/{payment_id}/review
func (h *PaymentHandler) Review(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req paymentv1.ReviewManualPaymentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		req.PaymentId = paymentID
		return h.client.ReviewManualPayment(ctx, &req)
	})
}

// GenerateReceipt handles POST /v1/payments/{payment_id}/generate-receipt
func (h *PaymentHandler) GenerateReceipt(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GenerateReceipt(ctx, &paymentv1.GenerateReceiptRequest{PaymentId: paymentID})
	})
}

// GetReceipt handles GET /v1/payments/{payment_id}/receipt
func (h *PaymentHandler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("payment_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPaymentReceipt(ctx, &paymentv1.GetPaymentReceiptRequest{PaymentId: paymentID})
	})
}

// GetByProviderRef handles GET /v1/payments/provider/{provider}/references/{provider_reference}
func (h *PaymentHandler) GetByProviderRef(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	ref := r.PathValue("provider_reference")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPaymentByProviderReference(ctx, &paymentv1.GetPaymentByProviderReferenceRequest{
			Provider:          provider,
			ProviderReference: ref,
		})
	})
}

// Reconcile handles POST /v1/payments/reconcile
func (h *PaymentHandler) Reconcile(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req paymentv1.ReconcilePaymentsRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ReconcilePayments(ctx, &req)
	})
}

// parseI32 safely converts a query string to int32 with a default value.
func parseI32(s string, defaultVal int32) int32 {
	if s == "" {
		return defaultVal
	}
	var v int32
	if _, err := fmt.Sscanf(s, "%d", &v); err != nil {
		return defaultVal
	}
	return v
}
