package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/Bharat1Rajput/flowpay/services/api-gateway/internal/client"
	"github.com/Bharat1Rajput/flowpay/services/api-gateway/internal/handler"
)

func main() {

	// gRPC connections
	orderConn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	paymentConn, _ := grpc.Dial("localhost:50052", grpc.WithInsecure())

	orderClient := client.NewOrderClient(orderConn)
	paymentClient := client.NewPaymentClient(paymentConn)

	h := handler.NewHandler(orderClient, paymentClient)

	http.HandleFunc("/orders", h.CreateOrder)
	http.HandleFunc("/payments", h.ProcessPayment)

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
