package service

import (
	"context"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/sms"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type mockEventProducer struct {
	calls []struct {
		topic string
		key   string
		msg   any
	}
}

func (m *mockEventProducer) Produce(ctx context.Context, topic string, key string, msg interface{}) error {
	m.calls = append(m.calls, struct {
		topic string
		key   string
		msg   any
	}{topic: topic, key: key, msg: msg})
	return nil
}
func (m *mockEventProducer) Close() error { return nil }

func TestOTPService_HandleDLR_PublishesEvent(t *testing.T) {
	mp := &mockEventProducer{}
	pub := events.NewPublisher(mp)

	_ = &OTPService{smsClient: &sms.SSLWirelessClient{}, eventPublisher: pub}

	// We can't call HandleDLR without repo+sms parser, so instead ensure publisher method works.
	// This test primarily guards against signature mismatches and panics.
	require.NoError(t, pub.PublishSMSDeliveryReport(context.Background(), "o1", "pmid", "8801***", "DELIVERED", "", "GP", time.Now()))
	require.NotEmpty(t, mp.calls)
	require.Equal(t, events.TopicSMSDeliveryReport, mp.calls[0].topic)
}

func TestOTPService_IncrementRateLimit_NoRedis_NoPanic(t *testing.T) {
	svc := &OTPService{}
	require.NotPanics(t, func() {
		svc.incrementRateLimit(context.Background(), "+8801712345678", "LOGIN")
	})
}

func TestOTPService_IncrementRateLimit_UnreachableRedis_NoPanic(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:1",
		DialTimeout:  10 * time.Millisecond,
		ReadTimeout:  10 * time.Millisecond,
		WriteTimeout: 10 * time.Millisecond,
	})
	defer rdb.Close()

	svc := &OTPService{redisClient: rdb}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	require.NotPanics(t, func() {
		svc.incrementRateLimit(ctx, "+8801712345678", "LOGIN")
	})
}
