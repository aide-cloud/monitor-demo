# Monitor Demo Service

[中文文档](README.zh-CN.md)

A microservice demo with Prometheus monitoring integration based on Kratos framework.

## Prerequisites

- Go 1.19+
- Docker & Docker Compose
- Kratos CLI

## Installation

### Install Kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## Service Startup

### Local Development
```bash
# Build the service
make build

# Run the service
./bin/server -conf ./configs
```

### Docker Deployment
```bash
# Build docker image
docker build -t monitor-demo .

# Run container
docker run --rm -p 8000:8000 -p 9000:9000 -v $(pwd)/configs:/data/conf monitor-demo
```

## Monitoring Stack Setup

### Start Prometheus Stack
```bash
cd prometheus
docker-compose up -d
```

This will start:
- Prometheus (http://localhost:9090)
- Alertmanager (http://localhost:9093)
- Grafana (http://localhost:3000)
  - Default credentials: admin/123456

## Available Metrics

The service exposes metrics at `/metrics` endpoint. Current implemented metrics include:

### Order Metrics

1. `monitor_demo_order_requests_total`
   - Type: Counter
   - Description: Total number of order requests
   - Labels: 
     - `method`: API method name
     - `status`: Request status
     - `error_type`: Type of error if any

2. `monitor_demo_order_request_duration_seconds`
   - Type: Histogram
   - Description: Order request duration in seconds
   - Labels:
     - `method`: API method name
   - Buckets: 0.1, 0.3, 0.5, 0.7, 0.9, 1.0, 1.5, 2.0

3. `monitor_demo_order_status_count`
   - Type: Gauge
   - Description: Current count of orders by status
   - Labels:
     - `status`: Order status

## Metric Types Guide

### 1. Counter
- **Use Case**: For values that only increase (e.g., request count, error count)
- **Definition**:
```go
prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: "monitor_demo",
        Subsystem: "requests",
        Name:      "total",
        Help:      "Total number of requests processed",
    },
    []string{"method", "status"},
)
```
- **Usage**:
```go
// Increment by 1
counter.WithLabelValues("GET", "success").Inc()
// Increment by specific value
counter.WithLabelValues("POST", "error").Add(2)
```

### 2. Gauge
- **Use Case**: For values that can go up and down (e.g., current memory usage, active requests)
- **Definition**:
```go
prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Namespace: "monitor_demo",
        Subsystem: "resources",
        Name:      "active_requests",
        Help:      "Number of active requests",
    },
    []string{"endpoint"},
)
```
- **Usage**:
```go
// Set specific value
gauge.WithLabelValues("/api/v1").Set(42)
// Increment/Decrement
gauge.WithLabelValues("/api/v1").Inc()
gauge.WithLabelValues("/api/v1").Dec()
```

### 3. Histogram
- **Use Case**: For measuring distributions of values (e.g., request duration, response size)
- **Definition**:
```go
prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Namespace: "monitor_demo",
        Subsystem: "requests",
        Name:      "duration_seconds",
        Help:      "Request duration in seconds",
        Buckets:   []float64{0.1, 0.5, 1, 2, 5}, // Define buckets
    },
    []string{"method"},
)
```
- **Usage**:
```go
// Observe a value
histogram.WithLabelValues("GET").Observe(0.42)
```

### 4. Summary
- **Use Case**: For calculating quantiles over time (e.g., request latency percentiles)
- **Definition**:
```go
prometheus.NewSummaryVec(
    prometheus.SummaryOpts{
        Namespace:  "monitor_demo",
        Subsystem:  "requests",
        Name:       "latency_seconds",
        Help:       "Request latency in seconds",
        Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // Quantiles
    },
    []string{"method"},
)
```
- **Usage**:
```go
// Observe a value
summary.WithLabelValues("GET").Observe(0.42)
```

### Key Differences
- **Counter**: Only increases, resets on restart
- **Gauge**: Can increase and decrease, current state
- **Histogram**: Measures distribution with predefined buckets
- **Summary**: Calculates configurable quantiles over sliding time window

## Adding New Metrics

### 1. Define New Metrics

Create new metrics in `internal/metrics/` directory. Example:

```go
var NewMetric = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: "monitor_demo",
        Subsystem: "your_subsystem",
        Name:      "metric_name",
        Help:      "Metric description",
    },
    []string{"label1", "label2"},
)
```

### 2. Register Metrics

Register in the `init()` function:

```go
func init() {
    prometheus.MustRegister(NewMetric)
}
```

### 3. Use Metrics in Code

```go
// Increment counter
NewMetric.WithLabelValues("label1_value", "label2_value").Inc()

// Add specific value
NewMetric.WithLabelValues("label1_value", "label2_value").Add(123)
```

### 4. Best Practices

- Use meaningful metric names following the pattern: `namespace_subsystem_name`
- Add descriptive help text
- Choose appropriate metric types:
  - Counter: For values that only increase
  - Gauge: For values that can go up and down
  - Histogram: For measuring distributions of values
  - Summary: For calculating quantiles

## Development

### Generate Code
```bash
# Download dependencies
make init

# Generate API files
make api

# Generate all files
make all
```

### Wire Generation
```bash
# Install wire
go get github.com/google/wire/cmd/wire

# Generate wire
cd cmd/monitor-demo
wire
```


## Metric Query Examples

### 1. Counter Queries
Using `monitor_demo_order_requests_total` as example:

#### Basic Queries
```promql
# Total requests in the last hour
sum(increase(monitor_demo_order_requests_total[1h]))

# Error rate over 5 minutes
sum(rate(monitor_demo_order_requests_total{status="error"}[5m])) 
  / 
sum(rate(monitor_demo_order_requests_total[5m])) * 100

# Requests per second by method
rate(monitor_demo_order_requests_total[5m])
```

#### Visualization
- **Graph Type**: Line graph
- **Use Case**: Request rate trends
```promql
# Request rate by method (dashboard query)
sum by(method) (rate(monitor_demo_order_requests_total[5m]))
```
- **Graph Type**: Pie chart
- **Use Case**: Error distribution
```promql
# Error distribution by type
sum by(error_type) (increase(monitor_demo_order_requests_total{status="error"}[24h]))
```

### 2. Gauge Queries
Using `monitor_demo_order_status_count` as example:

#### Basic Queries
```promql
# Current value
monitor_demo_order_status_count

# Max value over last hour
max_over_time(monitor_demo_order_status_count[1h])

# Average by status
avg by(status) (monitor_demo_order_status_count)
```

#### Visualization
- **Graph Type**: Gauge/Single stat
- **Use Case**: Current active orders
```promql
# Total current orders
sum(monitor_demo_order_status_count)
```
- **Graph Type**: Bar graph
- **Use Case**: Orders by status
```promql
# Distribution by status
monitor_demo_order_status_count
```

### 3. Histogram Queries
Using `monitor_demo_order_request_duration_seconds` as example:

#### Basic Queries
```promql
# 95th percentile response time
histogram_quantile(0.95, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))

# Average response time
rate(monitor_demo_order_request_duration_seconds_sum[5m])
  /
rate(monitor_demo_order_request_duration_seconds_count[5m])

# Request count by duration bucket
sum by(le) (increase(monitor_demo_order_request_duration_seconds_bucket[1h]))
```

#### Visualization
- **Graph Type**: Heatmap
- **Use Case**: Request duration distribution
```promql
# Duration distribution over time
sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le)
```
- **Graph Type**: Line graph
- **Use Case**: Multiple percentiles over time
```promql
# Multiple percentiles
histogram_quantile(0.50, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))
histogram_quantile(0.90, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))
histogram_quantile(0.95, sum(rate(monitor_demo_order_request_duration_seconds_bucket[5m])) by (le))
```

### 4. Summary Queries
Using a hypothetical `monitor_demo_request_latency_seconds` summary:

#### Basic Queries
```promql
# Direct percentile values
monitor_demo_request_latency_seconds{quantile="0.95"}

# Average latency
rate(monitor_demo_request_latency_seconds_sum[5m])
  /
rate(monitor_demo_request_latency_seconds_count[5m])

# Request count
rate(monitor_demo_request_latency_seconds_count[5m])
```

#### Visualization
- **Graph Type**: Line graph
- **Use Case**: Latency trends by percentile
```promql
# Multiple percentiles over time
monitor_demo_request_latency_seconds{quantile=~"0.5|0.9|0.99"}
```
- **Graph Type**: Stat panel
- **Use Case**: Current p99 latency
```promql
# Current p99 latency
monitor_demo_request_latency_seconds{quantile="0.99"}
```

## PromQL Function Reference

### 1. Aggregation Operators
Common aggregation operators that can be used with `by` and `without`:

| Operator | Description | Example |
|----------|-------------|---------|
| `sum` | Sum of values | `sum(http_requests_total)` |
| `avg` | Average value | `avg(node_memory_usage_bytes)` |
| `min` | Minimum value | `min(node_cpu_utilization)` |
| `max` | Maximum value | `max(node_temperature_celsius)` |
| `count` | Count of elements | `count(up)` |
| `stddev` | Standard deviation | `stddev(http_request_duration_seconds)` |
| `topk` | Top K elements | `topk(3, http_errors_total)` |
| `bottomk` | Bottom K elements | `bottomk(3, node_memory_free_bytes)` |

### 2. Time-based Functions

#### Rate Functions
| Function | Description | Example |
|----------|-------------|---------|
| `rate` | Per-second rate of increase (for counters) | `rate(http_requests_total[5m])` |
| `irate` | Instant rate based on last two samples | `irate(http_requests_total[5m])` |
| `increase` | Total increase in time window | `increase(http_requests_total[1h])` |

#### Time Window Functions
| Function | Description | Example |
|----------|-------------|---------|
| `avg_over_time` | Average over time window | `avg_over_time(node_cpu[10m])` |
| `min_over_time` | Minimum over time window | `min_over_time(node_temp[1h])` |
| `max_over_time` | Maximum over time window | `max_over_time(node_temp[1h])` |
| `sum_over_time` | Sum over time window | `sum_over_time(error_count[24h])` |
| `count_over_time` | Count of samples in window | `count_over_time(up[5m])` |

### 3. Histogram Functions
| Function | Description | Example |
|----------|-------------|---------|
| `histogram_quantile` | Calculate quantile from histogram | `histogram_quantile(0.95, http_request_duration_bucket)` |
| `rate` with histograms | Calculate rate of histogram samples | `rate(http_request_duration_bucket[5m])` |

### 4. Mathematical Functions
| Function | Description | Example |
|----------|-------------|---------|
| `abs` | Absolute value | `abs(delta(temperature[1h]))` |
| `ceil` | Round up | `ceil(node_cpu_usage)` |
| `floor` | Round down | `floor(node_cpu_usage)` |
| `round` | Round to nearest integer | `round(node_cpu_usage)` |
| `sqrt` | Square root | `sqrt(http_requests_total)` |
| `exp` | Exponential | `exp(node_pressure)` |
| `ln` | Natural logarithm | `ln(node_volume)` |
| `log2` | Binary logarithm | `log2(process_heap_bytes)` |
| `log10` | Decimal logarithm | `log10(process_heap_bytes)` |

### 5. Common Usage Patterns

#### Error Rate Calculation
```promql
# Error percentage over 5m
sum(rate(http_requests_total{status="error"}[5m]))
  /
sum(rate(http_requests_total[5m])) * 100
```

#### Apdex Score Calculation
```promql
# Apdex score (satisfied + tolerating/2) / total
(
  sum(rate(http_request_duration_bucket{le="0.3"}[5m])) +
  sum(rate(http_request_duration_bucket{le="1.2"}[5m])) / 2
) / sum(rate(http_request_duration_count[5m]))
```

#### Resource Utilization
```promql
# CPU utilization percentage
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
```

#### Service Availability (SLI)
```promql
# Success rate as SLI
sum(rate(http_requests_total{code=~"2.."}[5m]))
  /
sum(rate(http_requests_total[5m]))
```


