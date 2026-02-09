package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requests *prometheus.CounterVec
	latency  *prometheus.HistogramVec
}

func NewMetrics(registry *prometheus.Registry) *Metrics {
	m := &Metrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "validacion",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total HTTP requests",
		}, []string{"method", "path", "status"}),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "validacion",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),
	}

	registry.MustRegister(m.requests, m.latency)
	return m
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)
		status := strconv.Itoa(ww.status)
		m.requests.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		m.latency.WithLabelValues(r.Method, r.URL.Path, status).Observe(time.Since(start).Seconds())
	})
}
