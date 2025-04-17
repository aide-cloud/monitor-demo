# Monitor Demo 服务

[English](README.md)

基于 Kratos 框架的微服务示例，集成了 Prometheus 监控。

## 前置要求

- Go 1.19+
- Docker & Docker Compose
- Kratos CLI

## 安装

### 安装 Kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## 服务启动

### 本地开发
```bash
# 构建服务
make build

# 运行服务
./bin/server -conf ./configs
```

### Docker 部署
```bash
# 构建 docker 镜像
docker build -t monitor-demo .

# 运行容器
docker run --rm -p 8000:8000 -p 9000:9000 -v $(pwd)/configs:/data/conf monitor-demo
```

## 监控栈设置

### 启动 Prometheus 栈
```bash
cd prometheus
docker-compose up -d
```

这将启动以下服务：
- Prometheus (http://localhost:9090)
- Alertmanager (http://localhost:9093)
- Grafana (http://localhost:3000)
  - 默认登录凭据: admin/123456

## 可用指标

服务在 `/metrics` 端点暴露指标。当前实现的指标包括：

### 订单指标

1. `monitor_demo_order_requests_total`
   - 类型：计数器
   - 描述：订单请求总数
   - 标签：
     - `method`：API 方法名
     - `status`：请求状态
     - `error_type`：错误类型（如果有）

2. `monitor_demo_order_request_duration_seconds`
   - 类型：直方图
   - 描述：订单请求处理时间（秒）
   - 标签：
     - `method`：API 方法名
   - 分桶：0.1, 0.3, 0.5, 0.7, 0.9, 1.0, 1.5, 2.0

3. `monitor_demo_order_status_count`
   - 类型：仪表盘
   - 描述：按状态统计的当前订单数量
   - 标签：
     - `status`：订单状态

## 指标类型指南

### 1. Counter（计数器）
- **使用场景**：用于只增不减的值（如：请求总数、错误总数）
- **定义方式**：
```go
prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: "monitor_demo",
        Subsystem: "requests",
        Name:      "total",
        Help:      "处理的请求总数",
    },
    []string{"method", "status"},
)
```
- **使用方法**：
```go
// 增加1
counter.WithLabelValues("GET", "success").Inc()
// 增加指定值
counter.WithLabelValues("POST", "error").Add(2)
```

### 2. Gauge（仪表盘）
- **使用场景**：用于可增可减的值（如：当前内存使用量、活跃请求数）
- **定义方式**：
```go
prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Namespace: "monitor_demo",
        Subsystem: "resources",
        Name:      "active_requests",
        Help:      "当前活跃请求数",
    },
    []string{"endpoint"},
)
```
- **使用方法**：
```go
// 设置特定值
gauge.WithLabelValues("/api/v1").Set(42)
// 增加/减少
gauge.WithLabelValues("/api/v1").Inc()
gauge.WithLabelValues("/api/v1").Dec()
```

### 3. Histogram（直方图）
- **使用场景**：用于测量值的分布情况（如：请求持续时间、响应大小）
- **定义方式**：
```go
prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Namespace: "monitor_demo",
        Subsystem: "requests",
        Name:      "duration_seconds",
        Help:      "请求持续时间（秒）",
        Buckets:   []float64{0.1, 0.5, 1, 2, 5}, // 定义分桶
    },
    []string{"method"},
)
```
- **使用方法**：
```go
// 观察一个值
histogram.WithLabelValues("GET").Observe(0.42)
```

### 4. Summary（摘要）
- **使用场景**：用于计算时间窗口内的分位数（如：请求延迟百分位数）
- **定义方式**：
```go
prometheus.NewSummaryVec(
    prometheus.SummaryOpts{
        Namespace:  "monitor_demo",
        Subsystem:  "requests",
        Name:       "latency_seconds",
        Help:       "请求延迟（秒）",
        Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // 分位数
    },
    []string{"method"},
)
```
- **使用方法**：
```go
// 观察一个值
summary.WithLabelValues("GET").Observe(0.42)
```

### 主要区别
- **Counter**：只能增加，重启时重置
- **Gauge**：可增可减，反映当前状态
- **Histogram**：使用预定义的桶测量分布情况
- **Summary**：在滑动时间窗口内计算可配置的分位数

## 添加新指标

### 1. 定义新指标

在 `internal/metrics/` 目录下创建新指标。示例：

```go
var NewMetric = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: "monitor_demo",
        Subsystem: "your_subsystem",
        Name:      "metric_name",
        Help:      "指标描述",
    },
    []string{"label1", "label2"},
)
```

### 2. 注册指标

在 `init()` 函数中注册：

```go
func init() {
    prometheus.MustRegister(NewMetric)
}
```

### 3. 在代码中使用指标

```go
// 增加计数器
NewMetric.WithLabelValues("label1_value", "label2_value").Inc()

// 添加特定值
NewMetric.WithLabelValues("label1_value", "label2_value").Add(123)
```

### 4. 最佳实践

- 使用有意义的指标名称，遵循模式：`namespace_subsystem_name`
- 添加描述性的帮助文本
- 选择合适的指标类型：
  - Counter（计数器）：用于只增不减的值
  - Gauge（仪表盘）：用于可增可减的值
  - Histogram（直方图）：用于测量值的分布
  - Summary（摘要）：用于计算分位数

## 开发

### 生成代码
```bash
# 下载依赖
make init

# 生成 API 文件
make api

# 生成所有文件
make all
```

### Wire 生成
```bash
# 安装 wire
go get github.com/google/wire/cmd/wire

# 生成 wire
cd cmd/monitor-demo
wire
```

## 指标查询示例

### 1. Counter（计数器）查询
以 `monitor_demo_order_requests_total` 为例：

#### 基础查询
```promql
# 最近一小时的总请求数
sum(increase(monitor_demo_order_requests_total[1h]))

# 5分钟内的错误率
sum(rate(monitor_demo_order_requests_total{status="error"}[5m])) 
  / 
sum(rate(monitor_demo_order_requests_total[5m])) * 100

# 按方法统计每秒请求数
rate(monitor_demo_order_requests_total[5m])
```

#### 可视化
- **图表类型**：折线图
- **使用场景**：请求率趋势
```promql
# 按方法统计的请求率（仪表板查询）
sum by(method) (rate(monitor_demo_order_requests_total[5m]))
```
- **图表类型**：饼图
- **使用场景**：错误分布
```promql
# 按类型统计24小时内的错误分布
sum by(error_type) (increase(monitor_demo_order_requests_total{status="error"}[24h]))
```

### 2. Gauge（仪表盘）查询
以 `monitor_demo_order_status_count` 为例：

#### 基础查询
```promql
# 当前值
monitor_demo_order_status_count

# 最近一小时的最大值
max_over_time(monitor_demo_order_status_count[1h])

# 按状态统计平均值
avg by(status) (monitor_demo_order_status_count)
```

#### 可视化
- **图表类型**：仪表盘/单值统计
- **使用场景**：当前活跃订单数
```promql
# 当前订单总数
sum(monitor_demo_order_status_count)
```
- **图表类型**：柱状图
- **使用场景**：各状态订单数量
```promql
# 按状态分布
monitor_demo_order_status_count
```

### 3. Histogram（直方图）查询
以 `monitor_demo_order_request_duration_seconds` 为例：

#### 基础查询
```promql
# 95分位响应时间
histogram_quantile(0.95, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))

# 平均响应时间
rate(monitor_demo_order_request_duration_seconds_sum[5m])
  /
rate(monitor_demo_order_request_duration_seconds_count[5m])

# 按持续时间区间统计请求数
sum by(le) (increase(monitor_demo_order_request_duration_seconds_bucket[1h]))
```

#### 可视化
- **图表类型**：热力图
- **使用场景**：请求持续时间分布
```promql
# 随时间变化的持续时间分布
sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le)
```
- **图表类型**：折线图
- **使用场景**：多个百分位随时间变化
```promql
# 多个百分位
histogram_quantile(0.50, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))
histogram_quantile(0.90, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))
histogram_quantile(0.95, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))
```

### 4. Summary（摘要）查询
以假设的 `monitor_demo_request_latency_seconds` 摘要为例：

#### 基础查询
```promql
# 直接获取百分位数值
monitor_demo_request_latency_seconds{quantile="0.95"}

# 平均延迟
rate(monitor_demo_request_latency_seconds_sum[5m])
  /
rate(monitor_demo_request_latency_seconds_count[5m])

# 请求计数
rate(monitor_demo_request_latency_seconds_count[5m])
```

#### 可视化
- **图表类型**：折线图
- **使用场景**：各百分位延迟趋势
```promql
# 多个百分位随时间变化
monitor_demo_request_latency_seconds{quantile=~"0.5|0.9|0.99"}
```
- **图表类型**：统计面板
- **使用场景**：当前P99延迟
```promql
# 当前P99延迟
monitor_demo_request_latency_seconds{quantile="0.99"}
```

## PromQL 函数参考

### 1. 聚合运算符
可与 `by` 和 `without` 一起使用的常用聚合运算符：

| 运算符 | 描述 | 示例 |
|--------|------|------|
| `sum` | 值的总和 | `sum(http_requests_total)` |
| `avg` | 平均值 | `avg(node_memory_usage_bytes)` |
| `min` | 最小值 | `min(node_cpu_utilization)` |
| `max` | 最大值 | `max(node_temperature_celsius)` |
| `count` | 元素计数 | `count(up)` |
| `stddev` | 标准差 | `stddev(http_request_duration_seconds)` |
| `topk` | 前K个最大元素 | `topk(3, http_errors_total)` |
| `bottomk` | 前K个最小元素 | `bottomk(3, node_memory_free_bytes)` |

### 2. 基于时间的函数

#### 速率函数
| 函数 | 描述 | 示例 |
|------|------|------|
| `rate` | 每秒增长率（用于计数器） | `rate(http_requests_total[5m])` |
| `irate` | 基于最后两个样本的瞬时速率 | `irate(http_requests_total[5m])` |
| `increase` | 时间窗口内的总增长量 | `increase(http_requests_total[1h])` |

#### 时间窗口函数
| 函数 | 描述 | 示例 |
|------|------|------|
| `avg_over_time` | 时间窗口内的平均值 | `avg_over_time(node_cpu[10m])` |
| `min_over_time` | 时间窗口内的最小值 | `min_over_time(node_temp[1h])` |
| `max_over_time` | 时间窗口内的最大值 | `max_over_time(node_temp[1h])` |
| `sum_over_time` | 时间窗口内的总和 | `sum_over_time(error_count[24h])` |
| `count_over_time` | 时间窗口内的样本数 | `count_over_time(up[5m])` |

### 3. 直方图函数
| 函数 | 描述 | 示例 |
|------|------|------|
| `histogram_quantile` | 从直方图计算分位数 | `histogram_quantile(0.95, http_request_duration_bucket)` |
| `rate` 与直方图 | 计算直方图样本的速率 | `rate(http_request_duration_bucket[5m])` |

### 4. 数学函数
| 函数 | 描述 | 示例 |
|------|------|------|
| `abs` | 绝对值 | `abs(delta(temperature[1h]))` |
| `ceil` | 向上取整 | `ceil(node_cpu_usage)` |
| `floor` | 向下取整 | `floor(node_cpu_usage)` |
| `round` | 四舍五入 | `round(node_cpu_usage)` |
| `sqrt` | 平方根 | `sqrt(http_requests_total)` |
| `exp` | 指数函数 | `exp(node_pressure)` |
| `ln` | 自然对数 | `ln(node_volume)` |
| `log2` | 二进制对数 | `log2(process_heap_bytes)` |
| `log10` | 十进制对数 | `log10(process_heap_bytes)` |

### 5. 常用查询模式

#### 错误率计算
```promql
# 5分钟内的错误百分比
sum(rate(http_requests_total{status="error"}[5m]))
  /
sum(rate(http_requests_total[5m])) * 100
```

#### Apdex 分数计算
```promql
# Apdex 分数 (满意 + 可容忍/2) / 总数
(
  sum(rate(http_request_duration_bucket{le="0.3"}[5m])) +
  sum(rate(http_request_duration_bucket{le="1.2"}[5m])) / 2
) / sum(rate(http_request_duration_count[5m]))
```

#### 资源利用率
```promql
# CPU 使用率百分比
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
```

#### 服务可用性 (SLI)
```promql
# 作为 SLI 的成功率
sum(rate(http_requests_total{code=~"2.."}[5m]))
  /
sum(rate(http_requests_total[5m]))
```
