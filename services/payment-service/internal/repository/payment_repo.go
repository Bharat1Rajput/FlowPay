package repository

import (
	"context"
	"fmt"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/model"
	"github.com/google/uuid"
)

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrStaleUpdate = fmt.Errorf("stale update")
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, p *model.Payment) error
	GetByIdempotencyKey(ctx context.Context, key string) (*model.Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, expectedCurrent, newStatus model.PaymentStatus, gatewayRef, failureReason string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Payment, error)
}
