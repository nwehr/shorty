package main

import (
	"context"

	"github.com/mediocregopher/radix/v4"
)

type idStore struct {
	rds *radix.Cluster
}

func (s idStore) getNextID() (string, error) {
	id := ""
	err := s.rds.Do(context.Background(), radix.Cmd(&id, "LPOP", "shorty:ids"))
	return id, err
}
