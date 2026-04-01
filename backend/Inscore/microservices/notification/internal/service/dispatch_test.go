package service

import (
	"errors"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/delivery"
	"github.com/stretchr/testify/require"
)

func TestShouldRetryStopsOnPermanentError(t *testing.T) {
	svc := &Service{
		cfg: &config.Config{
			NotificationRetry: config.NotificationRetryConfig{
				MaxAttempts: 3,
				Backoff:     []time.Duration{time.Second},
			},
		},
	}

	require.False(t, svc.shouldRetry(1, delivery.Permanent(delivery.ErrPushNoActiveDeviceTokens)))
	require.True(t, svc.shouldRetry(1, errors.New("temporary failure")))
}

func TestWebhookRetryBackoffUsesLastValueWhenAttemptsExceedConfiguredSlots(t *testing.T) {
	svc := &Service{
		cfg: &config.Config{
			Webhook: config.WebhookConfig{
				Backoff: []time.Duration{time.Second, 5 * time.Second},
			},
		},
	}

	require.Equal(t, time.Second, svc.webhookRetryBackoff(1))
	require.Equal(t, 5*time.Second, svc.webhookRetryBackoff(3))
}
