package models

import "github.com/bcicen/go-units"

type DisplayConfig struct {
	Scale       float32
	GridSpacing units.Value
}

func NewDisplayConfig() DisplayConfig {
	return DisplayConfig{
		Scale:       2,
		GridSpacing: units.NewValue(12, units.Inch),
	}
}
