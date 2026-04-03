package main

import (
	"fmt"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
	"github.com/google/uuid"
)

func main() {
	orderID := uuid.New()

	item := model.NewOrderItem(orderID, "Pizza", 2, 5000)

	order := model.Order{
		ID:     orderID,
		UserID: uuid.New(),
		Status: model.StatusPending,
		Items:  []model.OrderItem{item},
	}

	order.ComputeTotal()

	fmt.Println("Total:", order.TotalAmount)
}
