package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	RequestsTotal(string, string)
	ResponseTime(string, string, time.Duration)
	CreatedPvzTotal(string)
	CreatedReceptionsTotal(string)
	AddedProductTotal(string)
}

type Metric struct {
	requestsTotal         *prometheus.CounterVec
	responseTime          *prometheus.HistogramVec
	createdPvzTotal       *prometheus.CounterVec
	createdReceptionTotal *prometheus.CounterVec
	addedProductTotal     *prometheus.CounterVec
}

func New() *Metric {
	requestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests",
	}, []string{"method", "path"})
	prometheus.MustRegister(requestsTotal)

	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "response_time_seconds",
		Help: "Duration of requests",
	}, []string{"method", "path"})
	prometheus.MustRegister(responseTime)

	createdPvzTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "pvz_created_total",
		Help: "Total number of created pvz",
	}, []string{"city"})
	prometheus.MustRegister(createdPvzTotal)

	createdReceptionTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "order_acceptances_created_total",
		Help: "Total number of order acceptances created",
	}, []string{"pvzId"})
	prometheus.MustRegister(createdReceptionTotal)

	addedProductTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "products_added_total",
		Help: "Total number of products added",
	}, []string{"type"})
	prometheus.MustRegister(addedProductTotal)

	return &Metric{
		requestsTotal:         requestsTotal,
		responseTime:          responseTime,
		createdPvzTotal:       createdPvzTotal,
		createdReceptionTotal: createdReceptionTotal,
		addedProductTotal:     addedProductTotal,
	}
}

func (m *Metric) RequestsTotal(method, path string) {
	m.requestsTotal.WithLabelValues(method, path).Inc()
}

func (m *Metric) ResponseTime(method, path string, duration time.Duration) {
	m.responseTime.WithLabelValues(method, path).Observe(duration.Seconds())
}

func (m *Metric) CreatedPvzTotal(city string) {
	m.createdPvzTotal.WithLabelValues(city).Inc()
}

func (m *Metric) CreatedReceptionsTotal(pvzId string) {
	m.createdReceptionTotal.WithLabelValues(pvzId).Inc()
}

func (m *Metric) AddedProductTotal(productType string) {
	m.addedProductTotal.WithLabelValues(productType).Inc()
}
