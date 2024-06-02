package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/controllers"
	"github.com/cpgillem/garden-planner/ui"
)

type FeatureEditor struct {
	// Data
	controller *controllers.PlantController

	// Window
	window fyne.Window

	// Containers
	mainContainer *fyne.Container
	plantForm     *fyne.Container

	// Widgets
	toolbar   *widget.Toolbar
	plantList *widget.List

	// Form Widgets
	idLabel           *widget.Label
	nameLabel         *widget.Label
	nameEntry         *widget.Entry
	interactionLabel  *widget.Label
	interactionEditor *ui.InteractionEditor
}
