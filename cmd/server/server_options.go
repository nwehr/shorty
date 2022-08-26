package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
)

type serverOptions struct {
	publicUrl   string
	postgresUrl string
	redisUrls   []string

	rds *radix.Cluster
	pg  *pgx.Conn

	oauthClientID     string
	oauthClientSecret string
	oauthAuthURL      string
	oauthTokenURL     string
	oauthRedirectURL  string
}

func getServerOptions(args []string) (serverOptions, error) {
	opts := serverOptions{
		publicUrl:         os.Getenv("PUBLIC_URL"),
		postgresUrl:       os.Getenv("POSTGRES_URL"),
		redisUrls:         strings.Split(os.Getenv("REDIS_URL"), ","),
		oauthClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		oauthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		oauthRedirectURL:  "http://localhost:48162/oauth/callback",
		oauthAuthURL:      "https://github.com/login/oauth/authorize",
		oauthTokenURL:     "https://github.com/login/oauth/access_token",
	}

	for i, arg := range args {
		switch arg {
		case "--public-url":
			opts.publicUrl = args[i+1]
		case "--postgres-url":
			opts.postgresUrl = args[i+1]
		case "--redis-url":
			opts.redisUrls = strings.Split(args[i+1], ",")
		case "--oauth-client-id":
			opts.oauthClientID = args[i+1]
		case "--oauth-client-secret":
			opts.oauthClientSecret = args[i+1]
		}
	}

	{
		pg, err := pgx.Connect(context.Background(), opts.postgresUrl)
		if err != nil {
			return opts, fmt.Errorf("unable to connect to postgres: %v", err)
		}

		opts.pg = pg
	}

	{
		rds, err := (radix.ClusterConfig{}).New(context.Background(), opts.redisUrls)
		if err != nil {
			return opts, fmt.Errorf("unable to connect to redis: %v", err)
		}

		opts.rds = rds
	}

	return opts, nil
}
