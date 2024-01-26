package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type PrometheusMetrics interface {
	SetupMetrics() error
	ExecutionTime() *prometheus.HistogramVec
	SpanTime() *prometheus.HistogramVec
	ErrorsHits() *prometheus.CounterVec
	SuccessHits() *prometheus.CounterVec
	TotalHits() prometheus.Counter
}

type prometheusMetrics struct {
	executionTime *prometheus.HistogramVec
	spanTime      *prometheus.HistogramVec
	errorsHits    *prometheus.CounterVec
	successHits   *prometheus.CounterVec
	totalHits     prometheus.Counter
}

func NewPrometheusMetrics(serviceName string) PrometheusMetrics {
	metrics := &prometheusMetrics{
		executionTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: serviceName + "_durations",
			Help: "Shows durations in minutes of request execution",
		}, []string{"status", "path", "method"}),
		spanTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: serviceName + "_span",
			Help: "Shows durations of SPAN",
		}, []string{"path", "method"}),
		errorsHits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: serviceName + "_errors_hits",
			Help: "Counts errors responses from service",
		}, []string{"status", "path", "method"}),
		successHits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: serviceName + "_success_hits",
			Help: "Counts success responses from service",
		}, []string{"status", "path", "method"}),
		totalHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: serviceName + "_total_hits",
			Help: "Counts all responses from service",
		}),
	}

	return metrics
}

func (m *prometheusMetrics) SetupMetrics() error {
	if err := prometheus.Register(m.executionTime); err != nil {
		return err
	}

	if err := prometheus.Register(m.spanTime); err != nil {
		return err
	}

	if err := prometheus.Register(m.errorsHits); err != nil {
		return err
	}

	if err := prometheus.Register(m.successHits); err != nil {
		return err
	}

	if err := prometheus.Register(m.totalHits); err != nil {
		return err
	}

	return nil
}

func (m *prometheusMetrics) ExecutionTime() *prometheus.HistogramVec {
	return m.executionTime
}

func (m *prometheusMetrics) SpanTime() *prometheus.HistogramVec {
	return m.spanTime
}

func (m *prometheusMetrics) ErrorsHits() *prometheus.CounterVec {
	return m.errorsHits
}
func (m *prometheusMetrics) SuccessHits() *prometheus.CounterVec {
	return m.successHits
}
func (m *prometheusMetrics) TotalHits() prometheus.Counter {
	return m.totalHits
}

func ServePrometheusHTTP(addr string) {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
