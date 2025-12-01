package database

import (
	"context"
	"time"
)

func CreateUsersTrackingTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users_tracking (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL UNIQUE,
		meal_count INTEGER NOT NULL DEFAULT 0,
		max_meals INTEGER NOT NULL DEFAULT 20,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_users_tracking_user_id ON users_tracking(user_id);
	`

	// Set timeout for table creation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx, query)
	return err
}
