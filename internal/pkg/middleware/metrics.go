package middleware

import (
	"go.uber.org/fx"
	"net/http"
	"time"

	"github.com/marrgancovka/pvzService/internal/pkg/metrics"
)

type In struct {
	fx.In

	Metrics metrics.Metrics
}

type MetricsMiddleware struct {
	metrics metrics.Metrics
}

func NewMetricsMiddleware(in In) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: in.Metrics,
	}
}

func (m *MetricsMiddleware) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		m.metrics.RequestsTotal(r.Method, r.URL.Path)
		m.metrics.ResponseTime(r.Method, r.URL.Path, duration)
	})
}
