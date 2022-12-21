package options

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
)

type Options struct {
	PublicURL   string
	PostgresURL string
	RedisURLs   []string

	RDS *radix.Cluster
	PG  *pgx.Conn

	ClientID         string
	ClientSecret     string
	AuthorizationURL string
	TokenURL         string
	LogoutURL        string
	RedirectURL      string

	Logger logger
}

func Get(args []string) (Options, error) {
	opts := Options{
		PublicURL:        os.Getenv("PUBLIC_URL"),
		PostgresURL:      os.Getenv("POSTGRES_URL"),
		RedisURLs:        strings.Split(os.Getenv("REDIS_URLS"), ","),
		ClientID:         os.Getenv("CLIENT_ID"),
		ClientSecret:     os.Getenv("CLIENT_SECRET"),
		RedirectURL:      coalesce(os.Getenv("REDIRECT_URL"), "http://localhost:48162/oauth/callback"),
		AuthorizationURL: coalesce(os.Getenv("AUTHORIZATION_URL"), "https://github.com/login/oauth/authorize"),
		TokenURL:         coalesce(os.Getenv("TOKEN_URL"), "https://github.com/login/oauth/access_token"),
		LogoutURL:        coalesce(os.Getenv("LOGOUT_URL")),
		Logger:           logger{},
	}

	for i, arg := range args {
		switch arg {
		case "--publicUrl":
			opts.PublicURL = args[i+1]
		case "--postgresUrl":
			opts.PostgresURL = args[i+1]
		case "--redisUrls":
			opts.RedisURLs = strings.Split(args[i+1], ",")
		case "--clientId":
			opts.ClientID = args[i+1]
		case "--clientSecret":
			opts.ClientSecret = args[i+1]
		case "--redirectUrl":
			opts.RedirectURL = args[i+1]
		case "--authorizationUrl":
			opts.AuthorizationURL = args[i+1]
		case "--tokenUrl":
			opts.TokenURL = args[i+1]
		case "--logoutUrl":
			opts.LogoutURL = args[i+1]
		case "--log":
			switch args[i+1] {
			case "debug":
				opts.Logger = logger{DEBUG}
			case "warning":
				opts.Logger = logger{WARNNING}
			case "error":
				opts.Logger = logger{ERROR}
			default:
				opts.Logger = logger{DEBUG}
			}
		}
	}

	{
		pg, err := pgx.Connect(context.Background(), opts.PostgresURL)
		if err != nil {
			return opts, fmt.Errorf("unable to connect to postgres: %v", err)
		}

		opts.PG = pg
	}

	{
		rds, err := (radix.ClusterConfig{}).New(context.Background(), opts.RedisURLs)
		if err != nil {
			return opts, fmt.Errorf("unable to connect to redis: %v", err)
		}

		opts.RDS = rds
	}

	return opts, nil
}

func coalesce(s ...string) string {
	for _, str := range s {
		if str != "" {
			return str
		}
	}

	return ""
}
