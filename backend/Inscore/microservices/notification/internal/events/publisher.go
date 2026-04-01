package events

import (
	"context"
	"reflect"
	"time"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	notificationeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, msg interface{}) error
	Close() error
}

type Publisher struct {
	producer EventProducer
}

const publishTimeout = 500 * time.Millisecond

func NewPublisher(producer EventProducer) *Publisher {
	if isNilProducer(producer) {
		producer = nil
	}
	return &Publisher{producer: producer}
}

func (p *Publisher) PublishNotificationSent(ctx context.Context, notificationID, recipientID, channel, notificationType, correlationID string) error {
	topic, _ := DescriptorForKey("notification_sent")
	evt := &notificationeventsv1.NotificationSentEvent{
		EventId:        uuid.NewString(),
		NotificationId: notificationID,
		RecipientId:    recipientID,
		Channel:        channel,
		Type:           notificationType,
		Timestamp:      timestamppb.New(time.Now()),
		CorrelationId:  correlationID,
	}
	if err := p.publish(ctx, topic.Topic, notificationID, evt); err != nil {
		appLogger.Warnf("Failed to publish NotificationSentEvent for notification %s: %v", notificationID, err)
	}
	return nil
}

func (p *Publisher) PublishNotificationDelivered(ctx context.Context, notificationID, recipientID, correlationID string) error {
	topic, _ := DescriptorForKey("notification_delivered")
	evt := &notificationeventsv1.NotificationDeliveredEvent{
		EventId:        uuid.NewString(),
		NotificationId: notificationID,
		RecipientId:    recipientID,
		Timestamp:      timestamppb.New(time.Now()),
		CorrelationId:  correlationID,
	}
	if err := p.publish(ctx, topic.Topic, notificationID, evt); err != nil {
		appLogger.Warnf("Failed to publish NotificationDeliveredEvent for notification %s: %v", notificationID, err)
	}
	return nil
}

func (p *Publisher) PublishNotificationFailed(ctx context.Context, notificationID, recipientID, errorMessage, correlationID string) error {
	topic, _ := DescriptorForKey("notification_failed")
	evt := &notificationeventsv1.NotificationFailedEvent{
		EventId:        uuid.NewString(),
		NotificationId: notificationID,
		RecipientId:    recipientID,
		ErrorMessage:   errorMessage,
		Timestamp:      timestamppb.New(time.Now()),
		CorrelationId:  correlationID,
	}
	if err := p.publish(ctx, topic.Topic, notificationID, evt); err != nil {
		appLogger.Warnf("Failed to publish NotificationFailedEvent for notification %s: %v", notificationID, err)
	}
	return nil
}

func (p *Publisher) publish(ctx context.Context, topic, key string, msg interface{}) error {
	if p == nil || isNilProducer(p.producer) {
		appLogger.Infof("Kafka producer not configured - event dropped (topic=%s, key=%s)", topic, key)
		return nil
	}
	if ctx != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	publishCtx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- p.producer.Produce(publishCtx, topic, key, msg)
	}()

	select {
	case err := <-errCh:
		return err
	case <-publishCtx.Done():
		return publishCtx.Err()
	}
}

func isNilProducer(producer EventProducer) bool {
	if producer == nil {
		return true
	}
	value := reflect.ValueOf(producer)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
