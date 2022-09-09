package main

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type user struct {
	ID          int64
	Login       string
	Name        string
	AccessToken string
}

type userStore struct {
	pg *pgx.Conn
}

func (s userStore) getUserByAccessToken(accessToken string) (user, error) {
	q := `
		select id, login, name from users where access_token = $1
	`

	u := user{
		AccessToken: accessToken,
	}

	err := s.pg.QueryRow(context.Background(), q, accessToken).Scan(&u.ID, &u.Login, &u.Name)

	return u, err
}

func (s userStore) getUserByID(id int64) (user, error) {
	q := `
		select login, name, access_token from users where id = $1
	`

	u := user{
		ID: id,
	}

	err := s.pg.QueryRow(context.Background(), q, id).Scan(&u.Login, &u.Name, &u.AccessToken)

	return u, err
}

func (s userStore) setUser(u user) error {
	_, err := s.getUserByID(u.ID)
	if err != nil {
		q := `
			insert into users (id, login, name, access_token) values ($1, $2, $3, $4)
		`

		_, err = s.pg.Exec(context.Background(), q, u.ID, u.Login, u.Name, u.AccessToken)
		return err
	}

	q := `
		update users set login = $2, name = $3, access_token = $4 where id = $1
	`

	_, err = s.pg.Exec(context.Background(), q, u.ID, u.Login, u.Name, u.AccessToken)
	return err
}
