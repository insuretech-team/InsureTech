package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// SMTP2GoClient sends email via the smtp2go HTTP API (port 443 — not blocked by DigitalOcean).
// Falls back to the standard SMTP client if SMTP2GO_API_KEY is not set.
type SMTP2GoClient struct {
	apiKey  string
	apiURL  string
	from    string
	httpCli *http.Client
}

type smtp2goRequest struct {
	APIKey   string   `json:"api_key"`
	To       []string `json:"to"`
	Sender   string   `json:"sender"`
	Subject  string   `json:"subject"`
	HTMLBody string   `json:"html_body,omitempty"`
	TextBody string   `json:"text_body,omitempty"`
}

type smtp2goResponse struct {
	Data struct {
		Error     string `json:"error"`
		ErrorCode int    `json:"error_code"`
		EmailID   string `json:"email_id"`
	} `json:"data"`
}

// NewSMTP2GoClient creates an smtp2go API email client.
func NewSMTP2GoClient(from string) *SMTP2GoClient {
	apiKey := os.Getenv("SMTP2GO_API_KEY")
	apiURL := os.Getenv("SMTP2GO_BASE_URL")
	if apiURL == "" {
		apiURL = "https://api.smtp2go.com/v3"
	}
	return &SMTP2GoClient{
		apiKey:  apiKey,
		apiURL:  apiURL,
		from:    from,
		httpCli: &http.Client{Timeout: 30 * time.Second},
	}
}

// IsConfigured returns true if the smtp2go API key is set.
func (c *SMTP2GoClient) IsConfigured() bool {
	return c.apiKey != "" && c.apiKey != "SMTP2GO_API_KEY_HERE"
}

// Send sends an email via the smtp2go HTTP API.
func (c *SMTP2GoClient) Send(ctx context.Context, to, subject, htmlBody, textBody string) error {
	if !c.IsConfigured() {
		return fmt.Errorf("smtp2go: API key not configured")
	}

	payload := smtp2goRequest{
		APIKey:   c.apiKey,
		To:       []string{to},
		Sender:   c.from,
		Subject:  subject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("smtp2go: marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+"/email/send", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("smtp2go: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("smtp2go: HTTP error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("smtp2go: API error %d: %s", resp.StatusCode, respBody)
	}

	var s2gResp smtp2goResponse
	if err := json.Unmarshal(respBody, &s2gResp); err == nil {
		if s2gResp.Data.ErrorCode != 0 {
			return fmt.Errorf("smtp2go: send error: %s (code %d)", s2gResp.Data.Error, s2gResp.Data.ErrorCode)
		}
	}

	return nil
}
