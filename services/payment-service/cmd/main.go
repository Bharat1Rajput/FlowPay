package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/Bharat1Rajput/flowpay/proto/payment"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/handler"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/kafka"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/repository"
	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/service"
)

func main() {

	// DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/payments_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresPaymentRepo(db)

	// TEMP: mock producer (no kafka yet)
	var producer kafka.EventProducer = &kafka.MockProducer{}

	svc := service.NewPaymentService(repo, producer)

	h := handler.NewPaymentHandler(svc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	pb.RegisterPaymentServiceServer(server, h)

	log.Println("Payment Service running on :50052")
	log.Fatal(server.Serve(lis))
}
