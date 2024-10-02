package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metrics

func Init(appName string) {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: appName + "_requests_total",
			Help: "count of requests to service",
		}),
		responseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: appName + "_responses_total",
			Help: "count of responses from service",
		}, []string{"status", "path"}),
		histogramResponseTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    appName + "_histogram_response_time_seconds",
			Help:    "histogram of response time from service",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 16),
		}, []string{"status", "path"}),
	}
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status int, path string) {
	metrics.responseCounter.WithLabelValues(fmt.Sprint(status), path).Inc()
}

func HistogramsResponseTimeObserve(status int, path string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(fmt.Sprint(status), path).Observe(time)
}
