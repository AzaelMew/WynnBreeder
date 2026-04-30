package models

import "fmt"

// FoodMaterial mirrors the WynnMounts MATERIALS array:
// [tier, name, speed, acceleration, altitude, energy, handling, toughness, boost, training]
type FoodMaterial struct {
	Tier          int
	Name          string
	Speed         int
	Acceleration  int
	Altitude      int
	Energy        int
	Handling      int
	Toughness     int
	Boost         int
	Training      int
}

// PendingFood is a user-selected food item with quantity.
type PendingFood struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// Materials ported from WynnMounts data.js.
var Materials = []FoodMaterial{
	{1, "Copper Ingot", 0, 0, 0, 4, 0, 8, 0, 0},
	{1, "Copper Gem", 4, 0, 0, 2, 0, 0, 0, 6},
	{1, "Oak Plank", 2, 6, 0, 0, 0, 4, 0, 0},
	{1, "Oak Paper", 0, 0, 8, 0, 0, 0, 4, 0},
	{1, "Wheat String", 0, 2, 0, 0, 4, 0, 6, 0},
	{1, "Wheat Grains", 8, 0, 4, 0, 0, 0, 0, 0},
	{1, "Gudgeon Oil", 0, 0, 2, 0, 6, 0, 0, 4},
	{1, "Gudgeon Meat", 0, 4, 0, 8, 0, 0, 0, 0},
	{10, "Granite Ingot", 0, 0, 0, 5, 0, 10, 0, 0},
	{10, "Granite Gem", 5, 0, 0, 2, 0, 0, 0, 8},
	{10, "Birch Plank", 2, 8, 0, 0, 0, 5, 0, 0},
	{10, "Birch Paper", 0, 0, 10, 0, 0, 0, 5, 0},
	{10, "Barley String", 0, 2, 0, 0, 5, 0, 8, 0},
	{10, "Barley Grains", 10, 0, 5, 0, 0, 0, 0, 0},
	{10, "Trout Oil", 0, 0, 2, 0, 8, 0, 0, 5},
	{10, "Trout Meat", 0, 5, 0, 10, 0, 0, 0, 0},
	{20, "Gold Ingot", 0, 0, 0, 5, 0, 12, 0, 0},
	{20, "Gold Gem", 6, 0, 0, 3, 0, 0, 0, 9},
	{20, "Willow Plank", 3, 9, 0, 0, 0, 6, 0, 0},
	{20, "Willow Paper", 0, 0, 12, 0, 0, 0, 5, 0},
	{20, "Oat String", 0, 3, 0, 0, 6, 0, 9, 0},
	{20, "Oat Grains", 12, 0, 5, 0, 0, 0, 0, 0},
	{20, "Salmon Oil", 0, 0, 3, 0, 9, 0, 0, 6},
	{20, "Salmon Meat", 0, 5, 0, 12, 0, 0, 0, 0},
	{30, "Sandstone Ingot", 0, 0, 0, 6, 0, 14, 0, 0},
	{30, "Sandstone Gem", 6, 0, 0, 3, 0, 0, 0, 11},
	{30, "Acacia Plank", 3, 11, 0, 0, 0, 6, 0, 0},
	{30, "Acacia Paper", 0, 0, 14, 0, 0, 0, 6, 0},
	{30, "Malt String", 0, 3, 0, 0, 6, 0, 11, 0},
	{30, "Malt Grains", 14, 0, 6, 0, 0, 0, 0, 0},
	{30, "Carp Oil", 0, 0, 3, 0, 11, 0, 0, 6},
	{30, "Carp Meat", 0, 6, 0, 14, 0, 0, 0, 0},
	{40, "Iron Ingot", 0, 0, 0, 6, 0, 16, 0, 0},
	{40, "Iron Gem", 7, 0, 0, 3, 0, 0, 0, 12},
	{40, "Spruce Plank", 3, 12, 0, 0, 0, 7, 0, 0},
	{40, "Spruce Paper", 0, 0, 16, 0, 0, 0, 6, 0},
	{40, "Hops String", 0, 3, 0, 0, 7, 0, 12, 0},
	{40, "Hops Grains", 16, 0, 6, 0, 0, 0, 0, 0},
	{40, "Icefish Oil", 0, 0, 3, 0, 12, 0, 0, 7},
	{40, "Icefish Meat", 0, 6, 0, 16, 0, 0, 0, 0},
	{50, "Silver Ingot", 0, 0, 0, 7, 0, 18, 0, 0},
	{50, "Silver Gem", 8, 0, 0, 4, 0, 0, 0, 14},
	{50, "Jungle Plank", 4, 14, 0, 0, 0, 8, 0, 0},
	{50, "Jungle Paper", 0, 0, 18, 0, 0, 0, 7, 0},
	{50, "Rye String", 0, 4, 0, 0, 8, 0, 14, 0},
	{50, "Rye Grains", 18, 0, 7, 0, 0, 0, 0, 0},
	{50, "Piranha Oil", 0, 0, 4, 0, 14, 0, 0, 8},
	{50, "Piranha Meat", 0, 7, 0, 18, 0, 0, 0, 0},
	{60, "Cobalt Ingot", 0, 0, 0, 8, 0, 20, 0, 0},
	{60, "Cobalt Gem", 9, 0, 0, 4, 0, 0, 0, 15},
	{60, "Dark Plank", 4, 15, 0, 0, 0, 9, 0, 0},
	{60, "Dark Paper", 0, 0, 20, 0, 0, 0, 8, 0},
	{60, "Millet String", 0, 4, 0, 0, 9, 0, 15, 0},
	{60, "Millet Grains", 20, 0, 8, 0, 0, 0, 0, 0},
	{60, "Koi Oil", 0, 0, 4, 0, 15, 0, 0, 9},
	{60, "Koi Meat", 0, 8, 0, 20, 0, 0, 0, 0},
	{70, "Kanderstone Ingot", 0, 0, 0, 8, 0, 22, 0, 0},
	{70, "Kanderstone Gem", 10, 0, 0, 4, 0, 0, 0, 17},
	{70, "Light Plank", 4, 17, 0, 0, 0, 10, 0, 0},
	{70, "Light Paper", 0, 0, 22, 0, 0, 0, 8, 0},
	{70, "Decay String", 0, 4, 0, 0, 10, 0, 17, 0},
	{70, "Decay Grains", 22, 0, 8, 0, 0, 0, 0, 0},
	{70, "Gylia Oil", 0, 0, 4, 0, 17, 0, 0, 10},
	{70, "Gylia Meat", 0, 8, 0, 22, 0, 0, 0, 0},
	{80, "Diamond Ingot", 0, 0, 0, 9, 0, 24, 0, 0},
	{80, "Diamond Gem", 10, 0, 0, 4, 0, 0, 0, 18},
	{80, "Pine Plank", 4, 18, 0, 0, 0, 10, 0, 0},
	{80, "Pine Paper", 0, 0, 24, 0, 0, 0, 9, 0},
	{80, "Rice String", 0, 4, 0, 0, 10, 0, 18, 0},
	{80, "Rice Grains", 24, 0, 9, 0, 0, 0, 0, 0},
	{80, "Bass Oil", 0, 0, 4, 0, 18, 0, 0, 10},
	{80, "Bass Meat", 0, 9, 0, 24, 0, 0, 0, 0},
	{90, "Molten Ingot", 0, 0, 0, 9, 0, 26, 0, 0},
	{90, "Molten Gem", 11, 0, 0, 5, 0, 0, 0, 20},
	{90, "Avo Plank", 5, 20, 0, 0, 0, 11, 0, 0},
	{90, "Avo Paper", 0, 0, 26, 0, 0, 0, 9, 0},
	{90, "Sorghum String", 0, 5, 0, 0, 11, 0, 20, 0},
	{90, "Sorghum Grains", 26, 0, 9, 0, 0, 0, 0, 0},
	{90, "Molten Oil", 0, 0, 5, 0, 20, 0, 0, 11},
	{90, "Molten Meat", 0, 9, 0, 26, 0, 0, 0, 0},
	{100, "Voidstone Ingot", 0, 0, 0, 10, 0, 28, 0, 0},
	{100, "Voidstone Gem", 12, 0, 0, 5, 0, 0, 0, 21},
	{100, "Sky Plank", 5, 21, 0, 0, 0, 12, 0, 0},
	{100, "Sky Paper", 0, 0, 28, 0, 0, 0, 10, 0},
	{100, "Hemp String", 0, 5, 0, 0, 12, 0, 21, 0},
	{100, "Hemp Grains", 28, 0, 10, 0, 0, 0, 0, 0},
	{100, "Starfish Oil", 0, 0, 5, 0, 21, 0, 0, 12},
	{100, "Starfish Meat", 0, 10, 0, 28, 0, 0, 0, 0},
	{105, "Dernic Ingot", 0, 0, 0, 10, 0, 29, 0, 0},
	{105, "Dernic Gem", 12, 0, 0, 5, 0, 0, 0, 22},
	{105, "Dernic Plank", 5, 22, 0, 0, 0, 12, 0, 0},
	{105, "Dernic Paper", 0, 0, 29, 0, 0, 0, 10, 0},
	{105, "Dernic String", 0, 5, 0, 0, 12, 0, 22, 0},
	{105, "Dernic Grains", 29, 0, 10, 0, 0, 0, 0, 0},
	{105, "Dernic Oil", 0, 0, 5, 0, 22, 0, 0, 12},
	{105, "Dernic Meat", 0, 10, 0, 29, 0, 0, 0, 0},
	{110, "Titanium Ingot", 0, 0, 0, 11, 0, 30, 0, 0},
	{110, "Titanium Gem", 13, 0, 0, 5, 0, 0, 0, 23},
	{110, "Maple Plank", 5, 23, 0, 0, 0, 13, 0, 0},
	{110, "Maple Paper", 0, 0, 30, 0, 0, 0, 11, 0},
	{110, "Jute String", 0, 5, 0, 0, 13, 0, 23, 0},
	{110, "Jute Grains", 30, 0, 11, 0, 0, 0, 0, 0},
	{110, "Sturgeon Oil", 0, 0, 5, 0, 23, 0, 0, 13},
	{110, "Sturgeon Meat", 0, 11, 0, 30, 0, 0, 0, 0},
	{115, "Cinnabar Ingot", 0, 0, 0, 11, 0, 31, 0, 0},
	{115, "Cinnabar Gem", 13, 0, 0, 5, 0, 0, 0, 23},
	{115, "Redwood Plank", 5, 23, 0, 0, 0, 13, 0, 0},
	{115, "Redwood Paper", 0, 0, 31, 0, 0, 0, 11, 0},
	{115, "Heather String", 0, 5, 0, 0, 13, 0, 23, 0},
	{115, "Heather Grains", 31, 0, 11, 0, 0, 0, 0, 0},
	{115, "Mahseer Oil", 0, 0, 5, 0, 23, 0, 0, 13},
	{115, "Mahseer Meat", 0, 11, 0, 31, 0, 0, 0, 0},
}

// FoodByName maps food name to material for fast lookup.
var FoodByName = func() map[string]FoodMaterial {
	m := make(map[string]FoodMaterial, len(Materials))
	for _, mat := range Materials {
		m[mat.Name] = mat
	}
	return m
}()

// ApplyFoods applies pending food boosts to a mount's limit values.
// Limits are capped at maxValue. Values are left unchanged.
func ApplyFoods(mj *MountJSON, foods []PendingFood) error {
	// Sum boosts per stat.
	var speed, accel, alt, energy, handling, toughness, boost, training int

	for _, f := range foods {
		if f.Quantity <= 0 {
			continue
		}
		mat, ok := FoodByName[f.Name]
		if !ok {
			return fmt.Errorf("unknown food: %s", f.Name)
		}
		speed += mat.Speed * f.Quantity
		accel += mat.Acceleration * f.Quantity
		alt += mat.Altitude * f.Quantity
		energy += mat.Energy * f.Quantity
		handling += mat.Handling * f.Quantity
		toughness += mat.Toughness * f.Quantity
		boost += mat.Boost * f.Quantity
		training += mat.Training * f.Quantity
	}

	// Apply to limits, capping at maxValue.
	addLimit := func(limit, maxVal, boost int) int {
		newLim := limit + boost
		if newLim > maxVal {
			return maxVal
		}
		return newLim
	}

	mj.Stats.Speed.Limit = addLimit(mj.Stats.Speed.Limit, mj.Stats.Speed.MaxValue, speed)
	mj.Stats.Acceleration.Limit = addLimit(mj.Stats.Acceleration.Limit, mj.Stats.Acceleration.MaxValue, accel)
	mj.Stats.Altitude.Limit = addLimit(mj.Stats.Altitude.Limit, mj.Stats.Altitude.MaxValue, alt)
	mj.Stats.Energy.Limit = addLimit(mj.Stats.Energy.Limit, mj.Stats.Energy.MaxValue, energy)
	mj.Stats.Handling.Limit = addLimit(mj.Stats.Handling.Limit, mj.Stats.Handling.MaxValue, handling)
	mj.Stats.Toughness.Limit = addLimit(mj.Stats.Toughness.Limit, mj.Stats.Toughness.MaxValue, toughness)
	mj.Stats.Boost.Limit = addLimit(mj.Stats.Boost.Limit, mj.Stats.Boost.MaxValue, boost)
	mj.Stats.Training.Limit = addLimit(mj.Stats.Training.Limit, mj.Stats.Training.MaxValue, training)

	return nil
}
