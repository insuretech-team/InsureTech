package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// Publisher publishes DocGen events to Kafka
type Publisher struct {
	producer *producer.EventProducer
	topic    string
}

// NewPublisher creates a new Kafka event publisher for DocGen
func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("brokers list cannot be empty")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic cannot be empty")
	}

	ep, err := producer.NewEventProducer(brokers, topic, "docgen-service")
	if err != nil {
		return nil, fmt.Errorf("failed to create event producer: %w", err)
	}

	return &Publisher{
		producer: ep,
		topic:    topic,
	}, nil
}

// templateCreatedEvent represents a document.template.created event
type templateCreatedEvent struct {
	EventType    string    `json:"event_type"`
	Timestamp    time.Time `json:"timestamp"`
	TenantID     string    `json:"tenant_id"`
	TemplateID   string    `json:"template_id"`
	TemplateName string    `json:"template_name"`
}

// PublishTemplateCreated publishes a document.template.created event
func (p *Publisher) PublishTemplateCreated(ctx context.Context, templateID, tenantID, templateName string) error {
	if templateID == "" || tenantID == "" {
		return fmt.Errorf("templateID and tenantID are required")
	}

	event := templateCreatedEvent{
		EventType:    "document.template.created",
		Timestamp:    time.Now().UTC(),
		TenantID:     tenantID,
		TemplateID:   templateID,
		TemplateName: templateName,
	}

	err := p.producer.Produce(ctx, p.topic, tenantID, event)
	if err != nil {
		logger.Errorf("Failed to publish template.created event: %v", err)
		return fmt.Errorf("failed to publish template.created event: %w", err)
	}

	logger.Infof("Published template.created event (templateID=%s, tenantID=%s)", templateID, tenantID)
	return nil
}

// templateUpdatedEvent represents a document.template.updated event
type templateUpdatedEvent struct {
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
	TenantID   string    `json:"tenant_id"`
	TemplateID string    `json:"template_id"`
}

// PublishTemplateUpdated publishes a document.template.updated event
func (p *Publisher) PublishTemplateUpdated(ctx context.Context, templateID, tenantID string) error {
	if templateID == "" || tenantID == "" {
		return fmt.Errorf("templateID and tenantID are required")
	}

	event := templateUpdatedEvent{
		EventType:  "document.template.updated",
		Timestamp:  time.Now().UTC(),
		TenantID:   tenantID,
		TemplateID: templateID,
	}

	err := p.producer.Produce(ctx, p.topic, tenantID, event)
	if err != nil {
		logger.Errorf("Failed to publish template.updated event: %v", err)
		return fmt.Errorf("failed to publish template.updated event: %w", err)
	}

	logger.Infof("Published template.updated event (templateID=%s, tenantID=%s)", templateID, tenantID)
	return nil
}

// documentGeneratedEvent represents a document.generated event
type documentGeneratedEvent struct {
	EventType    string    `json:"event_type"`
	Timestamp    time.Time `json:"timestamp"`
	TenantID     string    `json:"tenant_id"`
	GenerationID string    `json:"generation_id"`
	TemplateID   string    `json:"template_id"`
	EntityID     string    `json:"entity_id"`
	EntityType   string    `json:"entity_type"`
	Format       string    `json:"format"`
}

// PublishDocumentGenerated publishes a document.generated event
func (p *Publisher) PublishDocumentGenerated(ctx context.Context, generationID, templateID, tenantID, entityID, entityType, format string) error {
	if generationID == "" || templateID == "" || tenantID == "" {
		return fmt.Errorf("generationID, templateID, and tenantID are required")
	}

	event := documentGeneratedEvent{
		EventType:    "document.generated",
		Timestamp:    time.Now().UTC(),
		TenantID:     tenantID,
		GenerationID: generationID,
		TemplateID:   templateID,
		EntityID:     entityID,
		EntityType:   entityType,
		Format:       format,
	}

	err := p.producer.Produce(ctx, p.topic, tenantID, event)
	if err != nil {
		logger.Errorf("Failed to publish document.generated event: %v", err)
		return fmt.Errorf("failed to publish document.generated event: %w", err)
	}

	logger.Infof("Published document.generated event (generationID=%s, templateID=%s, tenantID=%s)", generationID, templateID, tenantID)
	return nil
}

// generationFailedEvent represents a document.generation.failed event
type generationFailedEvent struct {
	EventType    string    `json:"event_type"`
	Timestamp    time.Time `json:"timestamp"`
	TenantID     string    `json:"tenant_id"`
	GenerationID string    `json:"generation_id"`
	Reason       string    `json:"reason"`
}

// PublishGenerationFailed publishes a document.generation.failed event
func (p *Publisher) PublishGenerationFailed(ctx context.Context, generationID, tenantID, reason string) error {
	if generationID == "" || tenantID == "" {
		return fmt.Errorf("generationID and tenantID are required")
	}

	event := generationFailedEvent{
		EventType:    "document.generation.failed",
		Timestamp:    time.Now().UTC(),
		TenantID:     tenantID,
		GenerationID: generationID,
		Reason:       reason,
	}

	err := p.producer.Produce(ctx, p.topic, tenantID, event)
	if err != nil {
		logger.Errorf("Failed to publish generation.failed event: %v", err)
		return fmt.Errorf("failed to publish generation.failed event: %w", err)
	}

	logger.Infof("Published generation.failed event (generationID=%s, tenantID=%s)", generationID, tenantID)
	return nil
}

// Close closes the Kafka publisher gracefully
func (p *Publisher) Close() error {
	if p.producer != nil {
		if err := p.producer.Close(); err != nil {
			logger.Errorf("Failed to close Kafka producer: %v", err)
			return fmt.Errorf("failed to close producer: %w", err)
		}
		logger.Info("Kafka producer closed gracefully")
	}
	return nil
}
