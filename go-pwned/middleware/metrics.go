package middleware

import (
	"fmt"
	"net/http"

	"github.com/jruben-rg/password-service/go-pwned/metric"
)

const (
	Job = "pwned"
)

type statusResponseWriter struct {
	http.ResponseWriter
	code int
}

type MetricService interface {
	SaveMetrics(mi *metric.MetricInfo)
}

func NewStatusResponseWriter(rw http.ResponseWriter) *statusResponseWriter {
	// default status code is http.StatusOK
	return &statusResponseWriter{rw, http.StatusOK}
}

func (lrw *statusResponseWriter) WriteHeader(statusCode int) {
	lrw.code = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func Metrics(service MetricService, handler http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		appMetric := metric.NewMetricInfo(r.URL.Path, r.Method, Job)
		appMetric.Started()
		lrw := NewStatusResponseWriter(rw)
		handler.ServeHTTP(lrw, r)
		appMetric.Finished()
		appMetric.StatusCode = fmt.Sprint(lrw.code)
		service.SaveMetrics(appMetric)
	})

}
