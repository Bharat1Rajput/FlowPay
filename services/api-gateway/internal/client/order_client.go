package client

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/Bharat1Rajput/flowpay/proto/order"
)

type OrderClient struct {
	client pb.OrderServiceClient
}

func NewOrderClient(conn *grpc.ClientConn) *OrderClient {
	return &OrderClient{
		client: pb.NewOrderServiceClient(conn),
	}
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return c.client.CreateOrder(ctx, req)
}
