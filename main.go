package main

import (
	"fmt"

	orderpb "github.com/Bharat1Rajput/flowpay/proto/order"
)

func main() {
	req := orderpb.CreateOrderRequest{}
	fmt.Println("Proto working:", req)
}
