// Package renderer provides rich document output renderers for the docgen service.
// sidecar_client.go — HTTP client for the Python docrender sidecar.
//
// The sidecar runs as a companion process (FastAPI on port 8500) and handles:
//   - DOCX generation via python-docx (real OOXML, proper tables, styles)
//   - High-quality PDF generation via WeasyPrint (CSS-accurate, no browser)
//
// Communication is plain JSON over HTTP on localhost — no auth, no TLS.
package renderer

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

// SidecarClient calls the Python docrender sidecar service.
type SidecarClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewSidecarClient creates a client for the given sidecar base URL.
// baseURL example: "http://localhost:8500"
func NewSidecarClient(baseURL string, timeout time.Duration) *SidecarClient {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &SidecarClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// IsAvailable does a quick health-check against the sidecar.
func (c *SidecarClient) IsAvailable(ctx context.Context) bool {
	url := c.baseURL + "/health"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ─── DOCX ─────────────────────────────────────────────────────────────────────

// DocxRequest is the JSON body sent to POST /render/docx.
type DocxRequest struct {
	TemplateContent string         `json:"template_content"`
	Data            map[string]any `json:"data"`
	Title           string         `json:"title"`
	Author          string         `json:"author"`
	Subject         string         `json:"subject"`
}

// RenderDOCX sends a DOCX render request to the sidecar and returns raw .docx bytes.
func (c *SidecarClient) RenderDOCX(ctx context.Context, req DocxRequest) ([]byte, error) {
	return c.post(ctx, "/render/docx", req)
}

// ─── PDF (WeasyPrint) ─────────────────────────────────────────────────────────

// PdfRequest is the JSON body sent to POST /render/pdf.
type PdfRequest struct {
	HTML    string `json:"html"`
	BaseURL string `json:"base_url,omitempty"`
}

// RenderPDF sends an HTML→PDF render request to the sidecar via WeasyPrint.
func (c *SidecarClient) RenderPDF(ctx context.Context, html string) ([]byte, error) {
	return c.post(ctx, "/render/pdf", PdfRequest{HTML: html})
}

// ─── HTTP helper ──────────────────────────────────────────────────────────────

func (c *SidecarClient) post(ctx context.Context, path string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("sidecar: failed to marshal request: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("sidecar: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sidecar: request to %s failed: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024*1024)) // 64 MB cap
	if err != nil {
		return nil, fmt.Errorf("sidecar: failed to read response from %s: %w", path, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to extract a meaningful error message from JSON detail field
		var errResp struct {
			Detail string `json:"detail"`
		}
		if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.Detail != "" {
			return nil, fmt.Errorf("sidecar: %s returned %d: %s", path, resp.StatusCode, errResp.Detail)
		}
		snippet := strings.TrimSpace(string(respBody))
		if len(snippet) > 256 {
			snippet = snippet[:256] + "..."
		}
		return nil, fmt.Errorf("sidecar: %s returned %d: %s", path, resp.StatusCode, snippet)
	}

	if len(respBody) == 0 {
		return nil, fmt.Errorf("sidecar: %s returned empty body", path)
	}

	return respBody, nil
}
