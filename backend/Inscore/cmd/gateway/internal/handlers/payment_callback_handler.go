package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentCallbackHandler struct {
	client paymentservicev1.PaymentServiceClient
}

func NewPaymentCallbackHandler(conn *grpc.ClientConn) *PaymentCallbackHandler {
	return &PaymentCallbackHandler{client: paymentservicev1.NewPaymentServiceClient(conn)}
}

func (h *PaymentCallbackHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	resp, _, err := h.forwardCallback(r, "webhook")
	if err != nil {
		logger.Warn("payment webhook processing failed", zap.Error(err))
		http.Error(w, err.Error(), callbackHTTPStatus(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":         resp.GetAccepted(),
		"payment_id": resp.GetPaymentId(),
		"status":     resp.GetStatus(),
	})
}

func (h *PaymentCallbackHandler) Success(w http.ResponseWriter, r *http.Request) {
	h.handleBrowserReturn(w, r, "success", true)
}

func (h *PaymentCallbackHandler) Fail(w http.ResponseWriter, r *http.Request) {
	h.handleBrowserReturn(w, r, "failed", false)
}

func (h *PaymentCallbackHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	h.handleBrowserReturn(w, r, "cancelled", false)
}

func (h *PaymentCallbackHandler) handleBrowserReturn(w http.ResponseWriter, r *http.Request, statusValue string, verify bool) {
	callbackType := statusValue
	if callbackType == "failed" {
		callbackType = "fail"
	}
	resp, payload, err := h.forwardCallback(r, callbackType)
	if err != nil {
		logger.Warn("payment browser callback processing failed", zap.Error(err))
		http.Error(w, err.Error(), callbackHTTPStatus(err))
		return
	}

	verified := verify && resp.GetAccepted()
	if resp.GetStatus() != "" {
		statusValue = resp.GetStatus()
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	callbackURL := h.lookupCallbackURL(ctx, firstNonEmpty(resp.GetPaymentId(), payload.PaymentID))
	if callbackURL != "" {
		redirectURL := addQueryParams(callbackURL, map[string]string{
			"payment_id":     firstNonEmpty(resp.GetPaymentId(), payload.PaymentID),
			"status":         statusValue,
			"transaction_id": firstNonEmpty(payload.TranID, payload.ValID),
			"verified":       fmt.Sprintf("%t", verified),
		})
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"payment_id":     firstNonEmpty(resp.GetPaymentId(), payload.PaymentID),
		"status":         statusValue,
		"transaction_id": firstNonEmpty(payload.TranID, payload.ValID),
		"verified":       verified,
	})
}

func (h *PaymentCallbackHandler) forwardCallback(r *http.Request, callbackType string) (*paymentservicev1.HandleGatewayWebhookResponse, *paymentCallbackPayload, error) {
	if h.client == nil {
		return nil, nil, fmt.Errorf("payment service unavailable")
	}
	payload, rawPayload, err := parsePaymentCallback(r)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	resp, err := h.client.HandleGatewayWebhook(ctx, &paymentservicev1.HandleGatewayWebhookRequest{
		Provider:   "sslcommerz",
		Headers:    callbackHeaders(r, callbackType),
		RawPayload: rawPayload,
		RemoteAddr: r.RemoteAddr,
		ReceivedAt: timestamppb.Now(),
	})
	if err != nil {
		return nil, payload, err
	}
	return resp, payload, nil
}

func (h *PaymentCallbackHandler) lookupCallbackURL(ctx context.Context, paymentID string) string {
	if strings.TrimSpace(paymentID) == "" || h.client == nil {
		return ""
	}
	resp, err := h.client.GetPayment(ctx, &paymentservicev1.GetPaymentRequest{PaymentId: paymentID})
	if err != nil || resp.GetPayment() == nil {
		return ""
	}
	return paymentGatewayField(resp.GetPayment(), "callback_url")
}

type paymentCallbackPayload struct {
	PaymentID  string
	ValID      string
	TranID     string
	SessionKey string
}

func parsePaymentCallback(r *http.Request) (*paymentCallbackPayload, []byte, error) {
	if err := r.ParseForm(); err != nil {
		return nil, nil, fmt.Errorf("invalid callback payload")
	}
	payload := &paymentCallbackPayload{
		PaymentID:  firstNonEmpty(r.FormValue("value_a"), r.FormValue("payment_id")),
		ValID:      strings.TrimSpace(r.FormValue("val_id")),
		TranID:     strings.TrimSpace(r.FormValue("tran_id")),
		SessionKey: strings.TrimSpace(r.FormValue("sessionkey")),
	}
	if firstNonEmpty(payload.PaymentID, payload.ValID, payload.TranID, payload.SessionKey) == "" {
		return nil, nil, fmt.Errorf("payment callback must include payment or provider reference")
	}
	return payload, []byte(r.Form.Encode()), nil
}

func callbackHeaders(r *http.Request, callbackType string) map[string]string {
	headers := map[string]string{
		"x-payment-callback-type": callbackType,
	}
	for key, values := range r.Header {
		if len(values) == 0 || strings.TrimSpace(values[0]) == "" {
			continue
		}
		headers[strings.ToLower(key)] = strings.TrimSpace(values[0])
	}
	return headers
}

func callbackHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if strings.Contains(msg, "invalid callback payload") || strings.Contains(msg, "payment callback must include") {
		return http.StatusBadRequest
	}
	return http.StatusBadGateway
}

func paymentGatewayField(payment *paymententityv1.Payment, key string) string {
	if payment == nil || strings.TrimSpace(payment.GetGatewayResponse()) == "" {
		return ""
	}
	values := map[string]string{}
	if err := json.Unmarshal([]byte(payment.GetGatewayResponse()), &values); err != nil {
		return ""
	}
	return values[key]
}

func addQueryParams(raw string, params map[string]string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	query := parsed.Query()
	for key, value := range params {
		if strings.TrimSpace(value) != "" {
			query.Set(key, value)
		}
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func writeJSON(w http.ResponseWriter, status int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
