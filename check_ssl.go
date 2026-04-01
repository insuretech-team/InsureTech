package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

type sendSMSRequest struct {
	APIToken string `json:"api_token"`
	SID      string `json:"sid"`
	MSISDN   string `json:"msisdn"`
	SMS      string `json:"sms"`
	CSMSID   string `json:"csms_id"`
	SMSType  string `json:"sms_type,omitempty"`
}

func main() {
	otp, err := otpValue()
	if err != nil {
		fmt.Printf("failed to generate otp: %v\n", err)
		os.Exit(1)
	}

	payload := sendSMSRequest{
		APIToken: envOrDefault("SSLWIRELESS_API_KEY", "ngycabv6-uvz4o4o1-seg47hmg-xuswqyho-aa98j6ir"),
		SID:      envOrDefault("SSLWIRELESS_SID", "LIFEPLUSBDBRAND"),
		MSISDN:   envOrDefault("SSLWIRELESS_MSISDN", "8801913705269"),
		SMS:      envOrDefault("SSLWIRELESS_SMS", fmt.Sprintf("Your OTP is %s", otp)),
		CSMSID:   envOrDefault("SSLWIRELESS_CSMS_ID", fmt.Sprintf("cred-check-%d", time.Now().Unix())),
		SMSType:  envOrDefault("SSLWIRELESS_SMS_TYPE", "EN"),
	}

	apiBase := strings.TrimRight(envOrDefault("SSLWIRELESS_API_BASE", "https://smsplus.sslwireless.com"), "/")
	requestBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failed to marshal request: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, apiBase+"/api/v3/send-sms", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Printf("failed to build request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}

	fmt.Printf("checking SSL Wireless credentials with SID=%s token=%s msisdn=%s otp=%s\n", payload.SID, maskSecret(payload.APIToken), payload.MSISDN, otp)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("http_status=%s\n", resp.Status)
	fmt.Printf("response=%s\n", strings.TrimSpace(string(body)))

	if resp.StatusCode >= http.StatusBadRequest {
		os.Exit(1)
	}
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return strings.Trim(value, "\"")
}

func maskSecret(value string) string {
	if len(value) <= 8 {
		return "********"
	}

	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

func otpValue() (string, error) {
	if value := strings.TrimSpace(os.Getenv("SSLWIRELESS_OTP")); value != "" {
		return value, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%04d", n.Int64()), nil
}
