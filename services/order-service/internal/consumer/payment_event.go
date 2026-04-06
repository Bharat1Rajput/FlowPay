package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/service"
)

type PaymentEvent struct {
	EventID   string  `json:"event_id"`
	EventType string  `json:"event_type"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	OrderID string `json:"order_id"`
}

type PaymentConsumer struct {
	svc *service.OrderService
}

func NewPaymentConsumer(svc *service.OrderService) *PaymentConsumer {
	return &PaymentConsumer{svc: svc}
}

// REQUIRED for sarama
func (c *PaymentConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *PaymentConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *PaymentConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {

		log.Println("received message:", string(msg.Value))

		var event PaymentEvent

		// parse JSON
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid message:", err)
			session.MarkMessage(msg, "")
			continue
		}

		// process event
		if err := c.handleEvent(session.Context(), event); err != nil {
			log.Println("error processing event:", err)
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (c *PaymentConsumer) handleEvent(ctx context.Context, event PaymentEvent) error {

	orderID, err := uuid.Parse(event.Payload.OrderID)
	if err != nil {
		return err
	}
	log.Printf("processing event %s for order %s\n", event.EventType, orderID)

	switch event.EventType {

	case "payment.succeeded":
		return c.svc.ConfirmFromPayment(ctx, orderID)

	case "payment.failed":
		return c.svc.CancelFromPayment(ctx, orderID)

	default:
		log.Println("unknown event:", event.EventType)
		return nil
	}
}
