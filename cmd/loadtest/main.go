package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	opts, _ := getLoadTestOptions(os.Args[1:])

	numReqs := int64(0)
	lastReqs := int64(0)

	numErrors := int64(0)
	lastErrors := int64(0)

	done := false

	go func() {
		for i := 0; i < opts.duration; i++ {
			time.Sleep(time.Second)

			fmt.Printf("%d req/s, %d err/s\n", (numReqs - lastReqs), (numErrors - lastErrors))

			lastReqs = numReqs
			lastErrors = numErrors
		}

		done = true
	}()

	wg := &sync.WaitGroup{}
	wg.Add(opts.clients)

	for i := 0; i < opts.clients; i++ {
		go func() {
			defer wg.Done()

			for {
				if done {
					return
				}

				start := time.Now()

				transport, _ := http.DefaultTransport.(*http.Transport)
				transport.MaxIdleConns = 100
				transport.MaxConnsPerHost = 100

				client := &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
					Transport: transport,
				}

				resp, err := client.Get(opts.url)
				if err != nil {
					atomic.AddInt64(&numErrors, 1)
				} else {
					io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()
				}

				atomic.AddInt64(&numReqs, 1)

				if opts.rate > 0 {
					sleep := (int64(1_000_000) / int64(opts.rate)) - time.Since(start).Microseconds()
					if sleep > 0 {
						time.Sleep(time.Duration(sleep) * time.Microsecond)
					}
				}
			}
		}()
	}

	wg.Wait()
}

type loadTestOptions struct {
	url      string
	clients  int
	duration int
	rate     int
}

func getLoadTestOptions(args []string) (loadTestOptions, error) {
	opts := loadTestOptions{
		clients:  1,
		duration: 5,
	}

	for i, arg := range args {
		switch arg {
		case "--url":
			opts.url = args[i+1]
		case "--clients":
			opts.clients, _ = strconv.Atoi(args[i+1])
		case "--duration":
			opts.duration, _ = strconv.Atoi(args[i+1])
		case "--rate":
			opts.rate, _ = strconv.Atoi(args[i+1])
		}
	}

	return opts, nil
}
