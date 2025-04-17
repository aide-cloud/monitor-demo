package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// 订单请求计数器
	OrderRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "monitor_demo",
			Subsystem: "order",
			Name:      "requests_total",
			Help:      "Total number of order requests",
		},
		[]string{"method", "status", "error_type"},
	)

	// 订单请求处理时间
	OrderRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "monitor_demo",
			Subsystem: "order",
			Name:      "request_duration_seconds",
			Help:      "Order request duration in seconds",
			Buckets:   []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1.0, 1.5, 2.0},
		},
		[]string{"method"},
	)

	// 订单状态分布
	OrderStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "monitor_demo",
			Subsystem: "order",
			Name:      "status_count",
			Help:      "Current count of orders by status",
		},
		[]string{"status"},
	)

	// 新增: 订单处理时间分位数统计
	OrderProcessingTime = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "monitor_demo",
			Subsystem: "order",
			Name:      "processing_seconds",
			Help:      "Order processing time in seconds",
			Objectives: map[float64]float64{
				0.5:  0.05,  // 第50个百分位，允许5%的误差
				0.9:  0.01,  // 第90个百分位，允许1%的误差
				0.95: 0.005, // 第95个百分位，允许0.5%的误差
				0.99: 0.001, // 第99个百分位，允许0.1%的误差
			},
			MaxAge:     10 * time.Minute, // 样本最大存活时间
			AgeBuckets: 5,                // 用于计算衰减的时间桶数量
		},
		[]string{"method", "status"},
	)

	// 新增: 订单金额分布统计
	OrderAmountSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "monitor_demo",
			Subsystem: "order",
			Name:      "amount_distribution",
			Help:      "Distribution of order amounts",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"type"},
	)
)

func init() {
	// 注册指标
	prometheus.MustRegister(OrderRequests)
	prometheus.MustRegister(OrderRequestDuration)
	prometheus.MustRegister(OrderStatusGauge)
	prometheus.MustRegister(OrderProcessingTime)
	prometheus.MustRegister(OrderAmountSummary)
}
