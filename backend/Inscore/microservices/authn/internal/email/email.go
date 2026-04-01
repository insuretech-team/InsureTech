// Package email provides SMTP-based email sending for the authn microservice.
// Used exclusively for: Business Beneficiary and System User authentication flows.
// Handles: email verification OTP, email login OTP, password reset OTP.
package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// Client is the email client that supports both SMTP2Go HTTP API and fallback SMTP
type Client struct {
	config        Config
	smtp2goClient *SMTP2GoClient
}

// Config holds SMTP configuration
type Config struct {
	SMTPHost string
	SMTPPort int
	From     string
	Username string
	Password string
	TLS      bool
}

// NewClient creates a new email client
func NewClient(cfg Config) *Client {
	client := &Client{config: cfg}
	// Initialize SMTP2GoClient if API key is configured
	smtp2goClient := NewSMTP2GoClient(cfg.From)
	if smtp2goClient.IsConfigured() {
		client.smtp2goClient = smtp2goClient
	}
	return client
}

// SendOTPRequest is the request to send an email OTP
type SendOTPRequest struct {
	To        string // recipient email
	OTPCode   string // 6-digit OTP
	Purpose   string // email_verification, email_login, password_reset_email
	ExpiryMin int    // minutes until expiry
}

// SendOTPResponse is the response from sending an email OTP
type SendOTPResponse struct {
	MessageID string
	SentAt    time.Time
}

// SendOTP sends an OTP email to the recipient
func (c *Client) SendOTP(req *SendOTPRequest) (*SendOTPResponse, error) {
	subject, body := buildOTPEmail(req)

	// Try SMTP2Go first if configured
	if c.smtp2goClient != nil && c.smtp2goClient.IsConfigured() {
		ctx := context.Background()
		if err := c.smtp2goClient.Send(ctx, req.To, subject, body, ""); err != nil {
			// Log the failure but fall through to SMTP fallback
			logger.Errorf("SMTP2Go send failed, falling back to SMTP: %v", err)
		} else {
			return &SendOTPResponse{
				MessageID: fmt.Sprintf("email-%d", time.Now().UnixNano()),
				SentAt:    time.Now(),
			}, nil
		}
	}

	// Fallback to standard SMTP
	msg := buildMIMEMessage(c.config.From, req.To, subject, body)

	if err := c.send(req.To, msg); err != nil {
		return nil, fmt.Errorf("failed to send OTP email to %s: %w", MaskEmail(req.To), err)
	}

	return &SendOTPResponse{
		MessageID: fmt.Sprintf("email-%d", time.Now().UnixNano()),
		SentAt:    time.Now(),
	}, nil
}

// send sends an email via SMTP
func (c *Client) send(to, msg string) error {
	addr := fmt.Sprintf("%s:%d", c.config.SMTPHost, c.config.SMTPPort)
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.SMTPHost)

	if c.config.TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         c.config.SMTPHost,
			MinVersion:         tls.VersionTLS12,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			logger.Errorf("TLS dial failed: %v", err)
			return errors.New("TLS dial failed")
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, c.config.SMTPHost)
		if err != nil {
			logger.Errorf("SMTP client creation failed: %v", err)
			return errors.New("SMTP client creation failed")
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			logger.Errorf("SMTP auth failed: %v", err)
			return errors.New("SMTP auth failed")
		}
		if err = client.Mail(c.config.From); err != nil {
			logger.Errorf("SMTP MAIL FROM failed: %v", err)
			return errors.New("SMTP MAIL FROM failed")
		}
		if err = client.Rcpt(to); err != nil {
			logger.Errorf("SMTP RCPT TO failed: %v", err)
			return errors.New("SMTP RCPT TO failed")
		}
		w, err := client.Data()
		if err != nil {
			logger.Errorf("SMTP DATA failed: %v", err)
			return errors.New("SMTP DATA failed")
		}
		_, err = w.Write([]byte(msg))
		if err != nil {
			logger.Errorf("SMTP write failed: %v", err)
			return errors.New("SMTP write failed")
		}
		return w.Close()
	}

	// Plain SMTP (STARTTLS negotiated by smtp.SendMail)
	return smtp.SendMail(addr, auth, c.config.From, []string{to}, []byte(msg))
}

// buildOTPEmail builds subject and HTML body for OTP email based on purpose.
// Loads the branded HTML template from templates/common/otp_email.html.
// Falls back to inline HTML if the template cannot be found.
func buildOTPEmail(req *SendOTPRequest) (subject, body string) {
	switch req.Purpose {
	case "email_verification":
		subject = "ল্যাবএইড ইন্স্যুরটেক — ইমেইল যাচাই কোড"
	case "email_login":
		subject = "ল্যাবএইড ইন্স্যুরটেক — লগইন কোড"
	case "password_reset_email":
		subject = "ল্যাবএইড ইন্স্যুরটেক — পাসওয়ার্ড রিসেট কোড"
	default:
		subject = "ল্যাবএইড ইন্স্যুরটেক — যাচাই কোড"
	}

	// Try loading the branded HTML template.
	// Docker container: binary at /app/server, templates at /app/backend/inscore/templates/
	// Local dev:        templates relative to source or binary CWD.
	candidates := []string{
		"/app/backend/inscore/templates/common/otp_email.html", // Docker (absolute)
		"backend/inscore/templates/common/otp_email.html",      // CWD = /app in Docker
		"templates/common/otp_email.html",                      // CWD = templates parent
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		// Local dev: resolve relative to this source file
		candidates = append(candidates, filepath.Join(
			filepath.Dir(thisFile), "..", "..", "..", "..", "templates", "common", "otp_email.html",
		))
	}

	var tmplContent string
	for _, p := range candidates {
		b, readErr := os.ReadFile(p)
		if readErr == nil {
			tmplContent = string(b)
			logger.Infof("email: loaded OTP template from %s", p)
			break
		}
	}
	if tmplContent == "" {
		logger.Warnf("email: otp_email.html template not found in any candidate path — using fallback inline HTML")
	}

	if tmplContent != "" {
		tmpl, err := template.New("otp_email").Parse(tmplContent)
		if err == nil {
			var buf bytes.Buffer
			data := struct {
				Code          string
				ExpiryMinutes int
				Purpose       string
				Year          int
			}{Code: req.OTPCode, ExpiryMinutes: req.ExpiryMin, Purpose: req.Purpose, Year: time.Now().Year()}
			if err := tmpl.Execute(&buf, data); err == nil {
				return subject, buf.String()
			}
		}
		logger.Errorf("failed to render otp_email.html template: %v", err)
	}

	// Fallback: minimal inline HTML
	body = fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:560px;margin:40px auto;padding:24px;background:#f4f6f9;">
  <div style="background:#1a5276;padding:24px;border-radius:8px 8px 0 0;text-align:center;">
    <h1 style="color:#fff;margin:0;font-size:22px;">ল্যাবএইড ইন্স্যুরটেক</h1>
  </div>
  <div style="background:#fff;padding:32px;border-radius:0 0 8px 8px;border:1px solid #e5e7eb;">
    <h2 style="color:#1a5276;">আপনার যাচাই কোড</h2>
    <div style="background:#f0f7ff;border:2px dashed #2e86c1;border-radius:8px;text-align:center;padding:24px;margin:20px 0;">
      <span style="font-size:40px;font-weight:700;letter-spacing:10px;color:#1a5276;font-family:monospace;">%s</span>
      <p style="color:#e74c3c;margin-top:10px;font-size:13px;">মেয়াদ: %d মিনিট</p>
    </div>
    <p style="color:#e74c3c;font-size:13px;"><strong>এই কোড কাউকে শেয়ার করবেন না।</strong></p>
  </div>
</body></html>`, req.OTPCode, req.ExpiryMin)
	return subject, body
}

// buildMIMEMessage builds a MIME-formatted email message
func buildMIMEMessage(from, to, subject, htmlBody string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: =?UTF-8?B?4KaX4Ka+4KaH4KaoIOCmoeCmv+CmsOCnjeCmruCmvuCmsOCmvg==?= <%s>\r\n", from)) // ল্যাবএইড ইন্স্যুরটেক
	sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	sb.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(htmlBody)
	return sb.String()
}

// MaskEmail masks an email for safe logging: user@domain.com → u***@domain.com
func MaskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return "***"
	}
	local := parts[0]
	if len(local) <= 1 {
		return "***@" + parts[1]
	}
	return string(local[0]) + "***@" + parts[1]
}
