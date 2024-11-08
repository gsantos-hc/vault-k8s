package agent_inject

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type InjectorMetrics struct {
	requestQueue                prometheus.Gauge
	requestProcessingTime       prometheus.Summary
	injectionsByNamespace       *prometheus.CounterVec
	failedInjectionsByNamespace *prometheus.CounterVec
}

var metrics *InjectorMetrics

func InitMetrics(reg prometheus.Registerer) {
	const (
		namespace string = "vault"
		subsystem string = "sidecar_injector"
	)

	// Initialize Agent Sidecar Injector metrics
	metrics = &InjectorMetrics{
		requestQueue: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_queue_total",
			Help:      "Total count of webhook requests in the queue",
		}),

		requestProcessingTime: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_processing_duration_ms",
			Help:      "Summary of webhook request processing times in milliseconds (last 5 minutes)",
			MaxAge:    time.Duration(5 * time.Minute),
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.01,
				0.99: 0.001,
			},
		}),

		injectionsByNamespace: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "injections_by_namespace_total",
				Help:      "Total count of Agent Sidecar injections by namespace",
			},
			[]string{"namespace"},
		),

		failedInjectionsByNamespace: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "failed_injections_by_namespace_total",
				Help:      "Total count of failed Agent Sidecar injections by namespace",
			},
			[]string{"namespace"},
		),
	}

	// Register metrics in global registry
	reg.MustRegister(
		metrics.requestQueue,
		metrics.requestProcessingTime,
		metrics.injectionsByNamespace,
		metrics.failedInjectionsByNamespace,
	)
}
