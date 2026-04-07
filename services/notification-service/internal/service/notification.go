package service

import (
	"context"
	"log"
)

type Notifier interface {
	Send(ctx context.Context, message string) error
}

type LogNotifier struct{}

func (l *LogNotifier) Send(ctx context.Context, message string) error {
	log.Println("[Notification]", message)
	return nil
}
