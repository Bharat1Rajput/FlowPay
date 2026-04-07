package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"

	"github.com/Bharat1Rajput/flowpay/services/notification-service/internal/consumer"
	"github.com/Bharat1Rajput/flowpay/services/notification-service/internal/service"
)

func main() {

	ctx := context.Background()

	// ---------------- Kafka Config ----------------
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}
	var err error

	// ---------------- Producer (for DLQ) ----------------
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{broker}, producerConfig)
	if err != nil {
		log.Fatal("failed to create producer:", err)
	}

	// Retry + DLQ
	retryHandler := consumer.NewRetryHandler()
	dlqProducer := consumer.NewDLQProducer(producer, "notifications.DLQ")
	notifier := &service.LogNotifier{}
	notificationService := service.NewNotificationService(notifier)
	// ---------------- Consumer Group ----------------
	var consumerGroup sarama.ConsumerGroup

	for i := 0; i < 5; i++ {
		consumerGroup, err = sarama.NewConsumerGroup([]string{broker}, "notification-group", config)
		if err == nil {
			break
		}

		log.Println("retrying kafka connection...", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect to kafka:", err)
	}

	// ---------------- Consumer ----------------
	notificationConsumer := consumer.NewNotificationConsumer(
		retryHandler,
		dlqProducer,
		notificationService,
	)
	// ---------------- Start Consumer ----------------
	for {
		err := consumerGroup.Consume(ctx, []string{"payments_topic"}, notificationConsumer)
		if err != nil {
			log.Println("consumer error:", err)
		}
	}

}
