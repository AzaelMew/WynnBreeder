package database

import (
	"database/sql"
	"errors"
	"fmt"

	"wynnmounts/models"
)

var ErrSubmissionNotFound = errors.New("submission not found")
var ErrSubmissionAlreadyComplete = errors.New("submission already has offspring")

func (db *DB) CreateSubmission(userID int64, notes string, mounts []models.Mount, status string) (*models.Submission, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO submissions (user_id, notes, status) VALUES (?, ?, ?)`, userID, notes, status)
	if err != nil {
		return nil, fmt.Errorf("insert submission: %w", err)
	}
	subID, _ := res.LastInsertId()

	for i := range mounts {
		mounts[i].SubmissionID = subID
		if err := insertMount(tx, mounts[i]); err != nil {
			return nil, fmt.Errorf("insert mount: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return db.GetSubmission(subID)
}

// AddOffspring adds the offspring mount to a pending submission and marks it complete.
// Only the owner (or admin, enforced in handler) may call this.
func (db *DB) AddOffspring(subID int64, offspring models.Mount) (*models.Submission, error) {
	// Verify submission is still pending
	var status string
	err := db.QueryRow(`SELECT status FROM submissions WHERE id = ?`, subID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSubmissionNotFound
	}
	if err != nil {
		return nil, err
	}
	if status == "complete" {
		return nil, ErrSubmissionAlreadyComplete
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	offspring.SubmissionID = subID
	if err := insertMount(tx, offspring); err != nil {
		return nil, fmt.Errorf("insert offspring: %w", err)
	}
	if _, err := tx.Exec(`UPDATE submissions SET status = 'complete' WHERE id = ?`, subID); err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return db.GetSubmission(subID)
}

func insertMount(tx *sql.Tx, m models.Mount) error {
	_, err := tx.Exec(`
INSERT INTO mounts (
    submission_id, role, type, potential, color, name, energy_value, energy_max,
    speed_val, speed_lim, speed_max,
    accel_val, accel_lim, accel_max,
    altitude_val, altitude_lim, altitude_max,
    energy_stat_val, energy_stat_lim, energy_stat_max,
    handling_val, handling_lim, handling_max,
    toughness_val, toughness_lim, toughness_max,
    boost_val, boost_lim, boost_max,
    training_val, training_lim, training_max
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?
)`,
		m.SubmissionID, m.Role, m.Type, m.Potential, m.Color, m.Name, m.EnergyValue, m.EnergyMax,
		m.SpeedVal, m.SpeedLim, m.SpeedMax,
		m.AccelVal, m.AccelLim, m.AccelMax,
		m.AltitudeVal, m.AltitudeLim, m.AltitudeMax,
		m.EnergyStatVal, m.EnergyStatLim, m.EnergyStatMax,
		m.HandlingVal, m.HandlingLim, m.HandlingMax,
		m.ToughnessVal, m.ToughnessLim, m.ToughnessMax,
		m.BoostVal, m.BoostLim, m.BoostMax,
		m.TrainingVal, m.TrainingLim, m.TrainingMax,
	)
	return err
}

func (db *DB) GetSubmission(id int64) (*models.Submission, error) {
	s := &models.Submission{}
	err := db.QueryRow(`
SELECT s.id, s.user_id, u.username, s.notes, s.status, s.created_at
FROM submissions s JOIN users u ON s.user_id = u.id
WHERE s.id = ?`, id,
	).Scan(&s.ID, &s.UserID, &s.Username, &s.Notes, &s.Status, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSubmissionNotFound
	}
	if err != nil {
		return nil, err
	}

	mounts, err := db.getMountsForSubmission(id)
	if err != nil {
		return nil, err
	}
	s.Mounts = mounts
	return s, nil
}

func (db *DB) getMountsForSubmission(subID int64) ([]models.Mount, error) {
	rows, err := db.Query(`
SELECT id, submission_id, role, type, potential, color, name, energy_value, energy_max,
    speed_val, speed_lim, speed_max,
    accel_val, accel_lim, accel_max,
    altitude_val, altitude_lim, altitude_max,
    energy_stat_val, energy_stat_lim, energy_stat_max,
    handling_val, handling_lim, handling_max,
    toughness_val, toughness_lim, toughness_max,
    boost_val, boost_lim, boost_max,
    training_val, training_lim, training_max
FROM mounts WHERE submission_id = ? ORDER BY role`, subID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMounts(rows)
}

func scanMounts(rows *sql.Rows) ([]models.Mount, error) {
	var mounts []models.Mount
	for rows.Next() {
		var m models.Mount
		err := rows.Scan(
			&m.ID, &m.SubmissionID, &m.Role, &m.Type, &m.Potential, &m.Color, &m.Name, &m.EnergyValue, &m.EnergyMax,
			&m.SpeedVal, &m.SpeedLim, &m.SpeedMax,
			&m.AccelVal, &m.AccelLim, &m.AccelMax,
			&m.AltitudeVal, &m.AltitudeLim, &m.AltitudeMax,
			&m.EnergyStatVal, &m.EnergyStatLim, &m.EnergyStatMax,
			&m.HandlingVal, &m.HandlingLim, &m.HandlingMax,
			&m.ToughnessVal, &m.ToughnessLim, &m.ToughnessMax,
			&m.BoostVal, &m.BoostLim, &m.BoostMax,
			&m.TrainingVal, &m.TrainingLim, &m.TrainingMax,
		)
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, m)
	}
	return mounts, rows.Err()
}

func (db *DB) ListSubmissions(limit, offset int) ([]models.Submission, int, error) {
	var total int
	if err := db.QueryRow(`SELECT COUNT(*) FROM submissions`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := db.Query(`
SELECT s.id, s.user_id, u.username, s.notes, s.status, s.created_at
FROM submissions s JOIN users u ON s.user_id = u.id
ORDER BY s.created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subs []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Notes, &s.Status, &s.CreatedAt); err != nil {
			return nil, 0, err
		}
		subs = append(subs, s)
	}
	return subs, total, rows.Err()
}

// ListPendingByUser returns all pending submissions for a specific user.
func (db *DB) ListPendingByUser(userID int64) ([]models.Submission, error) {
	rows, err := db.Query(`
SELECT s.id, s.user_id, u.username, s.notes, s.status, s.created_at
FROM submissions s JOIN users u ON s.user_id = u.id
WHERE s.user_id = ? AND s.status = 'pending'
ORDER BY s.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Notes, &s.Status, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func (db *DB) DeleteSubmission(id int64) error {
	_, err := db.Exec(`DELETE FROM submissions WHERE id = ?`, id)
	return err
}
