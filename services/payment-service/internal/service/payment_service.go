package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/kafka"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/model"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/repository"
)

type PaymentService struct {
	repo     repository.PaymentRepository
	producer kafka.EventProducer
}

func NewPaymentService(repo repository.PaymentRepository, producer kafka.EventProducer) *PaymentService {
	return &PaymentService{repo: repo, producer: producer}
}

type ProcessPaymentInput struct {
	OrderID        uuid.UUID
	IdempotencyKey string
	Amount         int64
	Currency       string
}

func (s *PaymentService) ProcessPayment(ctx context.Context, in ProcessPaymentInput) (*model.Payment, error) {

	// Step 1: Idempotency check
	existing, err := s.repo.GetByIdempotencyKey(ctx, in.IdempotencyKey)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// Step 2: Validate
	if in.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	// Step 3: Create payment
	payment := &model.Payment{
		ID:             uuid.New(),
		OrderID:        in.OrderID,
		IdempotencyKey: in.IdempotencyKey,
		Amount:         in.Amount,
		Currency:       in.Currency,
		Status:         model.StatusCreated,
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	// Step 4: Move to PROCESSING
	_ = s.repo.UpdateStatus(ctx, payment.ID, model.StatusCreated, model.StatusProcessing, "", "")
	payment.Status = model.StatusProcessing

	// Step 5: Mock gateway call
	gatewayRef := "mock_txn_" + payment.ID.String()

	// Step 6: Final state
	err = s.repo.UpdateStatus(
		ctx,
		payment.ID,
		model.StatusProcessing,
		model.StatusSuccess,
		gatewayRef,
		"",
	)
	if err != nil {
		return nil, err
	}

	payment.Status = model.StatusSuccess
	payment.GatewayRef = gatewayRef
	// Step 7: Publish event
	if err := s.producer.PublishPaymentEvent(ctx, payment); err != nil {
		// log but DO NOT fail request
		fmt.Printf("failed to publish event: %v\n", err)
	}
	return payment, nil
}
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	return s.repo.GetByID(ctx, id)
}
 