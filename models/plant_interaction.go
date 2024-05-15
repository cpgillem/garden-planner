package models

const (
	NEUTRAL      = 0
	BENEFICIAL   = 1
	ANTAGONISTIC = 2
)

// Defines whether the target plant is good or bad to plant next to a subject plant.
type PlantInteraction struct {
	TargetPlantID   int    `json:"target_plant_id"`
	InteractionType uint16 `json:"interaction_type"`
}
