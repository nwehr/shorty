package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	reqeustsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shorty_requests_total",
		Help: "The total number of requests",
	})

	cacheMissesCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shorty_cache_misses_total",
		Help: "The total number of cache misses",
	})

	errorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "shorty_errors_total",
		Help: "The total number of errors",
	})
)
