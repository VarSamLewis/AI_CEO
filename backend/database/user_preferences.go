package database

import (
	"context"
	"time"
)

func CreateUserPreferenceTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS user_preference (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL UNIQUE,
		user_preference TEXT DEFAULT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_user_preference_user_id ON user_preference(user_id);
	`

	// Set timeout for table creation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx, query)
	return err
}
