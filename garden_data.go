package main

import (
	"fmt"

	"github.com/cpgillem/garden-planner/models"
)

type GardenData struct {
	Properties       map[string]models.Property
	FeatureTemplates map[string]models.FeatureTemplate
}

// Loads garden data from the files in the /data directory.
func NewGardenData() *GardenData {
	// Start with empty data.
	gardenData := GardenData{
		Properties:       map[string]models.Property{},
		FeatureTemplates: map[string]models.FeatureTemplate{},
	}

	// Load properties of any landscape feature.
	properties, err := ReadObjectFromFile[[]models.Property]("data/properties.json")
	if err != nil {
		fmt.Println("Could not load properties.")
	}

	// Map properties for easy retrieval.
	for _, p := range *properties {
		gardenData.Properties[p.Name] = p
	}

	// Load templates for landscaping features.
	featureTemplates, err := ReadObjectFromFile[[]models.FeatureTemplate]("data/feature_templates.json")
	if err != nil {
		fmt.Println("Could not load feature templates.")
	}

	// Map feature templates to names.
	for _, ft := range *featureTemplates {
		gardenData.FeatureTemplates[ft.Name] = ft
	}

	return &gardenData
}
