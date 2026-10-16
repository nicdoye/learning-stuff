package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run connects to Postgres, applies migrations, and serves the API.
func run() error {
	ctx := context.Background()

	url, err := databaseURL()
	if err != nil {
		return err
	}

	pool, err := connect(ctx, url)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := migrate(ctx, pool); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	addr := ":" + port()
	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, newHandler())
}

// port returns the port to listen on, defaulting to 8080.
func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// newHandler returns the HTTP routes served by the API.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	return mux
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
