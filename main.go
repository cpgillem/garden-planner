package main

import (
	"fmt"

	"github.com/cpgillem/garden-planner/models"
)

func main() {
	// Setup instance of UI.
	gardenPlanner := NewGardenPlanner()

	// Load test plan for now.
	testPlan, err := ReadObjectFromFile[models.Plan]("test_data/layout1.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	gardenPlanner.OpenPlan(testPlan)

	// Display UI.
	gardenPlanner.Start()
}
