package database

import (
	"wynnbreeder/models"
)

func (db *DB) GetStatInheritance() ([]models.StatRow, error) {
	stats := []struct {
		name   string
		valCol string
		limCol string
		maxCol string
	}{
		{"Speed", "speed_val", "speed_lim", "speed_max"},
		{"Acceleration", "accel_val", "accel_lim", "accel_max"},
		{"Altitude", "altitude_val", "altitude_lim", "altitude_max"},
		{"Energy", "energy_stat_val", "energy_stat_lim", "energy_stat_max"},
		{"Handling", "handling_val", "handling_lim", "handling_max"},
		{"Toughness", "toughness_val", "toughness_lim", "toughness_max"},
		{"Boost", "boost_val", "boost_lim", "boost_max"},
		{"Training", "training_val", "training_lim", "training_max"},
	}

	var rows []models.StatRow
	for _, s := range stats {
		row, err := db.getStatRow(s.name, s.valCol, s.limCol, s.maxCol)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (db *DB) getStatRow(statName, valCol, limCol, maxCol string) (models.StatRow, error) {
	query := `
SELECT
    AVG(pa.` + valCol + `) as avg_a_val,
    AVG(pa.` + limCol + `) as avg_a_lim,
    AVG(pa.` + maxCol + `) as avg_a_max,
    AVG(pb.` + valCol + `) as avg_b_val,
    AVG(pb.` + limCol + `) as avg_b_lim,
    AVG(pb.` + maxCol + `) as avg_b_max,
    AVG((pa.` + valCol + ` + pb.` + valCol + `) / 2.0) as avg_parent_avg_val,
    AVG((pa.` + limCol + ` + pb.` + limCol + `) / 2.0) as avg_parent_avg_lim,
    AVG((pa.` + maxCol + ` + pb.` + maxCol + `) / 2.0) as avg_parent_avg_max,
    AVG(off.` + valCol + `) as avg_off_val,
    AVG(off.` + limCol + `) as avg_off_lim,
    AVG(off.` + maxCol + `) as avg_off_max,
    AVG(off.` + valCol + ` - (pa.` + valCol + ` + pb.` + valCol + `) / 2.0) as delta_val,
    AVG(off.` + limCol + ` - (pa.` + limCol + ` + pb.` + limCol + `) / 2.0) as delta_lim,
    AVG(off.` + maxCol + ` - (pa.` + maxCol + ` + pb.` + maxCol + `) / 2.0) as delta_max,
    COUNT(*) as cnt
FROM submissions s
JOIN mounts pa ON pa.submission_id = s.id AND pa.role = 'parent_a'
JOIN mounts pb ON pb.submission_id = s.id AND pb.role = 'parent_b'
JOIN mounts off ON off.submission_id = s.id AND off.role = 'offspring'`

	var row models.StatRow
	row.StatName = statName
	err := db.QueryRow(query).Scan(
		&row.AvgParentAVal, &row.AvgParentALim, &row.AvgParentAMax,
		&row.AvgParentBVal, &row.AvgParentBLim, &row.AvgParentBMax,
		&row.AvgParentAvgVal, &row.AvgParentAvgLim, &row.AvgParentAvgMax,
		&row.AvgOffspringVal, &row.AvgOffspringLim, &row.AvgOffspringMax,
		&row.DeltaVal, &row.DeltaLim, &row.DeltaMax,
		&row.Count,
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
