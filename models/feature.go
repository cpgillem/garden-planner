package models

import "github.com/cpgillem/garden-planner/geometry"

type FeatureID int

// A landscaping feature, such as a row of plants, planter, garden bed, tree, or obstacle.
// The whole yard, fenced off area, etc. can serve as the root feature.
type Feature struct {
	Box  geometry.AxisAlignedBoundingBox `json:"box"`
	Name string                          `json:"name"`

	// Table of data properties depending on what type of feature this is.
	Properties map[string]any `json:"properties"`
}

// Creates a new plant feature according to the plant template.
func NewPlantFeature(name string, box geometry.AxisAlignedBoundingBox, plantType *PlantType) *Feature {
	plant := Feature{
		Name: name,
		Box:  box,
	}

	plant.Properties["plant_spacing"] = plantType.Spread
	plant.Properties["row_width"] = plantType.Spread

	return &plant
}
