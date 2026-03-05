package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// Publisher handles publishing media events to Kafka
type Publisher struct {
	producer *producer.EventProducer
	topic    string
}

// NewPublisher creates a new Kafka event publisher for media events
func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("at least one Kafka broker must be provided")
	}
	if topic == "" {
		return nil, fmt.Errorf("Kafka topic must be provided")
	}

	eventProducer, err := producer.NewEventProducer(brokers, topic, "media-service")
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	logger.Infof("Media Kafka publisher initialized (brokers=%v, topic=%s)", brokers, topic)

	return &Publisher{
		producer: eventProducer,
		topic:    topic,
	}, nil
}

// mediaEvent represents the base structure for all media events
type mediaEvent struct {
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
	TenantID   string    `json:"tenant_id"`
	MediaID    string    `json:"media_id,omitempty"`
	JobID      string    `json:"job_id,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// PublishFileUploaded publishes a media.file.uploaded event
func (p *Publisher) PublishFileUploaded(ctx context.Context, mediaID, tenantID, filename, mimeType string, sizeBytes int64) error {
	if mediaID == "" || tenantID == "" {
		return fmt.Errorf("mediaID and tenantID are required")
	}

	event := mediaEvent{
		EventType: "media.file.uploaded",
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		MediaID:   mediaID,
		Data: map[string]interface{}{
			"filename":   filename,
			"mime_type":  mimeType,
			"size_bytes": sizeBytes,
		},
	}

	if err := p.producer.Produce(ctx, p.topic, tenantID, event); err != nil {
		logger.Errorf("Failed to publish media.file.uploaded event (mediaID=%s, tenantID=%s): %v", mediaID, tenantID, err)
		return fmt.Errorf("failed to publish media.file.uploaded event: %w", err)
	}

	logger.Infof("Published media.file.uploaded event (mediaID=%s, tenantID=%s)", mediaID, tenantID)
	return nil
}

// PublishFileDeleted publishes a media.file.deleted event
func (p *Publisher) PublishFileDeleted(ctx context.Context, mediaID, tenantID string) error {
	if mediaID == "" || tenantID == "" {
		return fmt.Errorf("mediaID and tenantID are required")
	}

	event := mediaEvent{
		EventType: "media.file.deleted",
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		MediaID:   mediaID,
	}

	if err := p.producer.Produce(ctx, p.topic, tenantID, event); err != nil {
		logger.Errorf("Failed to publish media.file.deleted event (mediaID=%s, tenantID=%s): %v", mediaID, tenantID, err)
		return fmt.Errorf("failed to publish media.file.deleted event: %w", err)
	}

	logger.Infof("Published media.file.deleted event (mediaID=%s, tenantID=%s)", mediaID, tenantID)
	return nil
}

// PublishProcessingStarted publishes a media.processing.started event
func (p *Publisher) PublishProcessingStarted(ctx context.Context, jobID, mediaID, tenantID, jobType string) error {
	if jobID == "" || mediaID == "" || tenantID == "" {
		return fmt.Errorf("jobID, mediaID, and tenantID are required")
	}

	event := mediaEvent{
		EventType: "media.processing.started",
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		MediaID:   mediaID,
		JobID:     jobID,
		Data: map[string]interface{}{
			"job_type": jobType,
		},
	}

	if err := p.producer.Produce(ctx, p.topic, tenantID, event); err != nil {
		logger.Errorf("Failed to publish media.processing.started event (jobID=%s, mediaID=%s, tenantID=%s): %v", jobID, mediaID, tenantID, err)
		return fmt.Errorf("failed to publish media.processing.started event: %w", err)
	}

	logger.Infof("Published media.processing.started event (jobID=%s, mediaID=%s, tenantID=%s)", jobID, mediaID, tenantID)
	return nil
}

// PublishProcessingCompleted publishes a media.processing.completed event
func (p *Publisher) PublishProcessingCompleted(ctx context.Context, jobID, mediaID, tenantID, jobType string) error {
	if jobID == "" || mediaID == "" || tenantID == "" {
		return fmt.Errorf("jobID, mediaID, and tenantID are required")
	}

	event := mediaEvent{
		EventType: "media.processing.completed",
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		MediaID:   mediaID,
		JobID:     jobID,
		Data: map[string]interface{}{
			"job_type": jobType,
		},
	}

	if err := p.producer.Produce(ctx, p.topic, tenantID, event); err != nil {
		logger.Errorf("Failed to publish media.processing.completed event (jobID=%s, mediaID=%s, tenantID=%s): %v", jobID, mediaID, tenantID, err)
		return fmt.Errorf("failed to publish media.processing.completed event: %w", err)
	}

	logger.Infof("Published media.processing.completed event (jobID=%s, mediaID=%s, tenantID=%s)", jobID, mediaID, tenantID)
	return nil
}

// PublishProcessingFailed publishes a media.processing.failed event
func (p *Publisher) PublishProcessingFailed(ctx context.Context, jobID, mediaID, tenantID, jobType, reason string) error {
	if jobID == "" || mediaID == "" || tenantID == "" {
		return fmt.Errorf("jobID, mediaID, and tenantID are required")
	}

	event := mediaEvent{
		EventType: "media.processing.failed",
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		MediaID:   mediaID,
		JobID:     jobID,
		Data: map[string]interface{}{
			"job_type": jobType,
			"reason":   reason,
		},
	}

	if err := p.producer.Produce(ctx, p.topic, tenantID, event); err != nil {
		logger.Errorf("Failed to publish media.processing.failed event (jobID=%s, mediaID=%s, tenantID=%s): %v", jobID, mediaID, tenantID, err)
		return fmt.Errorf("failed to publish media.processing.failed event: %w", err)
	}

	logger.Infof("Published media.processing.failed event (jobID=%s, mediaID=%s, tenantID=%s, reason=%s)", jobID, mediaID, tenantID, reason)
	return nil
}

// PublishVirusDetected publishes a media.virus.detected event
func (p *Publisher) PublishVirusDetected(ctx context.Context, mediaID, tenantID, virusName string) error {
	if mediaID == "" || tenantID == "" {
		return fmt.Errorf("mediaID and tenantID are required")
	}

	event := mediaEvent{
		EventType: "media.virus.detected",
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		MediaID:   mediaID,
		Data: map[string]interface{}{
			"virus_name": virusName,
		},
	}

	if err := p.producer.Produce(ctx, p.topic, tenantID, event); err != nil {
		logger.Errorf("Failed to publish media.virus.detected event (mediaID=%s, tenantID=%s): %v", mediaID, tenantID, err)
		return fmt.Errorf("failed to publish media.virus.detected event: %w", err)
	}

	logger.Infof("Published media.virus.detected event (mediaID=%s, tenantID=%s, virus=%s)", mediaID, tenantID, virusName)
	return nil
}

// Close gracefully closes the Kafka publisher
func (p *Publisher) Close() error {
	if p.producer != nil {
		if err := p.producer.Close(); err != nil {
			logger.Errorf("Failed to close media Kafka publisher: %v", err)
			return fmt.Errorf("failed to close kafka producer: %w", err)
		}
		logger.Info("Media Kafka publisher closed gracefully")
	}
	return nil
}
