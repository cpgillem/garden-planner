package models

import "github.com/cpgillem/garden-planner/geometry"

type Plan struct {
	Name     string                 `json:"name"`
	Box      geometry.Box           `json:"box"`
	Features map[FeatureID]*Feature `json:"features"`
}

func NewPlan() *Plan {
	return &Plan{
		Name:     "",
		Box:      geometry.NewBoxZero(),
		Features: map[FeatureID]*Feature{},
	}
}
