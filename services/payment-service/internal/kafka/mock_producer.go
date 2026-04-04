package kafka

import (
	"context"
	"log"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/model"
)

type MockProducer struct{}

func (m *MockProducer) PublishPaymentEvent(ctx context.Context, p *model.Payment) error {
	log.Println("Mock Kafka Event Published:", p.ID)
	return nil
}
