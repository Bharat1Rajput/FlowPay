package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/Bharat1Rajput/flowpay/proto/order"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/consumer"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/handler"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/repository"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/service"
)

func main() {

	ctx := context.Background()

	// ---------------- DB ----------------
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/orders_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewPostgresOrderRepo(db)
	orderService := service.NewOrderService(orderRepo)

	handler := handler.NewOrderHandler(orderService)

	// ---------------- Kafka Config ----------------
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	// ---------------- Consumer Group (with retry) ----------------
	var consumerGroup sarama.ConsumerGroup

	for i := 0; i < 5; i++ {
		consumerGroup, err = sarama.NewConsumerGroup([]string{broker}, "order-service-group", config)
		if err == nil {
			break
		}

		log.Println("retrying kafka consumer connection...", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect to kafka:", err)
	}

	// ---------------- Kafka Consumer ----------------
	paymentConsumer := consumer.NewPaymentConsumer(orderService)

	// ---------------- Start Consumer ----------------
	go func() {
		for {
			err := consumerGroup.Consume(ctx, []string{"payments_topic"}, paymentConsumer)
			if err != nil {
				log.Println("consumer error:", err)
			}
		}
	}()

	// ---------------- gRPC Server ----------------
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, handler)

	log.Println("Order Service running on :50051")

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
