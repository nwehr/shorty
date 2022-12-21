package main

import (
	"context"
	"net/http"
	"os"

	"github.com/nwehr/shorty/cmd/server/handlers"
	"github.com/nwehr/shorty/cmd/server/options"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	opts, err := options.Get(os.Args[1:])
	if err != nil {
		opts.Logger.Errorf("could not get server options: %v\n", err)
	}

	go func() {
		opts.Logger.Debugf("starting metrics server on :2112\n")

		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	defer opts.PG.Close(context.Background())
	defer opts.RDS.Close()

	mux := http.NewServeMux()

	mux.Handle("/auth/", handlers.AuthHandler{Opts: opts})
	mux.Handle("/api/", handlers.NewAPIHandler(opts))
	mux.Handle("/admin", handlers.Admin(opts))
	mux.Handle("/", handlers.Redirect(opts))

	opts.Logger.Debugf("starting server on :8080\n")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		opts.Logger.Fatalf("could not listen: %v\n", err)
	}
}
