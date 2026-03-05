package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailHelpers(t *testing.T) {
	sub, body := buildOTPEmail(&SendOTPRequest{
		To:        "user@example.com",
		OTPCode:   "123456",
		Purpose:   "email_verification",
		ExpiryMin: 5,
	})
	require.Contains(t, sub, "Verify")
	require.Contains(t, body, "123456")

	sub2, body2 := buildOTPEmail(&SendOTPRequest{
		To:        "user@example.com",
		OTPCode:   "654321",
		Purpose:   "unknown",
		ExpiryMin: 5,
	})
	require.Contains(t, sub2, "OTP")
	require.Contains(t, body2, "654321")

	msg := buildMIMEMessage("from@example.com", "to@example.com", "Subj", "<b>Hi</b>")
	require.Contains(t, msg, "Subject: Subj")
	require.Contains(t, msg, "Content-Type: text/html")

	require.Equal(t, "u***@example.com", MaskEmail("user@example.com"))
	require.Equal(t, "***", MaskEmail("invalid"))
}

func TestEmailClient_SendOTP_ErrorPath(t *testing.T) {
	c := NewClient(Config{
		SMTPHost: "127.0.0.1",
		SMTPPort: 1,
		From:     "from@example.com",
		Username: "u",
		Password: "p",
		TLS:      false,
	})

	_, err := c.SendOTP(&SendOTPRequest{
		To:        "user@example.com",
		OTPCode:   "123456",
		Purpose:   "email_login",
		ExpiryMin: 5,
	})
	require.Error(t, err)
}
