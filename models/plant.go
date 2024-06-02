package models

// Plant species, cultivar, or whatever else (author is not a botanist).
// Examples: potato, cabbage, broccoli. Does not cover different variants, such as "better boy" tomatoes.
// In the future, this will be divided into variants where only some properties
// apply to all plants of the type.
type Plant struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// List of interactions with other plant types for intercropping.
	Interactions []PlantInteraction `json:"interactions"`
}
