package metric

import (
	"time"
)

type MetricInfo struct {
	Path       string
	Method     string
	StatusCode string
	Job        string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

func NewMetricInfo(path string, method string, job string) *MetricInfo {
	return &MetricInfo{
		Path:   path,
		Method: method,
		Job:    job,
	}
}

func (h *MetricInfo) Started() {
	h.StartedAt = time.Now()
}

func (h *MetricInfo) Finished() {
	h.FinishedAt = time.Now()
	h.Duration = time.Since(h.StartedAt).Seconds()
}
