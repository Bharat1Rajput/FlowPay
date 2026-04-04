package client

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/Bharat1Rajput/flowpay/proto/payment"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *PaymentClient {
	return &PaymentClient{
		client: pb.NewPaymentServiceClient(conn),
	}
}

func (c *PaymentClient) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	return c.client.ProcessPayment(ctx, req)
}
