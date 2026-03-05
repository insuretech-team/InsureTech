package sms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/stretchr/testify/require"
)

func TestPhoneHelpers(t *testing.T) {
	require.Equal(t, "8801712345678", NormalizeMSISDN("+8801712345678"))
	require.Equal(t, "8801712345678", NormalizeMSISDN("01712345678"))
	require.Equal(t, "", NormalizeMSISDN("12345"))

	require.True(t, ValidateMSISDN("+8801712345678"))
	require.False(t, ValidateMSISDN("invalid"))

	require.Equal(t, "GP", DetectCarrier("8801712345678"))
	require.Equal(t, "ROBI", DetectCarrier("8801612345678"))
	require.Equal(t, "BANGLALINK", DetectCarrier("8801912345678"))
	require.Equal(t, "TELETALK", DetectCarrier("8801512345678"))
	require.Equal(t, "UNKNOWN", DetectCarrier("8801212345678"))

	require.Equal(t, "8801XXX***678", MaskMSISDN("8801712345678"))
	require.Equal(t, "XXXXXXXXXXXXX", MaskMSISDN("short"))
}

func TestSSLWirelessClient_SendSMS_AndParseDLR(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","message_id":"mid-1"}`))
	}))
	defer srv.Close()

	cfg := &config.Config{}
	cfg.SMS.APIBase = srv.URL
	cfg.SMS.APIKey = "key"
	cfg.SMS.SID = "sid"
	cfg.SMS.MaskingEnabled = true
	cfg.SMS.MaskingSenderID = "LABAIDINS"
	cfg.SMS.NonMaskingEnabled = true
	cfg.SMS.NonMaskingSender = "12345"

	c := NewSSLWirelessClient(cfg)
	resp, err := c.SendSMS(context.Background(), &SendSMSRequest{
		MSISDN:     "01712345678",
		Message:    "hello",
		UseMasking: true,
		CSMSId:     "cs-1",
	})
	require.NoError(t, err)
	require.Equal(t, "mid-1", resp.MessageID)
	require.Equal(t, "PENDING", resp.Status)

	dlr, err := c.ParseDLRWebhook([]byte(`{"message_id":"m1","status":"DELIVERED"}`))
	require.NoError(t, err)
	require.Equal(t, "m1", dlr.MessageID)
	require.Equal(t, "DELIVERED", dlr.Status)
	_, err = c.ParseDLRWebhook([]byte(`{bad`))
	require.Error(t, err)
}
