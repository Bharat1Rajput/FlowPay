package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusCreated    PaymentStatus = "CREATED"
	StatusProcessing PaymentStatus = "PROCESSING"
	StatusSuccess    PaymentStatus = "SUCCESS"
	StatusFailed     PaymentStatus = "FAILED"
	StatusInvalid    PaymentStatus = "INVALID"
)

func (s PaymentStatus) IsTerminal() bool {
	return s == StatusSuccess || s == StatusFailed || s == StatusInvalid
}

type Payment struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	IdempotencyKey string
	Amount         int64
	Currency       string
	Status         PaymentStatus
	GatewayRef     string
	FailureReason  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
