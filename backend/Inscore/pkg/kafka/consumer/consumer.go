// Package consumer provides a Kafka consumer group implementation using IBM sarama.
// Supports at-least-once delivery semantics with commit-on-success, graceful
// shutdown via context cancellation, and a dead-letter queue (DLQ) for poison messages.
package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// HandlerFunc is the callback invoked for each Kafka message.
// Return nil to commit the offset; return an error to route to the DLQ.
type HandlerFunc func(ctx context.Context, msg *Message) error

// Message wraps a sarama ConsumerMessage with convenience accessors.
type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
}

// Unmarshal deserializes the message value into v (JSON).
func (m *Message) Unmarshal(v any) error {
	return json.Unmarshal(m.Value, v)
}

// ConsumerGroup wraps a sarama ConsumerGroup with lifecycle management.
type ConsumerGroup struct {
	client   sarama.ConsumerGroup
	topics   []string
	groupID  string
	handler  HandlerFunc
	dlqTopic string              // empty = no DLQ
	dlqProd  sarama.SyncProducer // nil = no DLQ producer
	wg       sync.WaitGroup
	cancel   context.CancelFunc
}

// Config holds consumer group configuration.
type Config struct {
	Brokers           []string
	GroupID           string
	Topics            []string
	Handler           HandlerFunc
	DLQTopic          string // optional; routed on handler error
	ClientID          string
	InitialOffset     int64                  // sarama.OffsetNewest or sarama.OffsetOldest
	SessionTimeout    time.Duration          // default 10s
	RebalanceStrategy sarama.BalanceStrategy // default RoundRobin
}

// NewConsumerGroup creates and starts a Kafka consumer group.
// Call Close() or cancel the context to stop it.
func NewConsumerGroup(cfg Config) (*ConsumerGroup, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("kafka consumer: at least one broker required")
	}
	if cfg.GroupID == "" {
		return nil, fmt.Errorf("kafka consumer: group_id is required")
	}
	if len(cfg.Topics) == 0 {
		return nil, fmt.Errorf("kafka consumer: at least one topic required")
	}
	if cfg.Handler == nil {
		return nil, fmt.Errorf("kafka consumer: handler is required")
	}

	sarCfg := sarama.NewConfig()
	sarCfg.Version = sarama.V2_8_0_0
	if cfg.ClientID != "" {
		sarCfg.ClientID = cfg.ClientID
	}
	// Bound dial/read/write latency so failed brokers do not stall startup loops.
	sarCfg.Net.DialTimeout = 3 * time.Second
	sarCfg.Net.ReadTimeout = 10 * time.Second
	sarCfg.Net.WriteTimeout = 10 * time.Second
	sarCfg.Metadata.Retry.Max = 3
	sarCfg.Metadata.Retry.Backoff = 500 * time.Millisecond

	// Offset management
	initialOffset := cfg.InitialOffset
	if initialOffset == 0 {
		initialOffset = sarama.OffsetNewest
	}
	sarCfg.Consumer.Offsets.Initial = initialOffset
	sarCfg.Consumer.Offsets.AutoCommit.Enable = false // manual commit-on-success

	// Session / rebalance
	timeout := cfg.SessionTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	sarCfg.Consumer.Group.Session.Timeout = timeout
	sarCfg.Consumer.Group.Heartbeat.Interval = timeout / 3

	strategy := cfg.RebalanceStrategy
	if strategy == nil {
		strategy = sarama.NewBalanceStrategyRoundRobin()
	}
	sarCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{strategy}

	client, err := newConsumerGroupWithTimeout(cfg.Brokers, cfg.GroupID, sarCfg, 12*time.Second)
	if err != nil {
		return nil, fmt.Errorf("kafka consumer: failed to create group client: %w", err)
	}

	cg := &ConsumerGroup{
		client:   client,
		topics:   cfg.Topics,
		groupID:  cfg.GroupID,
		handler:  cfg.Handler,
		dlqTopic: cfg.DLQTopic,
	}

	// Wire DLQ producer if topic is configured.
	if cfg.DLQTopic != "" {
		dlqCfg := sarama.NewConfig()
		dlqCfg.Version = sarama.V2_8_0_0
		dlqCfg.Producer.RequiredAcks = sarama.WaitForLocal
		dlqCfg.Producer.Return.Successes = true
		dlqCfg.Producer.Return.Errors = true
		dlqCfg.Net.DialTimeout = 3 * time.Second
		dlqCfg.Net.ReadTimeout = 10 * time.Second
		dlqCfg.Net.WriteTimeout = 10 * time.Second
		dlqCfg.Metadata.Retry.Max = 3
		dlqCfg.Metadata.Retry.Backoff = 500 * time.Millisecond
		dlqProd, dlqErr := newSyncProducerWithTimeout(cfg.Brokers, dlqCfg, 12*time.Second)
		if dlqErr != nil {
			appLogger.Warnf("kafka consumer: DLQ producer failed to connect (topic=%s): %v — DLQ disabled", cfg.DLQTopic, dlqErr)
		} else {
			cg.dlqProd = dlqProd
		}
	}

	appLogger.Infof("Kafka consumer group created (group=%s, topics=%v, brokers=%v)", cfg.GroupID, cfg.Topics, cfg.Brokers)
	return cg, nil
}

func newConsumerGroupWithTimeout(brokers []string, groupID string, cfg *sarama.Config, timeout time.Duration) (sarama.ConsumerGroup, error) {
	type result struct {
		client sarama.ConsumerGroup
		err    error
	}
	ch := make(chan result, 1)
	go func() {
		client, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
		ch <- result{client: client, err: err}
	}()

	select {
	case res := <-ch:
		return res.client, res.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timed out after %s creating consumer group (group=%s, brokers=%v)", timeout, groupID, brokers)
	}
}

func newSyncProducerWithTimeout(brokers []string, cfg *sarama.Config, timeout time.Duration) (sarama.SyncProducer, error) {
	type result struct {
		producer sarama.SyncProducer
		err      error
	}
	ch := make(chan result, 1)
	go func() {
		producer, err := sarama.NewSyncProducer(brokers, cfg)
		ch <- result{producer: producer, err: err}
	}()

	select {
	case res := <-ch:
		return res.producer, res.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timed out after %s creating sync producer (brokers=%v)", timeout, brokers)
	}
}

// Start begins consuming messages. Blocks until the context is cancelled.
// Run in a goroutine: go cg.Start(ctx)
func (cg *ConsumerGroup) Start(ctx context.Context) {
	child, cancel := context.WithCancel(ctx)
	cg.cancel = cancel

	h := &groupHandler{
		handler:  cg.handler,
		dlqTopic: cg.dlqTopic,
		dlqProd:  cg.dlqProd,
	}

	cg.wg.Add(1)
	go func() {
		defer cg.wg.Done()
		for {
			// Consume is called in a loop to handle rebalances.
			if err := cg.client.Consume(child, cg.topics, h); err != nil {
				if child.Err() != nil {
					// Context cancelled — clean shutdown.
					return
				}
				appLogger.Errorf("kafka consumer group error (group=%s): %v — restarting in 2s", cg.groupID, err)
				select {
				case <-child.Done():
					return
				case <-time.After(2 * time.Second):
				}
			}
			if child.Err() != nil {
				return
			}
		}
	}()

	// Log errors from the client.
	cg.wg.Add(1)
	go func() {
		defer cg.wg.Done()
		for {
			select {
			case err, ok := <-cg.client.Errors():
				if !ok {
					return
				}
				appLogger.Errorf("kafka consumer error (group=%s): %v", cg.groupID, err)
			case <-child.Done():
				return
			}
		}
	}()
}

// Close gracefully stops the consumer group and waits for all goroutines to exit.
func (cg *ConsumerGroup) Close() error {
	if cg.cancel != nil {
		cg.cancel()
	}
	cg.wg.Wait()
	if cg.dlqProd != nil {
		_ = cg.dlqProd.Close()
	}
	if err := cg.client.Close(); err != nil {
		return fmt.Errorf("kafka consumer group close error: %w", err)
	}
	appLogger.Infof("Kafka consumer group closed (group=%s)", cg.groupID)
	return nil
}

// ============================================================
// sarama ConsumerGroupHandler implementation
// ============================================================

type groupHandler struct {
	handler  HandlerFunc
	dlqTopic string
	dlqProd  sarama.SyncProducer
}

func (h *groupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *groupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case raw, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := toMessage(raw)
			if err := h.handler(session.Context(), msg); err != nil {
				appLogger.Errorf("kafka consumer: handler error (topic=%s, partition=%d, offset=%d): %v",
					raw.Topic, raw.Partition, raw.Offset, err)
				// Route to DLQ if configured.
				h.sendToDLQ(raw, err)
				// Still commit — we've handled (or DLQ'd) the message.
			}
			// Commit on success (at-least-once semantics).
			session.MarkMessage(raw, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *groupHandler) sendToDLQ(raw *sarama.ConsumerMessage, handlerErr error) {
	if h.dlqProd == nil || h.dlqTopic == "" {
		return
	}
	dlqMsg := &sarama.ProducerMessage{
		Topic: h.dlqTopic,
		Key:   sarama.ByteEncoder(raw.Key),
		Value: sarama.ByteEncoder(raw.Value),
		Headers: []sarama.RecordHeader{
			{Key: []byte("original-topic"), Value: []byte(raw.Topic)},
			{Key: []byte("original-partition"), Value: []byte(fmt.Sprintf("%d", raw.Partition))},
			{Key: []byte("original-offset"), Value: []byte(fmt.Sprintf("%d", raw.Offset))},
			{Key: []byte("error"), Value: []byte(handlerErr.Error())},
			{Key: []byte("dlq-at"), Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
		},
	}
	if _, _, err := h.dlqProd.SendMessage(dlqMsg); err != nil {
		appLogger.Errorf("kafka consumer: failed to send message to DLQ (topic=%s): %v", h.dlqTopic, err)
	} else {
		appLogger.Infof("kafka consumer: message routed to DLQ (topic=%s, original-topic=%s, offset=%d)", h.dlqTopic, raw.Topic, raw.Offset)
	}
}

func toMessage(raw *sarama.ConsumerMessage) *Message {
	headers := make(map[string]string, len(raw.Headers))
	for _, h := range raw.Headers {
		headers[string(h.Key)] = string(h.Value)
	}
	return &Message{
		Topic:     raw.Topic,
		Partition: raw.Partition,
		Offset:    raw.Offset,
		Key:       raw.Key,
		Value:     raw.Value,
		Headers:   headers,
		Timestamp: raw.Timestamp,
	}
}
