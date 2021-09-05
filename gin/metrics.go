package gin

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//InitMetrics

var (
	// histogram 桶分布
	probixBuckets = []float64{.00001, .0001, .001, .005, .01, .05, .1, .25, .5, 1, 2.5, 5, 10}
	// the upstream library supports it.
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "probix_apiserver_request_count",
			Help: "Counter of probix apiserver requests broken out for each  method, resource and code.",
		},
		[]string{"method", "resource", "code"},
	)
	RequestErrCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "probix_apiserver_requestErrors_total",
			Help: "Counter of probix apiserver requests broken out for each method, resource and code.",
		},
		[]string{"method", "resource", "code"},
	)
	RequestLatencies = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "probix_apiserver_request_latencies_seconds",
			Help: "Response latency distribution in microseconds for method, resource and code.",
			// 1ms 到10s进行分桶
			Buckets: probixBuckets,
		},
		[]string{"method", "resource", "code"},
	)
	RequestLatenciesSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "probix_apiserver_request_latencies_seconds_summary",
			Help: "Response latency summary in microseconds for each  method, resource and code.",
			// Make the sliding window of 10m.
			MaxAge:     10 * time.Minute,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "resource", "code"},
	)
)

func MonitorCount(resource, method, code string) {
	RequestCounter.With(prometheus.Labels{"method": method, "resource": resource, "code": code}).Inc()
}

func MonitorErrCount(resource, method, code string) {
	RequestErrCounter.With(prometheus.Labels{"method": method, "resource": resource, "code": code}).Inc()
}

func MonitorHistogram(resource, method, code string, start time.Time) {
	RequestLatencies.With(prometheus.Labels{"method": method, "resource": resource, "code": code}).Observe(time.Since(start).Seconds())
}

func MonitorSummary(resource, method, code string, start time.Time) {
	RequestLatenciesSummary.With(prometheus.Labels{"method": method, "resource": resource, "code": code}).Observe(time.Since(start).Seconds())
}
