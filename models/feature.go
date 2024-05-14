package models

import "github.com/cpgillem/garden-planner/geometry"

// A landscaping feature, such as a row of plants, planter, garden bed, tree, or obstacle.
// The whole yard, fenced off area, etc. can serve as the root feature.
type Feature struct {
	Box  geometry.AxisAlignedBoundingBox `json:"box"`
	Name string                          `json:"name"`
}
