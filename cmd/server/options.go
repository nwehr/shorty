package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
)

type options struct {
	publicURL   string
	postgresURL string
	redisURLs   []string

	rds *radix.Cluster
	pg  *pgx.Conn

	oauthClientID     string
	oauthClientSecret string
	oauthAuthURL      string
	oauthTokenURL     string
	oauthRedirectURL  string

	logger logger
}

func getOptions(args []string) (options, error) {
	opts := options{
		publicURL:         os.Getenv("PUBLIC_URL"),
		postgresURL:       os.Getenv("POSTGRES_URL"),
		redisURLs:         strings.Split(os.Getenv("REDIS_URL"), ","),
		oauthClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		oauthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		oauthRedirectURL:  "http://localhost:48162/oauth/callback",
		oauthAuthURL:      "https://github.com/login/oauth/authorize",
		oauthTokenURL:     "https://github.com/login/oauth/access_token",
		logger:            logger{},
	}

	for i, arg := range args {
		switch arg {
		case "--public-url":
			opts.publicURL = args[i+1]
		case "--postgres-url":
			opts.postgresURL = args[i+1]
		case "--redis-url":
			opts.redisURLs = strings.Split(args[i+1], ",")
		case "--oauth-client-id":
			opts.oauthClientID = args[i+1]
		case "--oauth-client-secret":
			opts.oauthClientSecret = args[i+1]
		case "--log":
			switch args[i+1] {
			case "debug":
				opts.logger = logger{DEBUG}
			case "warning":
				opts.logger = logger{WARNNING}
			case "error":
				opts.logger = logger{ERROR}
			default:
				opts.logger = logger{DEBUG}
			}
		}
	}

	{
		pg, err := pgx.Connect(context.Background(), opts.postgresURL)
		if err != nil {
			return opts, fmt.Errorf("unable to connect to postgres: %v", err)
		}

		opts.pg = pg
	}

	{
		rds, err := (radix.ClusterConfig{}).New(context.Background(), opts.redisURLs)
		if err != nil {
			return opts, fmt.Errorf("unable to connect to redis: %v", err)
		}

		opts.rds = rds
	}

	return opts, nil
}
