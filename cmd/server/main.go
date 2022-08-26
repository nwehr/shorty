package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("reading options...")
	opts, err := getServerOptions(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("creating server...")
	s, err := newServer(opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer s.opts.pg.Close(context.Background())
	defer s.opts.rds.Close()

	fmt.Println("starting req counter...")
	go s.startCountingReqs()

	fmt.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", s); err != nil {
		fmt.Println(err)
	}
}
