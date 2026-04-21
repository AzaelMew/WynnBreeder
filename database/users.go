package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"wynnbreeder/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUsernameTaken = errors.New("username already taken")

func (db *DB) CreateUser(username, passwordHash string, isAdmin bool) (*models.User, error) {
	return db.createUser(username, passwordHash, isAdmin, false)
}

func (db *DB) CreateSuperAdmin(username, passwordHash string) (*models.User, error) {
	return db.createUser(username, passwordHash, true, true)
}

func (db *DB) createUser(username, passwordHash string, isAdmin, isSuperAdmin bool) (*models.User, error) {
	res, err := db.Exec(
		`INSERT INTO users (username, password_hash, is_admin, is_superadmin) VALUES (?, ?, ?, ?)`,
		strings.ToLower(username), passwordHash, isAdmin, isSuperAdmin,
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	id, _ := res.LastInsertId()
	return db.GetUserByID(id)
}

func (db *DB) GetUserByID(id int64) (*models.User, error) {
	u := &models.User{}
	err := db.QueryRow(
		`SELECT id, username, password_hash, is_admin, is_superadmin, created_at FROM users WHERE id = ?`, id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.IsAdmin, &u.IsSuperAdmin, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	u := &models.User{}
	err := db.QueryRow(
		`SELECT id, username, password_hash, is_admin, is_superadmin, created_at FROM users WHERE username = ?`, strings.ToLower(username),
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.IsAdmin, &u.IsSuperAdmin, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (db *DB) ListUsers() ([]models.User, error) {
	rows, err := db.Query(
		`SELECT id, username, is_admin, is_superadmin, created_at FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsAdmin, &u.IsSuperAdmin, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (db *DB) PromoteToSuperAdmin(username string) (*models.User, error) {
	res, err := db.Exec(
		`UPDATE users SET is_admin = 1, is_superadmin = 1 WHERE username = ?`,
		strings.ToLower(username),
	)
	if err != nil {
		return nil, fmt.Errorf("promote user: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, ErrUserNotFound
	}
	u, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (db *DB) SetUserRole(id int64, isAdmin, isSuperAdmin bool) (*models.User, error) {
	res, err := db.Exec(
		`UPDATE users SET is_admin = ?, is_superadmin = ? WHERE id = ?`,
		isAdmin, isSuperAdmin, id,
	)
	if err != nil {
		return nil, fmt.Errorf("set user role: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, ErrUserNotFound
	}
	return db.GetUserByID(id)
}

func (db *DB) UpdatePassword(userID int64, newHash string) error {
	_, err := db.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, newHash, userID)
	return err
}

func (db *DB) DeleteUser(id int64) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (db *DB) CountAdmins() (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_admin = 1`).Scan(&count)
	return count, err
}

func isUniqueConstraint(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: users.username" ||
		containsStr(err.Error(), "UNIQUE constraint"))
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
