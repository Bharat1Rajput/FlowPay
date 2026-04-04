package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/Bharat1Rajput/flowpay/proto/order"
	pb2 "github.com/Bharat1Rajput/flowpay/proto/payment"
	"github.com/Bharat1Rajput/flowpay/services/api-gateway/internal/client"
)

type Handler struct {
	orderClient   *client.OrderClient
	paymentClient *client.PaymentClient
}

func NewHandler(o *client.OrderClient, p *client.PaymentClient) *Handler {
	return &Handler{
		orderClient:   o,
		paymentClient: p,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var req struct {
		UserID          string `json:"user_id"`
		DeliveryAddress string `json:"delivery_address"`
		Items           []struct {
			Name      string `json:"name"`
			Quantity  int32  `json:"quantity"`
			UnitPrice int64  `json:"unit_price"`
		} `json:"items"`
		Notes string `json:"notes"`
	}

	// ✅ Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 🔍 Debug (VERY IMPORTANT during development)
	fmt.Println("HTTP Items:", len(req.Items))

	//
	if len(req.Items) == 0 {
		http.Error(w, "items cannot be empty", http.StatusBadRequest)
		return
	}

	// Convert HTTP → gRPC (CRITICAL STEP)
	var items []*pb.OrderItem

	for _, i := range req.Items {
		items = append(items, &pb.OrderItem{
			Name:           i.Name,
			Quantity:       i.Quantity,
			UnitPricePaise: i.UnitPrice,
		})
	}

	// gRPC call with FULL mapping
	resp, err := h.orderClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		UserId:          req.UserID,
		DeliveryAddress: req.DeliveryAddress,
		Notes:           req.Notes,
		Items:           items, // 🔥 FIXED
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {

	var req struct {
		OrderID        string `json:"order_id"`
		IdempotencyKey string `json:"idempotency_key"`
		Amount         int64  `json:"amount"`
	}

	// ✅ Decode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 🔍 Debug
	fmt.Println("Processing payment for order:", req.OrderID)

	// ❌ Basic validation
	if req.OrderID == "" || req.IdempotencyKey == "" || req.Amount <= 0 {
		http.Error(w, "invalid payment request", http.StatusBadRequest)
		return
	}

	// ✅ gRPC call
	resp, err := h.paymentClient.ProcessPayment(r.Context(), &pb2.ProcessPaymentRequest{
		OrderId:        req.OrderID,
		IdempotencyKey: req.IdempotencyKey,
		AmountPaise:    req.Amount,
		Currency:       "INR",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
