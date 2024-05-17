package main

import "github.com/cpgillem/garden-planner/models"

func main() {
	// Load basic data the program needs.
	gardenData := NewGardenData()
	// Setup instance of UI.
	gardenPlanner := NewGardenPlanner(gardenData)

	// Load test plan for now.
	testPlan, _ := ReadObjectFromFile[models.Plan]("test_data/layout1.json")
	gardenPlanner.OpenPlan(testPlan)

	// Display UI.
	gardenPlanner.Start()
}
