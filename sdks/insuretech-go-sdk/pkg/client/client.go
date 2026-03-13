package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/newage-saint/insuretech-go-sdk/pkg/services"
)

// Client is the main InsureTech API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	
	// Service clients - auto-generated from OpenAPI spec
	// These will be populated during generation
		Audit *services.AuditService
	Fraud *services.FraudService
	Payment *services.PaymentService
	Product *services.ProductService
	Ai *services.AiService
	Beneficiary *services.BeneficiaryService
	Apikey *services.ApikeyService
	Mfs *services.MfsService
	Authz *services.AuthzService
	Renewal *services.RenewalService
	Task *services.TaskService
	Insurer *services.InsurerService
	Claim *services.ClaimService
	Refund *services.RefundService
	Report *services.ReportService
	Auth *services.AuthService
	Kyc *services.KycService
	Tenant *services.TenantService
	Support *services.SupportService
	Document *services.DocumentService
	Endorsement *services.EndorsementService
	Workflow *services.WorkflowService
	Partner *services.PartnerService
	Notification *services.NotificationService
	Billing *services.BillingService
	Analytics *services.AnalyticsService
	B2b *services.B2bService
	Underwriting *services.UnderwritingService
	Order *services.OrderService
	Voice *services.VoiceService
	Media *services.MediaService
	Iot *services.IotService
	Policy *services.PolicyService
	Commission *services.CommissionService

}

// ClientOption is a function that configures the Client
type ClientOption func(*Client)

// NewClient creates a new InsureTech API client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: "https://b2b.labaidinsuretech.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	
	// Apply options
	for _, opt := range opts {
		opt(c)
	}
	
	// Initialize service clients - auto-generated
		c.Audit = &services.AuditService{Client: c}
	c.Fraud = &services.FraudService{Client: c}
	c.Payment = &services.PaymentService{Client: c}
	c.Product = &services.ProductService{Client: c}
	c.Ai = &services.AiService{Client: c}
	c.Beneficiary = &services.BeneficiaryService{Client: c}
	c.Apikey = &services.ApikeyService{Client: c}
	c.Mfs = &services.MfsService{Client: c}
	c.Authz = &services.AuthzService{Client: c}
	c.Renewal = &services.RenewalService{Client: c}
	c.Task = &services.TaskService{Client: c}
	c.Insurer = &services.InsurerService{Client: c}
	c.Claim = &services.ClaimService{Client: c}
	c.Refund = &services.RefundService{Client: c}
	c.Report = &services.ReportService{Client: c}
	c.Auth = &services.AuthService{Client: c}
	c.Kyc = &services.KycService{Client: c}
	c.Tenant = &services.TenantService{Client: c}
	c.Support = &services.SupportService{Client: c}
	c.Document = &services.DocumentService{Client: c}
	c.Endorsement = &services.EndorsementService{Client: c}
	c.Workflow = &services.WorkflowService{Client: c}
	c.Partner = &services.PartnerService{Client: c}
	c.Notification = &services.NotificationService{Client: c}
	c.Billing = &services.BillingService{Client: c}
	c.Analytics = &services.AnalyticsService{Client: c}
	c.B2b = &services.B2bService{Client: c}
	c.Underwriting = &services.UnderwritingService{Client: c}
	c.Order = &services.OrderService{Client: c}
	c.Voice = &services.VoiceService{Client: c}
	c.Media = &services.MediaService{Client: c}
	c.Iot = &services.IotService{Client: c}
	c.Policy = &services.PolicyService{Client: c}
	c.Commission = &services.CommissionService{Client: c}

	
	return c
}

// WithAPIKey sets the API key for authentication
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithBaseURL sets the base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithRetry configures retry behavior
func WithRetry(maxRetries int, retryDelay time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Transport = &retryTransport{
			transport:  http.DefaultTransport,
			maxRetries: maxRetries,
			retryDelay: retryDelay,
		}
	}
}

// WithOAuth2 sets OAuth2 bearer token authentication
func WithOAuth2(token string) ClientOption {
	return func(c *Client) {
		c.apiKey = token
	}
}

// DoRequest performs an HTTP request with authentication and error handling
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Build URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	
	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = strings.NewReader(string(bodyBytes))
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	
	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	
	// Parse response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}
	
	return nil
}

// retryTransport implements retry logic for HTTP requests
type retryTransport struct {
	transport  http.RoundTripper
	maxRetries int
	retryDelay time.Duration
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	
	for i := 0; i <= t.maxRetries; i++ {
		resp, err = t.transport.RoundTrip(req)
		
		// Success or non-retryable error
		if err == nil && (resp.StatusCode < 500 || resp.StatusCode == 501) {
			return resp, nil
		}
		
		// Don't retry on last attempt
		if i < t.maxRetries {
			time.Sleep(t.retryDelay * time.Duration(i+1))
		}
	}
	
	return resp, err
}
