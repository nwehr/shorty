package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
)

type (
	mappedURL struct {
		ShortURL string `json:"short_url"`
		LongURL  string `json:"long_url"`
	}

	urlStore struct {
		pg  *pgx.Conn
		rds *radix.Cluster
	}
)

func (s urlStore) getURL(key string) (string, error) {
	url := ""
	err := s.rds.Do(context.Background(), radix.Cmd(&url, "GET", fmt.Sprintf("shorty:cache:%s", key)))
	if err != nil {
		return url, err
	}

	if url == "" {
		cacheMissesCounter.Inc()

		q := `
			select url from urls where key = $1
		`

		err = s.pg.QueryRow(context.Background(), q, key).Scan(&url)
		if err != nil {
			fmt.Printf("could not query database: %s\n", err.Error())
			return "", err
		}

		err = s.rds.Do(context.Background(), radix.Cmd(nil, "SET", fmt.Sprintf("shorty:cache:%s", key), url))
		if err != nil {
			fmt.Printf("could not set cache: %s\n", err.Error())
		}
	}

	return url, nil
}

func (s urlStore) setURL(userID int64, key, url string) error {
	q := `
		insert into urls (key, url, user_id) values ($1, $2, $3)
	`
	_, err := s.pg.Exec(context.Background(), q, key, url, userID)
	return err
}

func (s urlStore) listURLsByUserID(userID int64, publicURL string) ([]mappedURL, error) {
	q := `
		select key, url from urls where user_id = $1
	`

	rows, err := s.pg.Query(context.Background(), q, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	urls := []mappedURL{}

	for rows.Next() {
		key := ""
		longURL := ""

		rows.Scan(&key, &longURL)

		urls = append(urls, mappedURL{
			ShortURL: publicURL + "/" + key,
			LongURL:  longURL,
		})
	}

	return urls, nil
}
