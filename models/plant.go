package models

// Instance of a plant/plant row in a garden plan. Plants have all characteristics of landscaping features.
type Plant struct {
	Feature `json:"feature"`
	Type    *PlantType `json:"-"`

	// ID of the plant type when reading JSON.
	plantTypeID int `json:"plant_type_id"`
}
