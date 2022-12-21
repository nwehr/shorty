package models

import (
	"context"

	"github.com/mediocregopher/radix/v4"
)

type KeyStore struct {
	RDS *radix.Cluster
}

func (s KeyStore) Next() (string, error) {
	key := ""
	err := s.RDS.Do(context.Background(), radix.Cmd(&key, "LPOP", "shorty:keys"))
	return key, err
}
