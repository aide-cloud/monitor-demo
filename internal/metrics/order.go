package metrics

import (
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
)

func init() {
	// 注册指标
	prometheus.MustRegister(OrderRequests)
	prometheus.MustRegister(OrderRequestDuration)
	prometheus.MustRegister(OrderStatusGauge)
}
