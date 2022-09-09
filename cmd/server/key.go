package main

import (
	"context"

	"github.com/mediocregopher/radix/v4"
)

type keyStore struct {
	rds *radix.Cluster
}

func (s keyStore) getNextKey() (string, error) {
	key := ""
	err := s.rds.Do(context.Background(), radix.Cmd(&key, "LPOP", "shorty:keys"))
	return key, err
}
