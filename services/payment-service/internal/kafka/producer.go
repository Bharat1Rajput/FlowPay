package kafka

import (
	"context"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/model"
)

type EventProducer interface {
	PublishPaymentEvent(ctx context.Context, payment *model.Payment) error
}
