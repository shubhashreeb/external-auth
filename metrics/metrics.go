package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	loginReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "number_of_logins_received",
			Help:      "Number of logins received",
		})

	logoutReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "number_of_logouts_received",
			Help:      "Number of logouts received",
		})

	ratelimited = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "xds",
			Name:      "number_of_request_rate_limited",
			Help:      "number of request rejected due to rate limit",
		})

	invalidLoginReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "invalid_logins_received",
			Help:      "Number of invalid logins received",
		})

	cacheLoginToken = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "cache_login_token",
			Help:      "Number of token logins cached",
		})

	loginCacheHit = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "login_cache_hit",
			Help:      "Number of login cache hit",
		})

	loginCacheMissed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "login_cache_missed",
			Help:      "Number of login cache missed",
		})

	domains = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "xds",
			Name:      "domains",
			Help:      "number of clusters served",
		})

	eventsReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "number_of_event_received",
			Help:      "Number of event received",
		})

	eventsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "k8s",
			Name:      "number_of_event_processed",
			Help:      "Number of event processed",
		})

	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		})

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)

type Metrics struct {
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RunPrometheusServer() {
	fmt.Println("Running Prometheus server on port 8880")
	rand.Seed(time.Now().Unix())

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "prom_request_time",
		Help: "Time it has taken to retrieve the metrics",
	}, []string{"time"})

	prometheus.Register(histogramVec)

	http.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), histogramVec))

	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(loginReceived)
	prometheus.MustRegister(logoutReceived)
	prometheus.MustRegister(ratelimited)
	prometheus.MustRegister(invalidLoginReceived)
	prometheus.MustRegister(cacheLoginToken)
	prometheus.MustRegister(loginCacheHit)
	prometheus.MustRegister(loginCacheMissed)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func newHandlerWithHistogram(handler http.Handler, histogram *prometheus.HistogramVec) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		status := http.StatusOK

		defer func() {
			histogram.WithLabelValues(fmt.Sprintf("%d", status)).Observe(time.Since(start).Seconds())
		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}

func (m *Metrics) AddCounterStats(name string, value float64) {
	switch name {
	case "ratelimited":
		ratelimited.Add(value)
	case "invalidLoginReceived":
		invalidLoginReceived.Add(value)
	case "loginReceived":
		loginReceived.Add(value)
	case "logoutReceived":
		logoutReceived.Add(value)
	case "loginCacheHit":
		loginCacheHit.Add(value)
	case "cacheLoginToken":
		cacheLoginToken.Add(value)
	case "loginCacheMissed":
		loginCacheMissed.Add(1)
		// case "loginCacheHit":
		// 	loginCacheHit.Add(value)
		// case "loginCacheHit":
		// 	loginCacheHit.Add(value)
	}
}
