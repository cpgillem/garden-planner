package models

import "github.com/cpgillem/garden-planner/geometry"

// Represents a type of feature, e.g. plant row, plant pot, fountain, etc.
type FeatureTemplate struct {
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Properties  []string     `json:"properties"`
	Box         geometry.Box `json:"box"`
}
