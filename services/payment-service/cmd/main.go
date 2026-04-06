package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/Bharat1Rajput/flowpay/proto/payment"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/handler"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/kafka"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/repository"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/service"
)

func main() {

	// ---------------- DB ----------------
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/payments_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresPaymentRepo(db)

	// ---------------- Kafka Config ----------------
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_1_0_0 // safer default

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	// ---------------- Kafka Producer (with retry) ----------------
	var producer sarama.SyncProducer

	for i := 0; i < 5; i++ {
		producer, err = sarama.NewSyncProducer([]string{broker}, config)
		if err == nil {
			break
		}

		log.Println("retrying kafka connection...", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect to kafka:", err)
	}

	kafkaProducer := kafka.NewKafkaProducer(producer, "payments_topic")

	// ---------------- Service ----------------
	svc := service.NewPaymentService(repo, kafkaProducer)

	// ---------------- Handler ----------------
	h := handler.NewPaymentHandler(svc)

	// ---------------- gRPC Server ----------------
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterPaymentServiceServer(server, h)

	log.Println("Payment Service running on :50052")

	// ---------------- Start Server ----------------
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Println("server error:", err)
		}
	}()

	// ---------------- Block Forever ----------------
	select {}
}
