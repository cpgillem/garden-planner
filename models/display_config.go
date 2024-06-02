package models

import "github.com/bcicen/go-units"

type DisplayConfig struct {
	BaseUnit units.Unit
}

func NewDisplayConfig() DisplayConfig {
	return DisplayConfig{
		BaseUnit: units.Inch,
	}
}
