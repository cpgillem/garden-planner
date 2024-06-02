package models

import "github.com/cpgillem/garden-planner/geometry"

type FeatureID int

// A landscaping feature, such as a row of plants, planter, garden bed, tree, or obstacle.
// The whole yard, fenced off area, etc. can serve as the root feature.
type Feature struct {
	Box  geometry.Box `json:"box"`
	Name string       `json:"name"`

	// Table of data properties depending on what type of feature this is.
	Properties map[string]any `json:"properties"`
}

func NewFeature(propMap map[string]Property, template *FeatureTemplate) Feature {
	f := Feature{
		Name:       template.DisplayName,
		Box:        template.Box.Copy(),
		Properties: map[string]any{},
	}

	// Set default properties.
	for _, propName := range template.Properties {
		prop := propMap[propName]
		f.Properties[prop.Name] = prop.Default
	}

	return f
}
