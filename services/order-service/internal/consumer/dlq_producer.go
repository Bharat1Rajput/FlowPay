package consumer

import (
	"context"

	"github.com/IBM/sarama"
)

type DLQProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewDLQProducer(p sarama.SyncProducer, topic string) *DLQProducer {
	return &DLQProducer{
		producer: p,
		topic:    topic,
	}
}

func (d *DLQProducer) Send(ctx context.Context, msg *sarama.ConsumerMessage, reason string) error {

	dlqMsg := &sarama.ProducerMessage{
		Topic: d.topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}

	_, _, err := d.producer.SendMessage(dlqMsg)
	return err
}
