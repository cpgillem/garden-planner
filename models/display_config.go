package models

import "github.com/bcicen/go-units"

type DisplayConfig struct {
	Scale       float32
	BaseUnit    units.Unit
	GridSpacing units.Value
}

func NewDisplayConfig() DisplayConfig {
	return DisplayConfig{
		Scale:       2,
		BaseUnit:    units.Inch,
		GridSpacing: units.NewValue(12, units.Inch),
	}
}
