// Package email provides SMTP-based email sending for the authn microservice.
// Used exclusively for: Business Beneficiary and System User authentication flows.
// Handles: email verification OTP, email login OTP, password reset OTP.
package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// Client is the SMTP email client
type Client struct {
	config Config
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
	return &Client{config: cfg}
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

// buildOTPEmail builds subject and HTML body for OTP email based on purpose
func buildOTPEmail(req *SendOTPRequest) (subject, body string) {
	switch req.Purpose {
	case "email_verification":
		subject = "Verify Your Email Address - InsureTech"
		body = fmt.Sprintf(`
<html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="background:#1a56db;padding:20px;border-radius:8px 8px 0 0;">
    <h1 style="color:white;margin:0;font-size:24px;">InsureTech</h1>
  </div>
  <div style="background:#f9fafb;padding:30px;border-radius:0 0 8px 8px;border:1px solid #e5e7eb;">
    <h2 style="color:#111827;">Verify Your Email Address</h2>
    <p style="color:#6b7280;">Please use the following OTP to verify your email address. This code is valid for <strong>%d minutes</strong>.</p>
    <div style="background:#fff;border:2px solid #1a56db;border-radius:8px;padding:20px;text-align:center;margin:24px 0;">
      <span style="font-size:36px;font-weight:bold;letter-spacing:8px;color:#1a56db;">%s</span>
    </div>
    <p style="color:#ef4444;font-size:14px;"><strong>Do not share this code with anyone.</strong> InsureTech staff will never ask for your OTP.</p>
    <p style="color:#6b7280;font-size:12px;margin-top:24px;">If you did not request this, please ignore this email.</p>
  </div>
</body></html>`, req.ExpiryMin, req.OTPCode)

	case "email_login":
		subject = "Your Login OTP - InsureTech Portal"
		body = fmt.Sprintf(`
<html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="background:#1a56db;padding:20px;border-radius:8px 8px 0 0;">
    <h1 style="color:white;margin:0;font-size:24px;">InsureTech Portal</h1>
  </div>
  <div style="background:#f9fafb;padding:30px;border-radius:0 0 8px 8px;border:1px solid #e5e7eb;">
    <h2 style="color:#111827;">Your Login Code</h2>
    <p style="color:#6b7280;">Use this one-time code to log in to the InsureTech Portal. Valid for <strong>%d minutes</strong>.</p>
    <div style="background:#fff;border:2px solid #059669;border-radius:8px;padding:20px;text-align:center;margin:24px 0;">
      <span style="font-size:36px;font-weight:bold;letter-spacing:8px;color:#059669;">%s</span>
    </div>
    <p style="color:#ef4444;font-size:14px;"><strong>Do not share this code.</strong> InsureTech will never ask for your OTP.</p>
    <p style="color:#6b7280;font-size:12px;margin-top:24px;">If you did not attempt to log in, please contact support immediately and secure your account.</p>
  </div>
</body></html>`, req.ExpiryMin, req.OTPCode)

	case "password_reset_email":
		subject = "Password Reset Request - InsureTech"
		body = fmt.Sprintf(`
<html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="background:#dc2626;padding:20px;border-radius:8px 8px 0 0;">
    <h1 style="color:white;margin:0;font-size:24px;">InsureTech - Password Reset</h1>
  </div>
  <div style="background:#f9fafb;padding:30px;border-radius:0 0 8px 8px;border:1px solid #e5e7eb;">
    <h2 style="color:#111827;">Reset Your Password</h2>
    <p style="color:#6b7280;">We received a request to reset your password. Use this OTP to complete the process. Valid for <strong>%d minutes</strong>.</p>
    <div style="background:#fff;border:2px solid #dc2626;border-radius:8px;padding:20px;text-align:center;margin:24px 0;">
      <span style="font-size:36px;font-weight:bold;letter-spacing:8px;color:#dc2626;">%s</span>
    </div>
    <p style="color:#ef4444;font-size:14px;"><strong>Do not share this code.</strong></p>
    <p style="color:#6b7280;font-size:12px;margin-top:24px;">If you did not request a password reset, please ignore this email. Your password will not be changed.</p>
  </div>
</body></html>`, req.ExpiryMin, req.OTPCode)

	default:
		subject = "Your OTP - InsureTech"
		body = fmt.Sprintf(`
<html><body style="font-family:Arial,sans-serif;">
  <p>Your OTP is: <strong style="font-size:24px;letter-spacing:4px;">%s</strong></p>
  <p>Valid for %d minutes. Do not share this code.</p>
</body></html>`, req.OTPCode, req.ExpiryMin)
	}
	return subject, body
}

// buildMIMEMessage builds a MIME-formatted email message
func buildMIMEMessage(from, to, subject, htmlBody string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: InsureTech <%s>\r\n", from))
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
