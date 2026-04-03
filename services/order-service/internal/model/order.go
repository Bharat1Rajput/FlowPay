package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending        OrderStatus = "PENDING"
	StatusConfirmed      OrderStatus = "CONFIRMED"
	StatusPreparing      OrderStatus = "PREPARING"
	StatusOutForDelivery OrderStatus = "OUT_FOR_DELIVERY"
	StatusDelivered      OrderStatus = "DELIVERED"
	StatusCancelled      OrderStatus = "CANCELLED"
)

// IsTerminal → once reached, no further transitions allowed
func (s OrderStatus) IsTerminal() bool {
	return s == StatusDelivered || s == StatusCancelled
}

// CanUserCancel → business rule
func (s OrderStatus) CanUserCancel() bool {
	return s == StatusPending || s == StatusConfirmed
}

type Order struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Status       OrderStatus
	Items        []OrderItem
	TotalAmount  int64 // ALWAYS in paise
	Currency     string
	DeliveryAddr string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ComputeTotal ensures we NEVER trust client total
func (o *Order) ComputeTotal() {
	var total int64
	for _, item := range o.Items {
		total += item.TotalPrice
	}
	o.TotalAmount = total
}
