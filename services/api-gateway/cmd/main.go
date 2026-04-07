package main

import (
	"log"
	"net/http"
	"github.com/Bharat1Rajput/flowpay/services/api-gateway/internal/client"
	"github.com/Bharat1Rajput/flowpay/services/api-gateway/internal/handler"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {

	// gRPC connections
	orderConn, _ := grpc.Dial("order-service:50051", grpc.WithInsecure())
	paymentConn, _ := grpc.Dial("payment-service:50052", grpc.WithInsecure())

	orderClient := client.NewOrderClient(orderConn)
	paymentClient := client.NewPaymentClient(paymentConn)

	h := handler.NewHandler(orderClient, paymentClient)

	r := chi.NewRouter()  
	r.Post("/orders", h.CreateOrder)
	r.Post("/payments", h.ProcessPayment)
	r.Get("/orders/{id}", h.GetOrder)
	r.Get("/payments/{id}", h.GetPayment)
	r.Post("/orders/cancel/{id}", h.CancelOrder)
	
	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
