package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Open(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db := &DB{sqlDB}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func (db *DB) migrate() error {
	// v1: initial schema
	v1 := `
CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    submission_id INTEGER NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK(role IN ('parent_a', 'parent_b', 'offspring')),
    type TEXT NOT NULL,
    potential INTEGER NOT NULL,
    color TEXT NOT NULL,
    name TEXT NOT NULL,
    energy_value INTEGER NOT NULL,
    energy_max INTEGER NOT NULL,
    speed_val INTEGER NOT NULL DEFAULT 0,
    speed_lim INTEGER NOT NULL DEFAULT 0,
    speed_max INTEGER NOT NULL DEFAULT 0,
    accel_val INTEGER NOT NULL DEFAULT 0,
    accel_lim INTEGER NOT NULL DEFAULT 0,
    accel_max INTEGER NOT NULL DEFAULT 0,
    altitude_val INTEGER NOT NULL DEFAULT 0,
    altitude_lim INTEGER NOT NULL DEFAULT 0,
    altitude_max INTEGER NOT NULL DEFAULT 0,
    energy_stat_val INTEGER NOT NULL DEFAULT 0,
    energy_stat_lim INTEGER NOT NULL DEFAULT 0,
    energy_stat_max INTEGER NOT NULL DEFAULT 0,
    handling_val INTEGER NOT NULL DEFAULT 0,
    handling_lim INTEGER NOT NULL DEFAULT 0,
    handling_max INTEGER NOT NULL DEFAULT 0,
    toughness_val INTEGER NOT NULL DEFAULT 0,
    toughness_lim INTEGER NOT NULL DEFAULT 0,
    toughness_max INTEGER NOT NULL DEFAULT 0,
    boost_val INTEGER NOT NULL DEFAULT 0,
    boost_lim INTEGER NOT NULL DEFAULT 0,
    boost_max INTEGER NOT NULL DEFAULT 0,
    training_val INTEGER NOT NULL DEFAULT 0,
    training_lim INTEGER NOT NULL DEFAULT 0,
    training_max INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_mounts_submission ON mounts(submission_id);
CREATE INDEX IF NOT EXISTS idx_mounts_role ON mounts(role);
CREATE INDEX IF NOT EXISTS idx_submissions_user ON submissions(user_id);
`
	if _, err := db.Exec(v1); err != nil {
		return fmt.Errorf("v1 migration: %w", err)
	}

	// v2: add status column to submissions (pending/complete)
	var schemaVer int
	_ = db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&schemaVer)
	if schemaVer < 2 {
		v2 := `
ALTER TABLE submissions ADD COLUMN status TEXT NOT NULL DEFAULT 'complete';
CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status);
INSERT OR REPLACE INTO schema_version (version) VALUES (2);
`
		if _, err := db.Exec(v2); err != nil {
			return fmt.Errorf("v2 migration: %w", err)
		}
	}

	return nil
}
