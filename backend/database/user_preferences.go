package database

import (
	"context"
	"time"
)

func CreateUserPreferenceTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS user_preference (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_preference TEXT DEFAULT NULL, 
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE INDEX IF NOT EXISTS idx_user_preference ON user_preference(id);
	`

	// Set timeout for table creation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx, query)
	return err
}
