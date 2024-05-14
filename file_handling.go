package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cpgillem/garden-planner/models"
)

// Open a file as a garden plan.
// TODO: Genericize this to open any JSON file to avoid
// repeating file opening idioms.
func LoadPlan(path string) (*models.Plan, error) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("Could not load plan from " + path + "\n" + err.Error())
		return nil, err
	}

	defer f.Close()

	content, err := io.ReadAll(f)

	if err != nil {
		fmt.Println("Could not read file: " + path + "\n" + err.Error())
		return nil, err
	}

	var plan models.Plan

	err = json.Unmarshal(content, &plan)

	if err != nil {
		fmt.Println("Could not parse JSON.\n" + err.Error())
	}

	return &plan, nil
}
