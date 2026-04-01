package kafkaapp

import (
	"context"

	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
)

type ManagedConsumer struct {
	group  *kafkaconsumer.ConsumerGroup
	cancel context.CancelFunc
}

func StartConsumerGroup(cfg kafkaconsumer.Config) (*ManagedConsumer, error) {
	group, err := kafkaconsumer.NewConsumerGroup(cfg)
	if err != nil {
		return nil, err
	}
	consumerCtx, cancel := context.WithCancel(context.Background())
	go group.Start(consumerCtx)
	return &ManagedConsumer{
		group:  group,
		cancel: cancel,
	}, nil
}

func (m *ManagedConsumer) Close() error {
	if m == nil {
		return nil
	}
	if m.cancel != nil {
		m.cancel()
	}
	if m.group != nil {
		return m.group.Close()
	}
	return nil
}

func (m *ManagedConsumer) Group() *kafkaconsumer.ConsumerGroup {
	if m == nil {
		return nil
	}
	return m.group
}
