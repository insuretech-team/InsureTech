package delivery

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	From     string
	Username string
	Password string
	TLS      bool
}

type EmailClient struct {
	config EmailConfig
}

type EmailResponse struct {
	MessageID string
	SentAt    time.Time
}

func NewEmailClient(cfg EmailConfig) *EmailClient {
	return &EmailClient{config: cfg}
}

func (c *EmailClient) Send(ctx context.Context, to, subject, body string) (*EmailResponse, error) {
	if strings.TrimSpace(to) == "" {
		return nil, errors.New("email recipient is required")
	}
	if strings.TrimSpace(c.config.SMTPHost) == "" || c.config.SMTPPort == 0 {
		return nil, errors.New("email delivery is not configured")
	}

	msg := buildMIMEMessage(c.config.From, to, subject, body)
	if err := c.send(ctx, to, msg); err != nil {
		return nil, err
	}

	return &EmailResponse{
		MessageID: uuid.NewString(),
		SentAt:    time.Now(),
	}, nil
}

func (c *EmailClient) send(ctx context.Context, to, msg string) error {
	addr := net.JoinHostPort(c.config.SMTPHost, fmt.Sprintf("%d", c.config.SMTPPort))
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.SMTPHost)

	if c.config.TLS {
		dialer := &net.Dialer{Timeout: 10 * time.Second}
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			ServerName: c.config.SMTPHost,
			MinVersion: tls.VersionTLS12,
		})
		if err != nil {
			return fmt.Errorf("email TLS dial failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, c.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("email SMTP client creation failed: %w", err)
		}
		defer client.Close()

		if deadline, ok := ctx.Deadline(); ok {
			_ = conn.SetDeadline(deadline)
		}
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("email SMTP auth failed: %w", err)
		}
		if err := client.Mail(c.config.From); err != nil {
			return fmt.Errorf("email MAIL FROM failed: %w", err)
		}
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("email RCPT TO failed: %w", err)
		}
		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("email DATA failed: %w", err)
		}
		if _, err := writer.Write([]byte(msg)); err != nil {
			_ = writer.Close()
			return fmt.Errorf("email write failed: %w", err)
		}
		if err := writer.Close(); err != nil {
			return fmt.Errorf("email DATA close failed: %w", err)
		}
		return nil
	}

	if err := smtp.SendMail(addr, auth, c.config.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("email SendMail failed: %w", err)
	}
	appLogger.Infof("notification email accepted by SMTP server for %s", maskEmail(to))
	return nil
}

func buildMIMEMessage(from, to, subject, body string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", to))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	builder.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	builder.WriteString("\r\n")
	builder.WriteString(body)
	return builder.String()
}

func maskEmail(email string) string {
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
