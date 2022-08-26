package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mediocregopher/radix/v4"
)

func generateUniqueIds() {
	alphabet := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	id := ID{0, 0, 0, 0}
	end := ID{len(alphabet) - 1, len(alphabet) - 1, len(alphabet) - 1, len(alphabet) - 1}

	ids := []ID{}

	for {
		id = id.Inc(len(alphabet))
		ids = append(ids, id)

		if id == end {
			break
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	for _, id := range ids {
		fmt.Printf("RPUSH shorty:ids %s\n", id.String(alphabet))
	}
}

func importUniqueIds() error {
	opts, _ := getImportOptions(os.Args[1:])

	addrs := []string{
		"localhost:7000",
		"localhost:7001",
		"localhost:7002",
		"localhost:7003",
		"localhost:7004",
		"localhost:7005",
	}

	client, err := (radix.ClusterConfig{}).New(context.Background(), addrs)
	if err != nil {
		return err
	}

	defer client.Close()

	f, err := os.Open(opts.file)
	if err != nil {
		return err
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		err := client.Do(context.Background(), radix.Cmd(nil, "RPUSH", "shorty:ids", s.Text()))
		if err != nil {
			return err
		}
	}

	return nil
}

type importOptions struct {
	file  string
	redis string
}

func getImportOptions(args []string) (importOptions, error) {
	opts := importOptions{}

	for i, arg := range args {
		switch arg {
		case "--file":
			opts.file = args[i+1]
		case "--redis":
			opts.redis = args[i+1]
		}
	}

	return opts, nil
}

type ID [4]int

func (id ID) Inc(max int) ID {
	id[3]++

	if id[3] == max {
		id[3] = 0
		id[2]++
	}

	if id[2] == max {
		id[2] = 0
		id[1]++
	}

	if id[1] == max {
		id[1] = 0
		id[0]++
	}

	if id[0] == max {
		return ID{max, max, max, max}
	}

	return id
}

func (id ID) String(alphabet []byte) string {
	return string([]byte{alphabet[id[0]], alphabet[id[1]], alphabet[id[2]], alphabet[id[3]]})
}
