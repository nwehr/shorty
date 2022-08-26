package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
)

type shortUrl struct {
	ID  string
	URL string
}

type urlStore struct {
	pg  *pgx.Conn
	rds *radix.Cluster
}

func (s urlStore) getUrl(id string) (string, error) {
	url := ""
	err := s.rds.Do(context.Background(), radix.Cmd(&url, "GET", fmt.Sprintf("shorty:cache:%s", id)))
	if err != nil {
		return url, err
	}

	if url == "" {
		q := `
			select url from urls where id = $1
		`

		err = s.pg.QueryRow(context.Background(), q, id).Scan(&url)
		if err != nil {
			fmt.Printf("could not query database: %s\n", err.Error())
			return "", err
		}

		err = s.rds.Do(context.Background(), radix.Cmd(nil, "SET", fmt.Sprintf("shorty:cache:%s", id), url))
		if err != nil {
			fmt.Printf("could not set cache: %s\n", err.Error())
		}
	}

	return url, nil
}

func (s urlStore) setUrl(userId int64, id, url string) error {
	q := `
		insert into urls (id, url, user_id) values ($1, $2, $3)
	`
	_, err := s.pg.Exec(context.Background(), q, id, url, userId)
	return err
}

func (s urlStore) listUrlsByUserId(userId int64) ([]shortUrl, error) {
	urls := []shortUrl{}

	q := `
		select id, url from urls where user_id = $1
	`

	rows, err := s.pg.Query(context.Background(), q, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		url := shortUrl{}

		rows.Scan(&url.ID, &url.URL)

		urls = append(urls, url)
	}

	return urls, nil
}
