package repository

import (
	"context"
	"fmt"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
	"github.com/google/uuid"
)

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrStaleUpdate = fmt.Errorf("stale update")
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, expectedCurrent, newStatus model.OrderStatus) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
}
