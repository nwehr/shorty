package main

import (
	"context"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func main() {
	opts, err := getOptions(os.Args[1:])
	if err != nil {
		opts.logger.Errorf("could not get server options: %v\n", err)
	}

	go func() {
		opts.logger.Debugf("starting metrics server on :2112\n")

		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	defer opts.pg.Close(context.Background())
	defer opts.rds.Close()

	mux := http.NewServeMux()

	mux.Handle("/auth/", newAuthHandler(opts))
	mux.Handle("/api/", newAPIHandler(opts))
	mux.Handle("/", redirectHandler(opts))

	opts.logger.Debugf("starting server on :8080\n")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		opts.logger.Fatalf("could not listen: %v\n", err)
	}
}
