package service

import (
	"context"
	"fmt"
)

type NotificationService struct {
	notifier Notifier
}

func NewNotificationService(n Notifier) *NotificationService {
	return &NotificationService{
		notifier: n,
	}

}

func (s *NotificationService) SendPaymentSuccess(ctx context.Context, orderID string) error {

	msg := fmt.Sprintf("Payment successful for order %s", orderID)

	return s.notifier.Send(ctx, msg)
}

func (s *NotificationService) SendPaymentFailure(ctx context.Context, orderID string) error {

	msg := fmt.Sprintf("Payment failed for order %s", orderID)

	return s.notifier.Send(ctx, msg)
}


