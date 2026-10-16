package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDatabaseURLMissing(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := databaseURL(); err == nil {
		t.Fatal("expected an error when DATABASE_URL is not set")
	}
}

// testPool connects to the database named by DATABASE_URL, skipping the test
// when it is not set.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	url, err := databaseURL()
	if err != nil {
		t.Skip("DATABASE_URL not set; skipping database test")
	}

	pool, err := connect(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestMigrate(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()

	if err := migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	// A second run must be a no-op.
	if err := migrate(ctx, pool); err != nil {
		t.Fatalf("migrate again: %v", err)
	}

	var tableExists bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables
		WHERE table_name = 'bookmarks')`).Scan(&tableExists); err != nil {
		t.Fatalf("check bookmarks table: %v", err)
	}
	if !tableExists {
		t.Error("bookmarks table does not exist after migrate")
	}

	var recorded bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS (
		SELECT 1 FROM schema_migrations
		WHERE version = '0001_create_bookmarks.sql')`).Scan(&recorded); err != nil {
		t.Fatalf("check schema_migrations: %v", err)
	}
	if !recorded {
		t.Error("0001_create_bookmarks.sql not recorded in schema_migrations")
	}
}
