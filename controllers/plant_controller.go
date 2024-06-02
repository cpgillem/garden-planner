package controllers

import "github.com/cpgillem/garden-planner/models"

// Used to edit the collection of plant data used by the app for all plans.
type PlantController struct {
	plants map[int]models.Plant

	OnPlantAdded   func(models.Plant)
	OnPlantRemoved func(models.Plant)
}

func NewPlantController(plants *[]models.Plant) PlantController {
	c := PlantController{
		plants:         map[int]models.Plant{},
		OnPlantAdded:   func(p models.Plant) {},
		OnPlantRemoved: func(p models.Plant) {},
	}

	// Add initial plants.
	for _, p := range *plants {
		c.plants[p.ID] = p
	}

	return c
}

func (c *PlantController) AddPlant(plant models.Plant) {
	c.plants[plant.ID] = plant
	c.OnPlantAdded(plant)
}

func (c *PlantController) RemovePlant(id int) {
	// TODO
	c.OnPlantRemoved(c.plants[id])
}
