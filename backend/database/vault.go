package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

func InitDB() error {
	url := os.Getenv("TURSO_DATABASE_URL")
	if url == "" {
		return fmt.Errorf("TURSO_DATABASE_URL not set")
	}

	var err error
	DB, err = sql.Open("libsql", url)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	return nil
}
