package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	if templateID == "" {
		return fmt.Errorf("templateID is required")
	}

	event := templateCreatedEvent{
		EventType:    "document.template.created",
		Timestamp:    time.Now().UTC(),
		TenantID:     tenantID,
		TemplateID:   templateID,
		TemplateName: templateName,
	}

	key := tenantID
	if key == "" {
		key = templateID
	}
	err := p.producer.Produce(ctx, p.topic, key, event)
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
	if templateID == "" {
		return fmt.Errorf("templateID is required")
	}

	event := templateUpdatedEvent{
		EventType:  "document.template.updated",
		Timestamp:  time.Now().UTC(),
		TenantID:   tenantID,
		TemplateID: templateID,
	}

	key := tenantID
	if key == "" {
		key = templateID
	}
	err := p.producer.Produce(ctx, p.topic, key, event)
	if err != nil {
		logger.Errorf("Failed to publish template.updated event: %v", err)
		return fmt.Errorf("failed to publish template.updated event: %w", err)
	}

	logger.Infof("Published template.updated event (templateID=%s, tenantID=%s)", templateID, tenantID)
	return nil
}

// documentGenerationRequestedEvent aligns with document.events.v1.DocumentGenerationRequestedEvent.
type documentGenerationRequestedEvent struct {
	EventID              string    `json:"event_id"`
	DocumentGenerationID string    `json:"document_generation_id"`
	DocumentTemplateID   string    `json:"document_template_id"`
	EntityType           string    `json:"entity_type"`
	EntityID             string    `json:"entity_id"`
	CorrelationID        string    `json:"correlation_id"`
	Timestamp            time.Time `json:"timestamp"`
}

// PublishGenerationRequested publishes a document generation requested event.
func (p *Publisher) PublishGenerationRequested(ctx context.Context, generationID, templateID, tenantID, entityType, entityID, correlationID string) error {
	if generationID == "" || templateID == "" {
		return fmt.Errorf("generationID and templateID are required")
	}

	event := documentGenerationRequestedEvent{
		EventID:              uuid.NewString(),
		DocumentGenerationID: generationID,
		DocumentTemplateID:   templateID,
		EntityType:           entityType,
		EntityID:             entityID,
		CorrelationID:        correlationID,
		Timestamp:            time.Now().UTC(),
	}

	key := tenantID
	if key == "" {
		key = generationID
	}
	err := p.producer.Produce(ctx, p.topic, key, event)
	if err != nil {
		logger.Errorf("Failed to publish document generation requested event: %v", err)
		return fmt.Errorf("failed to publish document generation requested event: %w", err)
	}

	logger.Infof("Published document generation requested event (generationID=%s, templateID=%s)", generationID, templateID)
	return nil
}

// documentGeneratedEvent aligns with document.events.v1.DocumentGeneratedEvent.
type documentGeneratedEvent struct {
	EventID              string    `json:"event_id"`
	DocumentGenerationID string    `json:"document_generation_id"`
	EntityType           string    `json:"entity_type"`
	EntityID             string    `json:"entity_id"`
	FileURL              string    `json:"file_url"`
	CorrelationID        string    `json:"correlation_id"`
	Timestamp            time.Time `json:"timestamp"`
}

// PublishDocumentGenerated publishes a document generated event.
func (p *Publisher) PublishDocumentGenerated(ctx context.Context, generationID, tenantID, entityID, entityType, fileURL, correlationID string) error {
	if generationID == "" {
		return fmt.Errorf("generationID is required")
	}

	event := documentGeneratedEvent{
		EventID:              uuid.NewString(),
		DocumentGenerationID: generationID,
		EntityType:           entityType,
		EntityID:             entityID,
		FileURL:              fileURL,
		CorrelationID:        correlationID,
		Timestamp:            time.Now().UTC(),
	}

	key := tenantID
	if key == "" {
		key = generationID
	}
	err := p.producer.Produce(ctx, p.topic, key, event)
	if err != nil {
		logger.Errorf("Failed to publish document generated event: %v", err)
		return fmt.Errorf("failed to publish document generated event: %w", err)
	}

	logger.Infof("Published document generated event (generationID=%s)", generationID)
	return nil
}

// generationFailedEvent aligns with document.events.v1.DocumentGenerationFailedEvent.
type generationFailedEvent struct {
	EventID              string    `json:"event_id"`
	DocumentGenerationID string    `json:"document_generation_id"`
	ErrorMessage         string    `json:"error_message"`
	CorrelationID        string    `json:"correlation_id"`
	Timestamp            time.Time `json:"timestamp"`
}

// PublishGenerationFailed publishes a document generation failed event.
func (p *Publisher) PublishGenerationFailed(ctx context.Context, generationID, tenantID, errorMessage, correlationID string) error {
	if generationID == "" {
		return fmt.Errorf("generationID is required")
	}

	event := generationFailedEvent{
		EventID:              uuid.NewString(),
		DocumentGenerationID: generationID,
		ErrorMessage:         errorMessage,
		CorrelationID:        correlationID,
		Timestamp:            time.Now().UTC(),
	}

	key := tenantID
	if key == "" {
		key = generationID
	}
	err := p.producer.Produce(ctx, p.topic, key, event)
	if err != nil {
		logger.Errorf("Failed to publish document generation failed event: %v", err)
		return fmt.Errorf("failed to publish document generation failed event: %w", err)
	}

	logger.Infof("Published document generation failed event (generationID=%s)", generationID)
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
