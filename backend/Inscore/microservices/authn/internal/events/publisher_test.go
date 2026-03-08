package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type produceCall struct {
	topic string
	key   string
	msg   any
}

type mockProducer struct {
	calls []produceCall
}

func (m *mockProducer) Produce(ctx context.Context, topic string, key string, msg interface{}) error {
	m.calls = append(m.calls, produceCall{topic: topic, key: key, msg: msg})
	return nil
}

func (m *mockProducer) Close() error { return nil }

type blockingProducer struct{}

func (b *blockingProducer) Produce(ctx context.Context, topic string, key string, msg interface{}) error {
	<-ctx.Done()
	return ctx.Err()
}

func (b *blockingProducer) Close() error { return nil }

func TestPublisher_TypedNilProducer_DropsEvent(t *testing.T) {
	var nilProducer *mockProducer
	p := NewPublisher(nilProducer)

	require.NotPanics(t, func() {
		err := p.PublishUserLoggedIn(context.Background(), "u1", "s1", "SERVER_SIDE", "127.0.0.1", "WEB", "ua")
		require.NoError(t, err)
	})
}

func TestPublisher_AllMethods_Produce(t *testing.T) {
	ctx := context.Background()
	mp := &mockProducer{}
	p := NewPublisher(mp)

	require.NoError(t, p.PublishUserRegistered(ctx, "u1", "+8801xxx", "a@b.com", "1.1.1.1", "WEB"))
	require.NoError(t, p.PublishUserLoggedIn(ctx, "u1", "s1", "JWT", "1.1.1.1", "MOBILE", "ua"))
	require.NoError(t, p.PublishUserLoggedOut(ctx, "u1", "s1", "JWT", "user_initiated", "1.1.1.1", ""))
	require.NoError(t, p.PublishOTPSent(ctx, "o1", "+8801***", "login", "sms", "sslwireless", "", "", true))
	require.NoError(t, p.PublishOTPVerified(ctx, "o1", "u1", 1))
	require.NoError(t, p.PublishEmailVerificationSent(ctx, "u1", "user@domain.com", "o2", "email_verification", "1.1.1.1"))
	require.NoError(t, p.PublishEmailVerified(ctx, "u1", "user@domain.com"))
	require.NoError(t, p.PublishEmailLoginSucceeded(ctx, "u1", "s2", "user@domain.com", "SYSTEM_USER", "1.1.1.1", "ua", "dev"))
	require.NoError(t, p.PublishEmailLoginFailed(ctx, "user@domain.com", "invalid_otp", 2, "1.1.1.1", "ua"))
	require.NoError(t, p.PublishPasswordResetByEmailRequested(ctx, "u1", "user@domain.com", "o3", "1.1.1.1"))
	require.NoError(t, p.PublishTokenRefreshed(ctx, "u1", "s3", "old", "new", "newr", "1.1.1.1", "", "ua"))
	require.NoError(t, p.PublishSessionRevoked(ctx, "u1", "s3", "JWT", "system", "security"))
	require.NoError(t, p.PublishPasswordChanged(ctx, "u1", "1.1.1.1", "u1"))
	require.NoError(t, p.PublishPasswordResetRequested(ctx, "u1", "+8801xxx", "1.1.1.1", ""))
	require.NoError(t, p.PublishLoginFailed(ctx, "u1", "+8801xxx", "invalid_password", "1.1.1.1", "", "ua", 3))
	require.NoError(t, p.PublishCSRFValidationFailed(ctx, "u1", "s1", "eh", "rh", "1.1.1.1", "ua", "/x", "POST"))
	require.NoError(t, p.PublishSMSDeliveryReport(ctx, "o1", "pmid", "8801***", "DELIVERED", "", "GP", time.Now()))
	require.NoError(t, p.PublishSessionExpired(ctx, "u1", "s1", "JWT", time.Now(), 10))
	require.NoError(t, p.PublishAccountLocked(ctx, "u1", "failed_login", time.Now().Add(time.Minute)))

	require.GreaterOrEqual(t, len(mp.calls), 19)

	// Spot-check topics
	require.Equal(t, TopicUserRegistered, mp.calls[0].topic)
	require.Equal(t, "u1", mp.calls[0].key)
}

func TestPublisher_Publish_BoundsSlowProducer(t *testing.T) {
	p := NewPublisher(&blockingProducer{})

	start := time.Now()
	err := p.publish(context.Background(), TopicUserLoggedIn, "u1", map[string]string{"ok": "true"})
	elapsed := time.Since(start)

	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
	require.Less(t, elapsed, 2*time.Second)
}
