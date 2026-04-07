package handler

import (
	"context"
	"log"

	pb "github.com/Bharat1Rajput/flowpay/proto/order"
	"github.com/google/uuid"

	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/model"
	"github.com/Bharat1Rajput/flowpay/services/order-service/internal/service"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	var items []model.OrderItem
	for _, i := range req.Items {
		item := model.NewOrderItem(
			uuid.Nil, // will set later
			i.Name,
			int(i.Quantity),
			i.UnitPricePaise,
		)
		items = append(items, item)
	}

	order, err := h.svc.CreateOrder(ctx, service.CreateOrderInput{
		UserID:       userID,
		DeliveryAddr: req.DeliveryAddress,
		Notes:        req.Notes,
		Items:        items,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		OrderId:    order.ID.String(),
		Status:     string(order.Status),
		TotalPaise: order.TotalAmount,
	}, nil
}

func (h *OrderHandler) GetOrder(
	ctx context.Context,
	req *pb.GetOrderRequest,
) (*pb.GetOrderResponse, error) {

	id, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, err
	}

	order, err := h.svc.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	var items []*pb.OrderItem
	for _, i := range order.Items {
		items = append(items, &pb.OrderItem{
			Name:           i.ItemName,
			Quantity:       int32(i.Quantity),
			UnitPricePaise: i.UnitPrice,
			TotalPaise:     i.TotalPrice,
		})
	}

	return &pb.GetOrderResponse{
		OrderId:         order.ID.String(),
		UserId:          order.UserID.String(),
		Status:          string(order.Status),
		TotalPaise:      order.TotalAmount,
		Currency:        order.Currency,
		DeliveryAddress: order.DeliveryAddr,
		Items:           items,
		CreatedAt:       order.CreatedAt.String(),
	}, nil
}

func (h *OrderHandler) CancelOrder(
	ctx context.Context,
	req *pb.CancelOrderRequest,
) (*pb.CancelOrderResponse, error) {

	orderID, err := uuid.Parse(req.OrderId)
	log.Println("Cancelling order with ID:", orderID)
	if err != nil {
		return nil, err
	}
  
	userID, err := uuid.Parse(req.RequestingUserId)
	if err != nil {
		return nil, err
	}

	err = h.svc.CancelOrder(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}

	return &pb.CancelOrderResponse{
		Success: true,
		Message: "order cancelled",
	}, nil
}
