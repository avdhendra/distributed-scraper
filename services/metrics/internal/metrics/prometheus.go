package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	scrapeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scraper_requests_total",
			Help: "Total number of scrape requests",
		},
		[]string{"platform"},
	)
	scrapeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "scraper_duration_seconds",
			Help: "Duration of scrape operations",
		},
		[]string{"platform"},
	)
	scrapeErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scraper_errors_total",
			Help: "Total number of scrape errors",
		},
		[]string{"platform"},
	)
)

func Init() {
	prometheus.MustRegister(scrapeCounter, scrapeDuration, scrapeErrors)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9090", nil)
}