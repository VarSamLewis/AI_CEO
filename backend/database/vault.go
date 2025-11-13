package database

import (
	"database/sql"
	"fmt"
	"os"

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

  // Test connection
  if err := DB.Ping(); err != nil {
  	return fmt.Errorf("failed to ping db: %w", err)
  }

  return nil
  
}
