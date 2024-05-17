package models

import "github.com/cpgillem/garden-planner/geometry"

type Plan struct {
	Name     string                          `json:"name"`
	Box      geometry.AxisAlignedBoundingBox `json:"box"`
	Features []Feature                       `json:"features"`
}

func NewPlan() *Plan {
	return &Plan{
		Name:     "",
		Box:      geometry.NewBox(),
		Features: []Feature{},
	}
}
