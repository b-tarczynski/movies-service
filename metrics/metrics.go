package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "movies_service"

	apiSubsystem = "api"
)

type Metrics struct {
	registry *prometheus.Registry

	APIRequestDuration prometheus.Histogram
}

func New() *Metrics {
	metrics := &Metrics{
		registry: prometheus.NewRegistry(),
	}

	metrics.APIRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: apiSubsystem,
		Name:      "request_duration_milliseconds",
	})

	metrics.registry.MustRegister(metrics.APIRequestDuration)

	return metrics
}
