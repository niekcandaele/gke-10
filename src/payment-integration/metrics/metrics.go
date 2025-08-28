package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// Request metrics
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_requests_total",
			Help: "Total number of payment requests",
		},
		[]string{"status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_request_duration_seconds",
			Help:    "Duration of payment requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	// Transaction metrics
	transactionAmount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "payment_transaction_amount_cents_total",
			Help: "Total transaction amount in cents",
		},
	)

	// Error metrics
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_errors_total",
			Help: "Total number of payment errors",
		},
		[]string{"type"},
	)

	// Rate limiting metrics
	rateLimitRejections = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "payment_rate_limit_rejections_total",
			Help: "Total number of rate limit rejections",
		},
	)
)

// Metrics provides a simplified interface for metrics recording
type Metrics struct{}

// GetInstance returns the metrics instance
func GetInstance() *Metrics {
	return &Metrics{}
}

// RecordRequest records a payment request
func (m *Metrics) RecordRequest(success bool, latency time.Duration, amount int64, cardLast4 string) {
	status := "success"
	if !success {
		status = "failure"
	}

	requestsTotal.WithLabelValues(status).Inc()
	requestDuration.WithLabelValues(status).Observe(latency.Seconds())

	if success && amount > 0 {
		transactionAmount.Add(float64(amount))
	}
}

// RecordError records an error
func (m *Metrics) RecordError(errorType string) {
	errorsTotal.WithLabelValues(errorType).Inc()
}

// RecordRejection records a rate-limited rejection
func (m *Metrics) RecordRejection() {
	rateLimitRejections.Inc()
}

// GetStats returns current metrics as a map (for backward compatibility)
func (m *Metrics) GetStats() map[string]interface{} {
	// This is now handled by Prometheus metrics
	// Keeping for backward compatibility
	return map[string]interface{}{
		"metrics": "available at /metrics endpoint",
	}
}

// PrometheusHandler returns the Prometheus metrics handler
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
