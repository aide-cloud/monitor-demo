package service

import (
	"context"
	"time"

	pb "monitor-demo/api/order"
	"monitor-demo/internal/metrics"
)

type OrderService struct {
	pb.UnimplementedOrderServer
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	startTime := time.Now()
	defer func() {
		// 记录请求处理时间
		metrics.OrderRequestDuration.WithLabelValues("create").Observe(time.Since(startTime).Seconds())
	}()

	reply, err := s.createOrder(ctx, req)
	if err != nil {
		// 记录失败请求
		metrics.OrderRequests.WithLabelValues("create", "error").Inc()
		return nil, err
	}

	// 记录成功请求
	metrics.OrderRequests.WithLabelValues("create", "success").Inc()
	return reply, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	startTime := time.Now()
	defer func() {
		// 记录请求处理时间
		metrics.OrderRequestDuration.WithLabelValues("query").Observe(time.Since(startTime).Seconds())
	}()

	reply, err := s.getOrder(ctx, req)
	if err != nil {
		// 记录失败请求
		metrics.OrderRequests.WithLabelValues("query", "error").Inc()
		return nil, err
	}

	// 记录成功请求
	metrics.OrderRequests.WithLabelValues("query", "success").Inc()
	// 更新订单状态计数
	metrics.OrderStatusGauge.WithLabelValues(getStatusString(reply.Status)).Inc()

	return reply, nil
}

// 辅助函数：将状态码转换为字符串
func getStatusString(status int64) string {
	// 根据您的业务逻辑定义状态映射
	statusMap := map[int64]string{
		0: "pending",
		1: "processing",
		2: "completed",
		3: "failed",
		// 添加其他状态...
	}
	if s, ok := statusMap[status]; ok {
		return s
	}
	return "unknown"
}

// 实际的订单创建逻辑
func (s *OrderService) createOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	// 实现订单创建逻辑
	return &pb.CreateOrderReply{}, nil
}

// 实际的订单查询逻辑
func (s *OrderService) getOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	// 实现订单查询逻辑
	return &pb.GetOrderReply{}, nil
}
