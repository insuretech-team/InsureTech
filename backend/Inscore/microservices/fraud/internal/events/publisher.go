package events

import (
	"context"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	fraudeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventProducer abstracts Kafka event publication.
type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, msg interface{}) error
	Close() error
}

// Publisher publishes fraud domain events to Kafka.
type Publisher struct {
	producer EventProducer
	topic    string
}

func NewPublisher(producer EventProducer, topic string) *Publisher {
	return &Publisher{producer: producer, topic: topic}
}

func (p *Publisher) publish(ctx context.Context, key string, msg interface{}) error {
	topic := "fraud.events"
	if p != nil && p.topic != "" {
		topic = p.topic
	}
	if p == nil || p.producer == nil {
		appLogger.Infof("fraud: Kafka producer not configured, dropping event (topic=%s)", topic)
		return nil
	}
	return p.producer.Produce(ctx, topic, key, msg)
}

func (p *Publisher) PublishFraudAlertTriggered(ctx context.Context, alert *fraudv1.FraudAlert, correlationID string) error {
	if alert == nil {
		return nil
	}
	evt := &fraudeventsv1.FraudAlertTriggeredEvent{
		EventId:       uuid.NewString(),
		FraudAlertId:  alert.Id,
		AlertNumber:   alert.AlertNumber,
		EntityType:    alert.EntityType,
		EntityId:      alert.EntityId,
		RiskLevel:     alert.RiskLevel,
		FraudScore:    alert.FraudScore,
		CorrelationId: correlationID,
		Timestamp:     timestamppb.Now(),
	}
	if err := p.publish(ctx, alert.Id, evt); err != nil {
		appLogger.Warnf("fraud: failed to publish FraudAlertTriggeredEvent (alert_id=%s): %v", alert.Id, err)
	}
	return nil
}

func (p *Publisher) PublishFraudCaseCreated(ctx context.Context, fraudCase *fraudv1.FraudCase, correlationID string) error {
	if fraudCase == nil {
		return nil
	}
	evt := &fraudeventsv1.FraudCaseCreatedEvent{
		EventId:        uuid.NewString(),
		FraudCaseId:    fraudCase.Id,
		CaseNumber:     fraudCase.CaseNumber,
		FraudAlertId:   fraudCase.FraudAlertId,
		Priority:       fraudCase.Priority.String(),
		InvestigatorId: fraudCase.InvestigatorId,
		CorrelationId:  correlationID,
		Timestamp:      timestamppb.Now(),
	}
	if err := p.publish(ctx, fraudCase.Id, evt); err != nil {
		appLogger.Warnf("fraud: failed to publish FraudCaseCreatedEvent (case_id=%s): %v", fraudCase.Id, err)
	}
	return nil
}

func (p *Publisher) PublishFraudConfirmed(ctx context.Context, fraudCase *fraudv1.FraudCase, entityType string, entityID string, correlationID string) error {
	if fraudCase == nil {
		return nil
	}
	evt := &fraudeventsv1.FraudConfirmedEvent{
		EventId:       uuid.NewString(),
		FraudCaseId:   fraudCase.Id,
		EntityType:    entityType,
		EntityId:      entityID,
		CorrelationId: correlationID,
		Timestamp:     timestamppb.Now(),
	}
	if err := p.publish(ctx, fraudCase.Id, evt); err != nil {
		appLogger.Warnf("fraud: failed to publish FraudConfirmedEvent (case_id=%s): %v", fraudCase.Id, err)
	}
	return nil
}
