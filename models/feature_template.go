package models

// Represents a type of feature, e.g. plant row, plant pot, fountain, etc.
type FeatureTemplate struct {
	Name       string   `json:"name"`
	Properties []string `json:"properties"`
}
