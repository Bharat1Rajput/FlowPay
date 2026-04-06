package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/model"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(p sarama.SyncProducer, topic string) *KafkaProducer {
	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}
}

type PaymentEvent struct {
	EventID     string    `json:"event_id"`
	EventType   string    `json:"event_type"`
	PublishedAt time.Time `json:"published_at"`
	Version     string    `json:"version"`
	Payload     Payload   `json:"payload"`
}

type Payload struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
}

func (k *KafkaProducer) PublishPaymentEvent(ctx context.Context, payment *model.Payment) error {

	eventType := "payment.failed"
	if payment.Status == model.StatusSuccess {
		eventType = "payment.succeeded"
	}

	event := PaymentEvent{
		EventID:     uuid.New().String(),
		EventType:   eventType,
		PublishedAt: time.Now().UTC(),
		Version:     "1",
		Payload: Payload{
			PaymentID: payment.ID.String(),
			OrderID:   payment.OrderID.String(),
			Amount:    payment.Amount,
			Currency:  payment.Currency,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	} 

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(payment.OrderID.String()), // 🔥 important
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = k.producer.SendMessage(msg)
	return err
}
