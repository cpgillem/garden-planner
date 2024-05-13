package models

// Plant species, cultivar, or whatever else (author is not a botanist).
// Examples: potato, cabbage, broccoli. Does not cover different variants, such as "better boy" tomatoes.
type PlantType struct {
	id            uint
	DaysToHarvest uint16
	// Minimum spacing between rows in cm
	RowSpacing uint16
	// Plant spread, in cm
	Spread       uint16
	Interactions []PlantInteraction
}
