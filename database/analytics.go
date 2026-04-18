package database

import (
	"wynnmounts/models"
)

func (db *DB) GetStatInheritance() ([]models.StatRow, error) {
	stats := []struct {
		name   string
		aCol   string
		bCol   string
		offCol string
	}{
		{"Speed", "speed_val", "speed_val", "speed_val"},
		{"Acceleration", "accel_val", "accel_val", "accel_val"},
		{"Altitude", "altitude_val", "altitude_val", "altitude_val"},
		{"Energy", "energy_stat_val", "energy_stat_val", "energy_stat_val"},
		{"Handling", "handling_val", "handling_val", "handling_val"},
		{"Toughness", "toughness_val", "toughness_val", "toughness_val"},
		{"Boost", "boost_val", "boost_val", "boost_val"},
		{"Training", "training_val", "training_val", "training_val"},
	}

	var rows []models.StatRow
	for _, s := range stats {
		row, err := db.getStatRow(s.name, s.aCol)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (db *DB) getStatRow(statName, col string) (models.StatRow, error) {
	query := `
SELECT
    AVG(pa.` + col + `) as avg_a,
    AVG(pb.` + col + `) as avg_b,
    AVG((pa.` + col + ` + pb.` + col + `) / 2.0) as avg_parent_avg,
    AVG(off.` + col + `) as avg_off,
    AVG(off.` + col + ` - (pa.` + col + ` + pb.` + col + `) / 2.0) as delta,
    COUNT(*) as cnt
FROM submissions s
JOIN mounts pa ON pa.submission_id = s.id AND pa.role = 'parent_a'
JOIN mounts pb ON pb.submission_id = s.id AND pb.role = 'parent_b'
JOIN mounts off ON off.submission_id = s.id AND off.role = 'offspring'`

	var row models.StatRow
	row.StatName = statName
	err := db.QueryRow(query).Scan(
		&row.AvgParentAVal, &row.AvgParentBVal, &row.AvgParentAvg,
		&row.AvgOffspringVal, &row.Delta, &row.Count,
	)
	return row, err
}

func (db *DB) GetColorInheritance() ([]models.ColorInheritance, error) {
	rows, err := db.Query(`
SELECT pa.color, pb.color, off.color, COUNT(*) as cnt
FROM submissions s
JOIN mounts pa ON pa.submission_id = s.id AND pa.role = 'parent_a'
JOIN mounts pb ON pb.submission_id = s.id AND pb.role = 'parent_b'
JOIN mounts off ON off.submission_id = s.id AND off.role = 'offspring'
GROUP BY pa.color, pb.color, off.color
ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ColorInheritance
	for rows.Next() {
		var ci models.ColorInheritance
		if err := rows.Scan(&ci.ParentAColor, &ci.ParentBColor, &ci.OffspringColor, &ci.Count); err != nil {
			return nil, err
		}
		results = append(results, ci)
	}
	return results, rows.Err()
}

func (db *DB) GetPotentialStats() (*models.PotentialStats, error) {
	ps := &models.PotentialStats{}
	err := db.QueryRow(`
SELECT AVG(pa.potential), AVG(pb.potential), AVG(off.potential), COUNT(*)
FROM submissions s
JOIN mounts pa ON pa.submission_id = s.id AND pa.role = 'parent_a'
JOIN mounts pb ON pb.submission_id = s.id AND pb.role = 'parent_b'
JOIN mounts off ON off.submission_id = s.id AND off.role = 'offspring'`,
	).Scan(&ps.AvgParentAPotential, &ps.AvgParentBPotential, &ps.AvgOffspringPotential, &ps.Count)
	return ps, err
}
