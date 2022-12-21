package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
)

type (
	URLMap struct {
		Key     string `json:"key"`
		LongURL string `json:"long_url"`
		Visits  uint   `json:"visits"`
	}

	URLMapStore struct {
		PG  *pgx.Conn
		RDS *radix.Cluster
	}
)

func (s URLMapStore) Get(key string) (URLMap, error) {
	m := URLMap{
		Key: key,
	}

	// err := s.RDS.Do(context.Background(), radix.Cmd(&m.LongURL, "GET", fmt.Sprintf("shorty:cache:%s", key)))
	// if err != nil {
	// 	return m, err
	// }

	if m.LongURL == "" {
		// cacheMissesCounter.Inc()

		q := `
			select url, visits from urls where key = $1
		`

		err := s.PG.QueryRow(context.Background(), q, key).Scan(&m.LongURL, &m.Visits)
		if err != nil {
			fmt.Printf("could not query database: %s\n", err.Error())
			return m, err
		}

		// err = s.RDS.Do(context.Background(), radix.Cmd(nil, "SET", fmt.Sprintf("shorty:cache:%s", key), m.LongURL))
		// if err != nil {
		// 	fmt.Printf("could not set cache: %s\n", err.Error())
		// }
	}

	return m, nil
}

func (s URLMapStore) Set(issuer string, m URLMap) (URLMap, error) {
	q := `
		update urls set url = $1, issuer = $2 where key = (select key from urls where url is null limit 1) returning key
	`
	err := s.PG.QueryRow(context.Background(), q, m.LongURL, issuer).Scan(&m.Key)
	return m, err
}

func (s URLMapStore) List(issuer string) ([]URLMap, error) {
	q := `
		select key, url, visits from urls where issuer = $1
	`

	rows, err := s.PG.Query(context.Background(), q, issuer)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	urls := []URLMap{}

	for rows.Next() {
		key := ""
		longURL := ""
		visits := uint(0)

		err = rows.Scan(&key, &longURL, &visits)
		if err != nil {
			return urls, err
		}

		urls = append(urls, URLMap{
			Key:     key,
			LongURL: longURL,
			Visits:  visits,
		})
	}

	return urls, nil
}
