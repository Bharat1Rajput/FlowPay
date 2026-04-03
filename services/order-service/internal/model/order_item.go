package model

import "github.com/google/uuid"

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	ItemName   string
	Quantity   int
	UnitPrice  int64
	TotalPrice int64 // Quantity * UnitPrice
}

// Constructor → ensures invariant
func NewOrderItem(orderID uuid.UUID, name string, qty int, unitPrice int64) OrderItem {
	return OrderItem{
		ID:         uuid.New(),
		OrderID:    orderID,
		ItemName:   name,
		Quantity:   qty,
		UnitPrice:  unitPrice,
		TotalPrice: int64(qty) * unitPrice,
	}
}
