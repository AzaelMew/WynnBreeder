package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	IsAdmin      bool      `json:"is_admin"`
	IsSuperAdmin bool      `json:"is_superadmin"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Submission struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Notes     string    `json:"notes"`
	Status    string    `json:"status"` // "pending" or "complete"
	CreatedAt time.Time `json:"created_at"`
	Mounts    []Mount   `json:"mounts,omitempty"`
}

type MountRole string

const (
	RoleParentA   MountRole = "parent_a"
	RoleParentB   MountRole = "parent_b"
	RoleOffspring MountRole = "offspring"
)

type StatValues struct {
	Value    int `json:"value"`
	Limit    int `json:"limit"`
	MaxValue int `json:"maxValue"`
}

type EnergyValues struct {
	Value    int `json:"value"`
	MaxValue int `json:"maxValue"`
}

// MountJSON matches imported JSON structure
type MountJSON struct {
	Type      string       `json:"type"`
	Potential int          `json:"potential"`
	Color     string       `json:"color"`
	Name      string       `json:"name"`
	Energy    EnergyValues `json:"energy"`
	Stats     struct {
		Speed        StatValues `json:"speed"`
		Acceleration StatValues `json:"acceleration"`
		Altitude     StatValues `json:"altitude"`
		Energy       StatValues `json:"energy"`
		Handling     StatValues `json:"handling"`
		Toughness    StatValues `json:"toughness"`
		Boost        StatValues `json:"boost"`
		Training     StatValues `json:"training"`
	} `json:"stats"`
}

type Mount struct {
	ID           int64     `json:"id"`
	SubmissionID int64     `json:"submission_id"`
	Role         MountRole `json:"role"`
	Type         string    `json:"type"`
	Potential    int       `json:"potential"`
	Color        string    `json:"color"`
	Name         string    `json:"name"`
	EnergyValue  int       `json:"energy_value"`
	EnergyMax    int       `json:"energy_max"`
	SpeedVal     int       `json:"speed_val"`
	SpeedLim     int       `json:"speed_lim"`
	SpeedMax     int       `json:"speed_max"`
	AccelVal     int       `json:"accel_val"`
	AccelLim     int       `json:"accel_lim"`
	AccelMax     int       `json:"accel_max"`
	AltitudeVal  int       `json:"altitude_val"`
	AltitudeLim  int       `json:"altitude_lim"`
	AltitudeMax  int       `json:"altitude_max"`
	EnergyStatVal int      `json:"energy_stat_val"`
	EnergyStatLim int      `json:"energy_stat_lim"`
	EnergyStatMax int      `json:"energy_stat_max"`
	HandlingVal  int       `json:"handling_val"`
	HandlingLim  int       `json:"handling_lim"`
	HandlingMax  int       `json:"handling_max"`
	ToughnessVal int       `json:"toughness_val"`
	ToughnessLim int       `json:"toughness_lim"`
	ToughnessMax int       `json:"toughness_max"`
	BoostVal     int       `json:"boost_val"`
	BoostLim     int       `json:"boost_lim"`
	BoostMax     int       `json:"boost_max"`
	TrainingVal  int       `json:"training_val"`
	TrainingLim  int       `json:"training_lim"`
	TrainingMax  int       `json:"training_max"`
}

func MountFromJSON(mj MountJSON, submissionID int64, role MountRole) Mount {
	return Mount{
		SubmissionID:  submissionID,
		Role:          role,
		Type:          mj.Type,
		Potential:     mj.Potential,
		Color:         mj.Color,
		Name:          mj.Name,
		EnergyValue:   mj.Energy.Value,
		EnergyMax:     mj.Energy.MaxValue,
		SpeedVal:      mj.Stats.Speed.Value,
		SpeedLim:      mj.Stats.Speed.Limit,
		SpeedMax:      mj.Stats.Speed.MaxValue,
		AccelVal:      mj.Stats.Acceleration.Value,
		AccelLim:      mj.Stats.Acceleration.Limit,
		AccelMax:      mj.Stats.Acceleration.MaxValue,
		AltitudeVal:   mj.Stats.Altitude.Value,
		AltitudeLim:   mj.Stats.Altitude.Limit,
		AltitudeMax:   mj.Stats.Altitude.MaxValue,
		EnergyStatVal: mj.Stats.Energy.Value,
		EnergyStatLim: mj.Stats.Energy.Limit,
		EnergyStatMax: mj.Stats.Energy.MaxValue,
		HandlingVal:   mj.Stats.Handling.Value,
		HandlingLim:   mj.Stats.Handling.Limit,
		HandlingMax:   mj.Stats.Handling.MaxValue,
		ToughnessVal:  mj.Stats.Toughness.Value,
		ToughnessLim:  mj.Stats.Toughness.Limit,
		ToughnessMax:  mj.Stats.Toughness.MaxValue,
		BoostVal:      mj.Stats.Boost.Value,
		BoostLim:      mj.Stats.Boost.Limit,
		BoostMax:      mj.Stats.Boost.MaxValue,
		TrainingVal:   mj.Stats.Training.Value,
		TrainingLim:   mj.Stats.Training.Limit,
		TrainingMax:   mj.Stats.Training.MaxValue,
	}
}

type SubmitRequest struct {
	ParentA   MountJSON  `json:"parent_a"`
	ParentB   MountJSON  `json:"parent_b"`
	Offspring *MountJSON `json:"offspring"` // nil = save as pending (breed in progress)
	Notes     string     `json:"notes"`
}

// StatRow used in analytics tables
type StatRow struct {
	StatName       string  `json:"stat_name"`
	AvgParentAVal  float64 `json:"avg_parent_a_val"`
	AvgParentBVal  float64 `json:"avg_parent_b_val"`
	AvgParentAvg   float64 `json:"avg_parent_avg"`
	AvgOffspringVal float64 `json:"avg_offspring_val"`
	Delta          float64 `json:"delta"`
	Count          int     `json:"count"`
}

type ColorInheritance struct {
	ParentAColor    string `json:"parent_a_color"`
	ParentBColor    string `json:"parent_b_color"`
	OffspringColor  string `json:"offspring_color"`
	Count           int    `json:"count"`
}

type PotentialStats struct {
	AvgParentAPotential  float64 `json:"avg_parent_a_potential"`
	AvgParentBPotential  float64 `json:"avg_parent_b_potential"`
	AvgOffspringPotential float64 `json:"avg_offspring_potential"`
	Count                int     `json:"count"`
}
