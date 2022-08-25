package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func main() {
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

	go s.startCountingReqs()

	if err := http.ListenAndServe(":8080", s); err != nil {
		fmt.Println(err)
	}
}
