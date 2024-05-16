package main

func main() {
	// Load basic data the program needs.
	gardenData := NewGardenData()
	// Setup instance of UI.
	gardenPlanner := NewGardenPlanner(gardenData)

	// Display UI.
	gardenPlanner.Start()
}
