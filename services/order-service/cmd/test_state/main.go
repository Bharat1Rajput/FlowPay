package main

import (
	"fmt"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/statemachine"
)

func main() {

	// ✅ VALID transition
	err := statemachine.CanTransition(
		model.StatusPending,
		model.StatusConfirmed,
	)

	fmt.Println("Pending → Confirmed:", err)

	// ❌ INVALID transition
	err = statemachine.CanTransition(
		model.StatusDelivered,
		model.StatusPreparing,
	)

	fmt.Println("Delivered → Preparing:", err)
}
