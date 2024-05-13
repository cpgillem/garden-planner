package models

// Instance of a plant in a garden plan. Plants have all characteristics of landscaping features.
type Plant struct {
	Feature
	Type *PlantType
}
