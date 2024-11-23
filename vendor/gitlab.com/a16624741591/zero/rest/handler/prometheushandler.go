package handler

import (
	"gitlab.com/a16624741591/zero/rest/internal/response"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/a16624741591/zero/core/metric"
	"gitlab.com/a16624741591/zero/core/prometheus"
	"gitlab.com/a16624741591/zero/core/timex"
)

const serverNamespace = "http_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "code", "method"},
	})
)

// PrometheusHandler returns a middleware that reports stats to prometheus.
func PrometheusHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !prometheus.Enabled() {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := timex.Now()
			cw := &response.WithCodeResponseWriter{Writer: w}
			defer func() {
				metricServerReqDur.Observe(int64(timex.Since(startTime)/time.Millisecond), r.URL.Path, r.Method)
				metricServerReqCodeTotal.Inc(r.URL.Path, strconv.Itoa(cw.Code), r.Method)
			}()

			next.ServeHTTP(cw, r)
		})
	}
}
