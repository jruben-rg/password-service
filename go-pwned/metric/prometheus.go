package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusService struct {
	httpRequestHistogram *prometheus.HistogramVec
}

func NewPrometheusService() (*PrometheusService, error) {
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"path", "method", "code", "job"})

	s := &PrometheusService{
		httpRequestHistogram: http,
	}

	err := prometheus.Register(s.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

func (ps *PrometheusService) SaveMetrics(mi *MetricInfo) {
	ps.httpRequestHistogram.WithLabelValues(mi.Path, mi.Method, mi.StatusCode, mi.Job).Observe(mi.Duration)
}
