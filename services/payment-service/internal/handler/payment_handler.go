package handler

import (
	"context"

	pb "github.com/Bharat1Rajput/flowpay/proto/payment"
	"github.com/google/uuid"

	"github.com/Bharat1Rajput/flowpay/services/payment-service/internal/service"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) ProcessPayment(
	ctx context.Context,
	req *pb.ProcessPaymentRequest,
) (*pb.ProcessPaymentResponse, error) {

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, err
	}

	payment, err := h.svc.ProcessPayment(ctx, service.ProcessPaymentInput{
		OrderID:        orderID,
		IdempotencyKey: req.IdempotencyKey,
		Amount:         req.AmountPaise,
		Currency:       req.Currency,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ProcessPaymentResponse{
		PaymentId:  payment.ID.String(),
		Status:     string(payment.Status),
		GatewayRef: payment.GatewayRef,
	}, nil
}

func (h *PaymentHandler) GetPayment(
	ctx context.Context,
	req *pb.GetPaymentRequest,
) (*pb.GetPaymentResponse, error) {

	id, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, err
	}

	payment, err := h.svc.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetPaymentResponse{
		PaymentId:     payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		Status:        string(payment.Status),
		AmountPaise:   payment.Amount,
		Currency:      payment.Currency,
		GatewayRef:    payment.GatewayRef,
		FailureReason: payment.FailureReason,
		CreatedAt:     payment.CreatedAt.String(),
	}, nil
}
