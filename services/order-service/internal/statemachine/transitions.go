package statemachine

import (
	"fmt"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
)

// Define allowed transitions
var allowedTransitions = map[model.OrderStatus][]model.OrderStatus{
	model.StatusPending: {
		model.StatusConfirmed,
		model.StatusCancelled,
	},
	model.StatusConfirmed: {
		model.StatusPreparing,
		model.StatusCancelled,
	},
	model.StatusPreparing: {
		model.StatusOutForDelivery,
	},
	model.StatusOutForDelivery: {
		model.StatusDelivered,
	},
}

// Custom error type
type InvalidTransitionError struct {
	From model.OrderStatus
	To   model.OrderStatus
}

func (e InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid transition from %s to %s", e.From, e.To)
}

// Core function
func CanTransition(from, to model.OrderStatus) error {
	// Terminal states cannot transition
	if from.IsTerminal() {
		return InvalidTransitionError{From: from, To: to}
	}

	allowed, exists := allowedTransitions[from]
	if !exists {
		return InvalidTransitionError{From: from, To: to}
	}

	for _, s := range allowed {
		if s == to {
			return nil
		}
	}

	return InvalidTransitionError{From: from, To: to}
}
