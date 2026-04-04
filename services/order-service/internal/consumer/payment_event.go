package consumer

import (
	"context"
	"encoding/json"

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

func (c *PaymentConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {

		var event PaymentEvent

		// 1. Parse JSON
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			// skip bad message
			session.MarkMessage(msg, "")
			continue
		}

		// 2. Handle event
		c.handleEvent(session.Context(), event)

		// 3. Mark as processed
		session.MarkMessage(msg, "")
	}

	return nil
}

func (c *PaymentConsumer) handleEvent(ctx context.Context, event PaymentEvent) {

	orderID, err := uuid.Parse(event.Payload.OrderID)
	if err != nil {
		return
	}

	switch event.EventType {

	case "payment.succeeded":
		_ = c.svc.ConfirmFromPayment(ctx, orderID)

	case "payment.failed":
		_ = c.svc.CancelOrder(ctx, orderID, uuid.Nil) // system cancel

	default:
		// ignore unknown events

	}
}
