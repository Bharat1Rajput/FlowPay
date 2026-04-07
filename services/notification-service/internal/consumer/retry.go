package consumer

import (
	"context"
	"time"
)

type RetryHandler struct {
	maxAttempts int
	baseDelay   time.Duration
}

func NewRetryHandler() *RetryHandler {
	return &RetryHandler{
		maxAttempts: 3,
		baseDelay:   time.Second,
	}
}

func (r *RetryHandler) Execute(ctx context.Context, fn func() error) error {

	var err error

	for attempt := 0; attempt < r.maxAttempts; attempt++ {

		err = fn()
		if err == nil {
			return nil
		}

		delay := r.baseDelay * (1 << attempt) // 1s, 2s, 4s
		time.Sleep(delay)
	}

	return err
}
