package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	pb "monitor-demo/api/order"
	"monitor-demo/internal/metrics"

	"github.com/go-kratos/kratos/v2/errors"
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
		metrics.OrderRequestDuration.WithLabelValues("create").Observe(time.Since(startTime).Seconds())
	}()

	reply, err := s.createOrder(ctx, req)
	if err != nil {
		// 增加错误类型标签
		errCode := errors.FromError(err).GetReason()
		switch errCode {
		case pb.ErrorReason_ORDER_TIMEOUT.String():
			metrics.OrderRequests.WithLabelValues("create", "error", "timeout").Inc()
		case pb.ErrorReason_ORDER_GOODS_UNAUTHORIZED.String():
			metrics.OrderRequests.WithLabelValues("create", "error", "unauthorized").Inc()
		case pb.ErrorReason_ORDER_GOODS_PRICE_ERROR.String():
			metrics.OrderRequests.WithLabelValues("create", "error", "invalid_price").Inc()
		case pb.ErrorReason_ORDER_GOODS_QUANTITY_ERROR.String():
			metrics.OrderRequests.WithLabelValues("create", "error", "invalid_quantity").Inc()
		default:
			metrics.OrderRequests.WithLabelValues("create", "error", "unknown").Inc()
		}
		return nil, err
	}

	metrics.OrderRequests.WithLabelValues("create", "success", "").Inc()
	return reply, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	startTime := time.Now()
	defer func() {
		metrics.OrderRequestDuration.WithLabelValues("query").Observe(time.Since(startTime).Seconds())
	}()

	reply, err := s.getOrder(ctx, req)
	if err != nil {
		if pb.IsOrderNotFound(err) {
			metrics.OrderRequests.WithLabelValues("query", "error", "not_found").Inc()
		} else {
			metrics.OrderRequests.WithLabelValues("query", "error", "unknown").Inc()
		}
		return nil, err
	}

	metrics.OrderRequests.WithLabelValues("query", "success", "").Inc()
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
	// 模拟业务校验
	if req.Price <= 0 {
		return nil, pb.ErrorOrderGoodsPriceError("商品价格必须大于0")
	}
	if req.Quantity <= 0 {
		return nil, pb.ErrorOrderGoodsQuantityError("商品数量必须大于0")
	}

	// 模拟随机延迟 (100ms-1s)
	delay := 100 + rand.Intn(900)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// 模拟随机错误 (10%概率)
	if rand.Float32() < 0.1 {
		return nil, pb.ErrorOrderTimeout("订单创建超时")
	}

	// 模拟商品权限校验 (5%概率无权限)
	if rand.Float32() < 0.05 {
		return nil, pb.ErrorOrderGoodsUnauthorized("无权购买该商品")
	}

	// 生成订单号 (简单模拟)
	orderNo := fmt.Sprintf("ORD%d%d", time.Now().Unix(), rand.Intn(1000))

	return &pb.CreateOrderReply{
		OrderNo:      orderNo,
		OuterTradeNo: req.OuterTradeNo,
	}, nil
}

// 实际的订单查询逻辑
func (s *OrderService) getOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	// 模拟随机延迟 (50ms-500ms)
	delay := 50 + rand.Intn(450)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// 模拟订单不存在 (15%概率)
	if rand.Float32() < 0.15 {
		return nil, pb.ErrorOrderNotFound("订单不存在")
	}

	// 模拟随机订单状态
	status := rand.Int63n(4) // 0-3 对应 pending, processing, completed, failed

	// 模拟订单数据
	return &pb.GetOrderReply{
		OrderNo:      req.OrderNo,
		OuterTradeNo: req.OuterTradeNo,
		GoodsCode:    req.GoodsCode,
		GoodsName:    "模拟商品",
		Price:        rand.Int63n(10000) + 1,
		Quantity:     rand.Int63n(10) + 1,
		Status:       status,
	}, nil
}

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}
