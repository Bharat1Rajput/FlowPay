package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/Bharat1Rajput/flowpay/proto/order"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/handler"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/repository"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/service"
)

func main() {

	// Connect DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/orders_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// service, repo, handler wiring
	repo := repository.NewPostgresOrderRepo(db)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	pb.RegisterOrderServiceServer(server, h)

	log.Println("Order Service running on :50051")
	log.Fatal(server.Serve(lis))
}
