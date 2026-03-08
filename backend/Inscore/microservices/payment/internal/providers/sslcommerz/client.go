package sslcommerz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	paymentcfg "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/domain"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type Client struct {
	httpClient        *http.Client
	storeID           string
	storePassword     string
	apiBaseURL        string
	validationBaseURL string
	refundBaseURL     string
}

func NewClient(cfg *paymentcfg.Config) *Client {
	timeout := cfg.HTTPTimeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Client{
		httpClient:        &http.Client{Timeout: timeout},
		storeID:           cfg.SSLCommerzStoreID,
		storePassword:     cfg.SSLCommerzStorePassword,
		apiBaseURL:        strings.TrimRight(cfg.SSLCommerzAPIBaseURL, "/"),
		validationBaseURL: strings.TrimRight(cfg.SSLCommerzValidationBaseURL, "/"),
		refundBaseURL:     strings.TrimRight(cfg.SSLCommerzRefundBaseURL, "/"),
	}
}

func (c *Client) InitSession(ctx context.Context, req *domain.GatewaySessionRequest) (*domain.GatewaySessionResponse, error) {
	values := url.Values{
		"store_id":         {c.storeID},
		"store_passwd":     {c.storePassword},
		"total_amount":     {formatMoney(req.Amount)},
		"currency":         {firstNonEmpty(req.Currency, "BDT")},
		"tran_id":          {req.TransactionID},
		"success_url":      {req.SuccessURL},
		"fail_url":         {req.FailURL},
		"cancel_url":       {req.CancelURL},
		"ipn_url":          {req.IPNURL},
		"cus_name":         {req.CustomerName},
		"cus_email":        {req.CustomerEmail},
		"cus_phone":        {req.CustomerPhone},
		"cus_add1":         {req.CustomerAddr1},
		"cus_city":         {req.CustomerCity},
		"cus_postcode":     {req.CustomerPostcode},
		"cus_country":      {req.CustomerCountry},
		"shipping_method":  {"NO"},
		"product_name":     {firstNonEmpty(req.Metadata["product_name"], "Insurance Premium")},
		"product_category": {firstNonEmpty(req.Metadata["product_category"], "Insurance")},
		"product_profile":  {firstNonEmpty(req.Metadata["product_profile"], "general")},
		"value_a":          {req.PaymentID},
		"value_b":          {req.OrderID},
		"value_c":          {req.TenantID},
		"value_d":          {req.Metadata["correlation_id"]},
	}

	payload, err := c.postForm(ctx, c.apiBaseURL+"/gwprocess/v4/api.php", values)
	if err != nil {
		return nil, err
	}
	if !statusLooksSuccessful(payload["status"]) {
		return nil, fmt.Errorf("sslcommerz init failed: %s", firstNonEmpty(payload["failedreason"], payload["failed_reason"], payload["status"]))
	}
	return &domain.GatewaySessionResponse{
		Provider:       "SSLCOMMERZ",
		Status:         payload["status"],
		GatewayPageURL: firstNonEmpty(payload["GatewayPageURL"], payload["gatewayPageURL"]),
		SessionKey:     firstNonEmpty(payload["sessionkey"], payload["sessionKey"]),
		TranID:         req.TransactionID,
		RawFields:      payload,
	}, nil
}

func (c *Client) ValidatePayment(ctx context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
	if strings.TrimSpace(req.TransactionID) == "" {
		return nil, fmt.Errorf("sslcommerz validation requires transaction identifier")
	}
	values := url.Values{
		"val_id":       {req.TransactionID},
		"store_id":     {c.storeID},
		"store_passwd": {c.storePassword},
		"format":       {"json"},
	}
	payload, err := c.postForm(ctx, c.validationBaseURL+"/validator/api/validationserverAPI.php", values)
	if err != nil {
		return nil, err
	}
	return validationFromPayload(payload), nil
}

func (c *Client) QueryPayment(ctx context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
	values := url.Values{
		"store_id":     {c.storeID},
		"store_passwd": {c.storePassword},
		"format":       {"json"},
	}
	if strings.TrimSpace(req.TransactionID) != "" {
		values.Set("tran_id", req.TransactionID)
	}
	if strings.TrimSpace(req.SessionKey) != "" {
		values.Set("sessionkey", req.SessionKey)
	}
	payload, err := c.postForm(ctx, c.validationBaseURL+"/validator/api/merchantTransIDvalidationAPI.php", values)
	if err != nil {
		return nil, err
	}
	return validationFromPayload(payload), nil
}

func (c *Client) InitiateRefund(ctx context.Context, req *domain.GatewayRefundRequest) (*domain.GatewayRefundResponse, error) {
	values := url.Values{
		"bank_tran_id":   {req.BankTransactionID},
		"refund_amount":  {formatMoney(req.Amount)},
		"refund_remarks": {firstNonEmpty(req.Reason, "Refund initiated")},
		"store_id":       {c.storeID},
		"store_passwd":   {c.storePassword},
		"format":         {"json"},
	}
	payload, err := c.postForm(ctx, c.refundBaseURL+"/validator/api/merchantTransIDvalidationAPI.php", values)
	if err != nil {
		return nil, err
	}
	return &domain.GatewayRefundResponse{
		Provider:    "SSLCOMMERZ",
		Status:      normalizeProviderStatus(payload["status"]),
		RefundRefID: payload["refund_ref_id"],
		RawFields:   payload,
	}, nil
}

func (c *Client) postForm(ctx context.Context, endpoint string, values url.Values) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sslcommerz request failed with status %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return stringifyMap(payload), nil
}

func stringifyMap(in map[string]any) map[string]string {
	out := make(map[string]string, len(in))
	for key, value := range in {
		switch v := value.(type) {
		case string:
			out[key] = v
		case float64:
			out[key] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			out[key] = strconv.FormatBool(v)
		default:
			raw, _ := json.Marshal(v)
			out[key] = string(raw)
		}
	}
	return out
}

func validationFromPayload(payload map[string]string) *domain.GatewayValidationResponse {
	return &domain.GatewayValidationResponse{
		Provider:          "SSLCOMMERZ",
		Status:            normalizeProviderStatus(payload["status"]),
		TransactionID:     firstNonEmpty(payload["tran_id"], payload["transaction_id"]),
		ValidationID:      firstNonEmpty(payload["val_id"], payload["validation_id"]),
		BankTransactionID: payload["bank_tran_id"],
		Amount:            parseMoney(payload["amount"], payload["currency"]),
		CardType:          payload["card_type"],
		CardBrand:         payload["card_brand"],
		CardIssuer:        payload["card_issuer"],
		CardIssuerCountry: payload["card_issuer_country"],
		RiskLevel:         payload["risk_level"],
		RiskTitle:         payload["risk_title"],
		ValidatedAt:       time.Now().UTC(),
		RawFields:         payload,
	}
}

func parseMoney(amount, currency string) *commonv1.Money {
	if strings.TrimSpace(amount) == "" {
		return nil
	}
	decimalAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return &commonv1.Money{Currency: firstNonEmpty(currency, "BDT")}
	}
	return &commonv1.Money{
		Amount:        int64(decimalAmount * 100),
		DecimalAmount: decimalAmount,
		Currency:      firstNonEmpty(currency, "BDT"),
	}
}

func formatMoney(value *commonv1.Money) string {
	if value == nil {
		return "0.00"
	}
	if value.DecimalAmount > 0 {
		return strconv.FormatFloat(value.DecimalAmount, 'f', 2, 64)
	}
	return strconv.FormatFloat(float64(value.Amount)/100, 'f', 2, 64)
}

func statusLooksSuccessful(value string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	switch normalized {
	case "SUCCESS", "VALID", "VALIDATED", "INITIATED":
		return true
	default:
		return false
	}
}

func normalizeProviderStatus(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
