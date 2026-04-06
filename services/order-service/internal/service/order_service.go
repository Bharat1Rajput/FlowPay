package service

import (
	"context"
	"fmt"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/repository"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/statemachine"
	"github.com/google/uuid"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

type CreateOrderInput struct {
	UserID       uuid.UUID
	DeliveryAddr string
	Notes        string
	Items        []model.OrderItem
}

func (s *OrderService) CreateOrder(ctx context.Context, in CreateOrderInput) (*model.Order, error) {

	// 1. Validation
	fmt.Println("items : ", in.Items)
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("items cannot be empty")
	}
	if in.DeliveryAddr == "" {
		return nil, fmt.Errorf("delivery address required")
	}

	// 2. Build order
	order := &model.Order{
		ID:           uuid.New(),
		UserID:       in.UserID,
		Status:       model.StatusPending,
		Items:        in.Items,
		Currency:     "INR",
		DeliveryAddr: in.DeliveryAddr,
		Notes:        in.Notes,
	}

	// 3. Compute total (CRITICAL)
	order.ComputeTotal()

	// 4. Persist
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Ownership check
	if order.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Business rule
	if !order.Status.CanUserCancel() {
		return fmt.Errorf("cannot cancel at this stage")
	}

	// State machine validation
	if err := statemachine.CanTransition(order.Status, model.StatusCancelled); err != nil {
		return err
	}

	// Update with optimistic locking
	return s.repo.UpdateStatus(ctx, id, order.Status, model.StatusCancelled)
}

func (s *OrderService) ConfirmFromPayment(ctx context.Context, orderID uuid.UUID) error {

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Idempotency guard
	if order.Status == model.StatusConfirmed {
		return nil
	}

	// Validate transition
	if err := statemachine.CanTransition(order.Status, model.StatusConfirmed); err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, orderID, order.Status, model.StatusConfirmed)
}

func (s *OrderService) CancelFromPayment(ctx context.Context, orderID uuid.UUID) error {

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Idempotency (critical for Kafka retries)
	if order.Status == model.StatusCancelled {
		return nil
	}

	// Validate transition
	if err := statemachine.CanTransition(order.Status, model.StatusCancelled); err != nil {
		return err
	}

	// Update with optimistic locking
	return s.repo.UpdateStatus(ctx, orderID, order.Status, model.StatusCancelled)
}
