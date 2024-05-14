package main

import (
	"os"

	"github.com/cpgillem/garden-planner/models"
)

func main() {
	// Setup UI.
	gardenPlanner := NewGardenPlanner()

	// Setup application data.
	// For now, load the test file.
	file, _ := os.Open("test_data/layout1.json")
	plan, _ := ReadObject[models.Plan](file)
	gardenPlanner.OpenPlan(plan)
	file.Close()

	// Display UI.
	gardenPlanner.Start()
}
