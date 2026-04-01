package kyc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ── FLVEAdapter Interface ────────────────────────────────────────────────────

// FLVEAdapter is the interface for FLVE eKYC session management.
type FLVEAdapter interface {
	StartEKYC(ctx context.Context, req *FLVEStartRequest) (*FLVEStartResponse, error)
}

// ── Request / Response Types ─────────────────────────────────────────────────

type FLVEStartRequest struct {
	UserID            string `json:"user_id"`
	TenantID          string `json:"tenant_id,omitempty"`
	UserType          string `json:"user_type,omitempty"`
	Portal            string `json:"portal,omitempty"`
	KYCVerificationID string `json:"kyc_verification_id,omitempty"`
}

type FLVEStartResponse struct {
	SessionID           string `json:"session_id"`
	State               string `json:"state"`
	TotalTimeoutSeconds int    `json:"total_timeout_seconds"`
	Error               string `json:"error,omitempty"`
}

// ── HTTP Adapter Implementation ──────────────────────────────────────────────

type flveHTTPAdapter struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewFLVEAdapter creates an FLVEAdapter backed by FLVE HTTP endpoints.
func NewFLVEAdapter(baseURL, token string, timeout time.Duration) FLVEAdapter {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &flveHTTPAdapter{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (a *flveHTTPAdapter) StartEKYC(ctx context.Context, req *FLVEStartRequest) (*FLVEStartResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/ekyc/start", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(a.token) != "" {
		httpReq.Header.Set("Authorization", "Bearer "+a.token)
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("flve request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check for HTML error responses from HuggingFace
	if resp.StatusCode != 200 && strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		switch resp.StatusCode {
		case 401, 403:
			return nil, fmt.Errorf("FLVE authentication failed (HTTP %d): check token", resp.StatusCode)
		case 404:
			return nil, fmt.Errorf("FLVE endpoint not found (HTTP 404): Space may be sleeping or token is wrong")
		default:
			return nil, fmt.Errorf("FLVE returned HTML error (HTTP %d)", resp.StatusCode)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("flve start ekyc status %d: %s", resp.StatusCode, string(respBody))
	}

	var result FLVEStartResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("flve start ekyc error: %s", result.Error)
	}
	return &result, nil
}
