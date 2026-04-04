package handler

import (
	"encoding/json"
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
	}

	_ = json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.orderClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		UserId:          req.UserID,
		DeliveryAddress: req.DeliveryAddress,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
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

	_ = json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.paymentClient.ProcessPayment(r.Context(), &pb2.ProcessPaymentRequest{
		OrderId:        req.OrderID,
		IdempotencyKey: req.IdempotencyKey,
		AmountPaise:    req.Amount,
		Currency:       "INR",
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
