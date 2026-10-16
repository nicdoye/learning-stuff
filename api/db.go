package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// databaseURL returns the Postgres connection string from the environment.
func databaseURL() (string, error) {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url, nil
	}
	return "", errors.New("DATABASE_URL is not set: start the local Postgres and export it (see README.md)")
}

// connect opens a Postgres connection pool and verifies it with a ping.
func connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return pool, nil
}
