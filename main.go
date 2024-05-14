package main

func main() {
	// Setup UI.
	gardenPlanner := NewGardenPlanner()

	// Setup application data.
	// For now, load the test file.
	plan, _ := LoadPlan("test_data/layout1.json")
	gardenPlanner.OpenPlan(plan)

	// Display UI.
	gardenPlanner.Start()
}
