package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	opts, err := getServerOptions(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	s, err := newServer(opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer s.opts.pg.Close(context.Background())
	defer s.opts.rds.Close()

	if err := http.ListenAndServe(":8080", s); err != nil {
		fmt.Println(err)
	}
}
