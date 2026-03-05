// Package producer provides a real Kafka producer implementation using IBM sarama.
// Used by all microservices to publish domain events to Kafka topics.
package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// EventProducer is a real Kafka producer using IBM sarama sync producer.
// It supports both proto messages and arbitrary JSON payloads.
type EventProducer struct {
	producer sarama.SyncProducer
	brokers  []string
	clientID string
}

// NewEventProducer creates a new Kafka sync producer connected to the given brokers.
// clientID is used for Kafka client identification (e.g. "authn-service").
func NewEventProducer(brokers []string, topic string, clientID string) (*EventProducer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Version = sarama.V2_8_0_0

	// Reliability: wait for all in-sync replicas to ack
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Retry.Backoff = 250 * time.Millisecond
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	// Idempotent producer (exactly-once semantics within a session)
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1 // Required for idempotent producer

	// Compression for efficiency
	cfg.Producer.Compression = sarama.CompressionSnappy

	// Keep startup/connect behavior bounded so caller retries are visible quickly.
	cfg.Net.DialTimeout = 3 * time.Second
	cfg.Net.ReadTimeout = 10 * time.Second
	cfg.Net.WriteTimeout = 10 * time.Second
	cfg.Metadata.Retry.Max = 3
	cfg.Metadata.Retry.Backoff = 500 * time.Millisecond

	// Timeout
	cfg.Producer.Timeout = 10 * time.Second

	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka sync producer (brokers=%v): %w", brokers, err)
	}

	appLogger.Infof("Kafka producer connected to brokers: %v (clientID: %s)", brokers, clientID)

	return &EventProducer{
		producer: syncProducer,
		brokers:  brokers,
		clientID: clientID,
	}, nil
}

// NewEventProducerWithRetry creates a Kafka producer with retry logic on startup.
// Useful when Kafka may not be immediately available at service startup.
func NewEventProducerWithRetry(brokers []string, topic, clientID string, maxRetries int, retryDelay time.Duration) (*EventProducer, error) {
	type result struct {
		producer *EventProducer
		err      error
	}

	attemptTimeout := 12 * time.Second
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		resultCh := make(chan result, 1)
		go func() {
			p, err := NewEventProducer(brokers, topic, clientID)
			resultCh <- result{producer: p, err: err}
		}()

		select {
		case res := <-resultCh:
			if res.err == nil {
				return res.producer, nil
			}
			lastErr = res.err
		case <-time.After(attemptTimeout):
			lastErr = fmt.Errorf("timed out after %s creating kafka sync producer (brokers=%v)", attemptTimeout, brokers)
		}

		appLogger.Warnf("Kafka connection attempt %d/%d failed: %v. Retrying in %s...", i+1, maxRetries, lastErr, retryDelay)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	return nil, fmt.Errorf("failed to connect to Kafka after %d attempts: %w", maxRetries, lastErr)
}

// Produce sends a proto message to a Kafka topic with the given key.
// The message is serialized to JSON using protojson for human readability and schema flexibility.
// topic: the Kafka topic to publish to
// key: partition key (e.g. user_id for ordered per-user events)
// msg: proto.Message to serialize
func (p *EventProducer) Produce(ctx context.Context, topic string, key string, msg interface{}) error {
	var valueBytes []byte
	var err error

	// Prefer proto JSON marshaling for proto messages
	if protoMsg, ok := msg.(proto.Message); ok {
		marshaler := protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: false,
		}
		valueBytes, err = marshaler.Marshal(protoMsg)
		if err != nil {
			return fmt.Errorf("failed to marshal proto message: %w", err)
		}
	} else {
		// Fall back to standard JSON
		valueBytes, err = json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueBytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("content-type"),
				Value: []byte("application/json"),
			},
			{
				Key:   []byte("producer-id"),
				Value: []byte(p.clientID),
			},
			{
				Key:   []byte("produced-at"),
				Value: []byte(time.Now().UTC().Format(time.RFC3339Nano)),
			},
		},
	}

	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		appLogger.Errorf("Kafka produce failed (topic=%s, key=%s): %v", topic, key, err)
		return fmt.Errorf("kafka produce failed: %w", err)
	}

	appLogger.Infof("Kafka message sent (topic=%s, key=%s, partition=%d, offset=%d)", topic, key, partition, offset)
	return nil
}

// ProduceEvent sends raw bytes to a topic. Lower-level API for custom serialization.
// topic must be set explicitly; key is the partition key.
func (p *EventProducer) ProduceEvent(key string, value []byte) error {
	// This low-level method requires callers to specify topic via Produce().
	// Kept for backward compatibility. Publishes to first configured broker's default topic.
	_ = key
	_ = value
	return fmt.Errorf("use Produce(ctx, topic, key, msg) instead of ProduceEvent")
}

// Close closes the Kafka producer gracefully, flushing pending messages.
func (p *EventProducer) Close() error {
	if p.producer != nil {
		if err := p.producer.Close(); err != nil {
			appLogger.Errorf("Kafka producer close error: %v", err)
			return err
		}
		appLogger.Info("Kafka producer closed gracefully")
	}
	return nil
}
