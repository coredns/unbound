package unbound

import (
	"sync"

	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics exported by the unbound plugin.
var (
	RequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: plugin.Namespace,
		Subsystem: "unbound",
		Name:      "request_duration_seconds",
		Buckets:   plugin.TimeBuckets,
		Help:      "Histogram of the time each request took.",
	})

	RcodeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "unbound",
		Name:      "response_rcode_count_total",
		Help:      "Counter of rcodes made per request.",
	}, []string{"rcode"})
)

var once sync.Once
