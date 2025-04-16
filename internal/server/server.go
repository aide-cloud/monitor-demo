package server

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"monitor-demo/api/order"
	"monitor-demo/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

// 模拟请求的辅助函数
func simulateRequests(orderService *service.OrderService) {
	// 创建一个上下文
	ctx := context.Background()

	// 随机生成订单数据
	generateOrderRequest := func() *order.CreateOrderRequest {
		return &order.CreateOrderRequest{
			OuterTradeNo: generateTradeNo(),
			GoodsCode:    generateGoodsCode(),
			GoodsName:    "测试商品",
			Price:        rand.Int63n(10000) + 1, // 1-10000
			Quantity:     rand.Int63n(10) + 1,    // 1-10
		}
	}

	// 随机生成查询请求
	generateQueryRequest := func() *order.GetOrderRequest {
		return &order.GetOrderRequest{
			OrderNo:      generateOrderNo(),
			OuterTradeNo: generateTradeNo(),
			GoodsCode:    generateGoodsCode(),
		}
	}

	// 启动定时任务
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// 模拟创建订单请求
			go func() {
				_, err := orderService.CreateOrder(ctx, generateOrderRequest())
				log.Debugw("msg", "create order", "err", err)
			}()

			// 模拟查询订单请求
			go func() {
				_, err := orderService.GetOrder(ctx, generateQueryRequest())
				log.Debugw("msg", "query order", "err", err)
			}()
		}
	}()
}

// 生成随机订单号
func generateOrderNo() string {
	return fmt.Sprintf("ORD%d%d", time.Now().UnixNano(), rand.Intn(1000))
}

// 生成随机交易号
func generateTradeNo() string {
	return fmt.Sprintf("TRD%d%d", time.Now().UnixNano(), rand.Intn(1000))
}

// 生成随机商品编码
func generateGoodsCode() string {
	return fmt.Sprintf("GOODS%d", rand.Intn(100))
}
