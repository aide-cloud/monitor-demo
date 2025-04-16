package service

import (
	"context"

	pb "monitor-demo/api/order"
)

type OrderService struct {
	pb.UnimplementedOrderServer
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	return &pb.CreateOrderReply{}, nil
}
func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	return &pb.GetOrderReply{}, nil
}
