package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Bharat1Rajput/flowpay/services/notification-service/internal/service"
	"github.com/IBM/sarama"
)

type PaymentEvent struct {
	EventID   string  `json:"event_id"`
	EventType string  `json:"event_type"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	OrderID string `json:"order_id"`
}

type NotificationConsumer struct {
	retry *RetryHandler
	dlq   *DLQProducer
	svc   *service.NotificationService
}

func NewNotificationConsumer(retry *RetryHandler, dlq *DLQProducer, svc *service.NotificationService) *NotificationConsumer {
	return &NotificationConsumer{
		retry: retry,
		dlq:   dlq,
		svc:   svc,
	}
}

// required for sarama
func (c *NotificationConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *NotificationConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *NotificationConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {

		log.Println("notification received:", string(msg.Value))

		var event PaymentEvent

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid message:", err)
			session.MarkMessage(msg, "")
			continue
		}

		err := c.retry.Execute(session.Context(), func() error {
			return c.handleEvent(session.Context(), event)
		})

		if err != nil {
			log.Println("notification failed after retries, sending to DLQ:", err)

			if dlqErr := c.dlq.Send(session.Context(), msg, err.Error()); dlqErr != nil {
				log.Println("failed to send to DLQ:", dlqErr)
			}

			session.MarkMessage(msg, "")
			continue
		}

		session.MarkMessage(msg, "")

	}

	return nil
}

func (c *NotificationConsumer) handleEvent(ctx context.Context, event PaymentEvent) error {

	switch event.EventType {

	case "payment.succeeded":
		return c.svc.SendPaymentSuccess(ctx, event.Payload.OrderID)

	case "payment.failed":
		return c.svc.SendPaymentFailure(ctx, event.Payload.OrderID)

	default:
		return fmt.Errorf("unknown event: %s", event.EventType)
	}
}
