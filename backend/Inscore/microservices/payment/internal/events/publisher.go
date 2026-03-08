package events

import (
    "context"

    kafkaproducer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
    appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
    "google.golang.org/protobuf/proto"
)

type Publisher struct {
    producer *kafkaproducer.EventProducer
}

func NewPublisher(producer *kafkaproducer.EventProducer) *Publisher {
    return &Publisher{producer: producer}
}

func (p *Publisher) Publish(ctx context.Context, topic string, key string, msg proto.Message) {
    if p == nil || p.producer == nil {
        appLogger.Warnf("Kafka producer not available - skipping event on topic %s", topic)
        return
    }
    if err := p.producer.Produce(ctx, topic, key, msg); err != nil {
        appLogger.Errorf("failed to publish event to topic %s: %v", topic, err)
    }
}
