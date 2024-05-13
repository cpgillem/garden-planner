package models

// Plant species, cultivar, or whatever else (author is not a botanist).
// Examples: potato, cabbage, broccoli. Does not cover different variants, such as "better boy" tomatoes.
// In the future, this will be divided into variants where only some properties
// apply to all plants of the type.
type PlantType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Minimum spacing between rows in inches
	RowSpacing uint16 `json:"row_spacing"`
	// Plant spread, in inches, as a radius
	Spread       uint16             `json:"spread"`
	Interactions []PlantInteraction `json:"interactions"`
}
