package delivery

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WebhookConfig struct {
	Enabled   bool
	Timeout   time.Duration
	UserAgent string
}

type WebhookRequest struct {
	TargetURL    string
	Secret       string
	EventType    string
	Payload      json.RawMessage
	Timeout      time.Duration
	Subscription string
}

type WebhookResponse struct {
	StatusCode int
	Body       string
	SentAt     time.Time
}

type WebhookClient struct {
	config     WebhookConfig
	httpClient *http.Client
}

func NewWebhookClient(cfg WebhookConfig) *WebhookClient {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if strings.TrimSpace(cfg.UserAgent) == "" {
		cfg.UserAgent = "insuretech-notification-webhook/1.0"
	}
	return &WebhookClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *WebhookClient) Send(ctx context.Context, req *WebhookRequest) (*WebhookResponse, error) {
	if req == nil {
		return nil, Permanent(fmt.Errorf("webhook request is required"))
	}
	if !c.config.Enabled {
		return nil, Permanent(fmt.Errorf("webhook delivery is disabled"))
	}
	if strings.TrimSpace(req.TargetURL) == "" {
		return nil, Permanent(fmt.Errorf("webhook target URL is required"))
	}
	if strings.TrimSpace(req.Secret) == "" {
		return nil, Permanent(fmt.Errorf("webhook secret is required"))
	}
	if len(req.Payload) == 0 {
		req.Payload = json.RawMessage(`{}`)
	}
	parsedURL, err := url.Parse(req.TargetURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, Permanent(fmt.Errorf("invalid webhook target URL"))
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	signature := signWebhookPayload(req.Secret, timestamp, req.Payload)
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = c.config.Timeout
	}
	requestCtx := ctx
	var cancel context.CancelFunc
	if timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	httpReq, err := http.NewRequestWithContext(requestCtx, http.MethodPost, req.TargetURL, bytes.NewReader(req.Payload))
	if err != nil {
		return nil, Permanent(fmt.Errorf("create webhook request: %w", err))
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", c.config.UserAgent)
	httpReq.Header.Set("X-InsureTech-Event", req.EventType)
	httpReq.Header.Set("X-InsureTech-Timestamp", timestamp)
	httpReq.Header.Set("X-InsureTech-Signature", "sha256="+signature)
	if strings.TrimSpace(req.Subscription) != "" {
		httpReq.Header.Set("X-InsureTech-Subscription", req.Subscription)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send webhook request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32*1024))
	if err != nil {
		return nil, fmt.Errorf("read webhook response: %w", err)
	}
	response := &WebhookResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		SentAt:     time.Now().UTC(),
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return response, nil
	}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode >= 500 {
		return response, fmt.Errorf("webhook target returned status %d", resp.StatusCode)
	}
	return response, Permanent(fmt.Errorf("webhook target returned status %d", resp.StatusCode))
}

func signWebhookPayload(secret, timestamp string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
