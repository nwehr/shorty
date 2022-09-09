package main

import (
	"fmt"
	"math/rand"
	"time"
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
