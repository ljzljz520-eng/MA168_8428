package api

import (
	"net/http"
	"sync/atomic"
)

type Metrics struct {
	requests atomic.Int64
	failures atomic.Int64
}

func (m *Metrics) Observe(status int) {
	m.requests.Add(1)
	if status >= 400 {
		m.failures.Add(1)
	}
}

func (m *Metrics) Snapshot() map[string]int64 {
	return map[string]int64{"requests": m.requests.Load(), "failures": m.failures.Load()}
}

func MetricsHandler(metrics *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { encodeEnvelope(w, http.StatusOK, metrics.Snapshot()) }
}
