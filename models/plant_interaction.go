package models

const (
	NEUTRAL      uint16 = 0
	BENEFICIAL          = 1
	ANTAGONISTIC        = 2
)

// Defines whether the target plant is good or bad to plant next to a subject plant.
type PlantInteraction struct {
	TargetPlant     *PlantType
	InteractionType uint16
}
