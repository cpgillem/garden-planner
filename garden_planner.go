package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/models"
	"github.com/cpgillem/garden-planner/ui"
)

// Represents the state of the application.
type GardenPlanner struct {
	App    fyne.App
	Window fyne.Window

	// Containers
	MainContainer   *fyne.Container
	Sidebar         *fyne.Container
	GardenContainer *fyne.Container

	// Permanent Widgets
	Toolbar   *widget.Toolbar
	StatusBar *widget.Label

	// Data
	currentPlan *Plan
}

// Represents a garden plan loaded from a JSON file.
type Plan struct {
	Name        string         `json:"name"`
	RootFeature models.Feature `json:"root_feature"`
}

func (p *GardenPlanner) Start() {
	p.Window.Show()
	p.App.Run()
}

func (instance *GardenPlanner) OpenPlan(plan *Plan) {
	instance.currentPlan = plan

	// Create a feature widget.
	featureWidget := ui.NewFeatureWidget(&instance.currentPlan.RootFeature)

	instance.GardenContainer.Add(featureWidget)

	// Setup Sidebar
	featureList := widget.NewList(
		// Length
		func() int {
			if instance.currentPlan == nil {
				return 0
			}
			return int(len(instance.currentPlan.RootFeature.Features))
		},
		// Create
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		// Update
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(instance.currentPlan.RootFeature.Features[id].Name)
		},
	)

	instance.Sidebar.Add(featureList)
}

// Cleans up the UI elements depending on a current plan.
func (instance *GardenPlanner) ClosePlan() {
	instance.GardenContainer.RemoveAll()
	instance.Sidebar.RemoveAll()
	instance.currentPlan = nil
}

func NewGardenPlanner() *GardenPlanner {
	// Setup UI elements
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Garden Planner")
	sidebar := container.NewVBox()
	gardenContainer := container.New(&ui.GardenLayout{})
	toolbar := widget.NewToolbar()
	statusBar := widget.NewLabel("")
	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, gardenContainer)

	// Setup Toolbar
	toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {

	}))
	toolbar.Append(widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
		fmt.Println("Clicked")
	}))
	toolbar.Append(widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {}))
	toolbar.Append(widget.NewToolbarSeparator())

	mainWindow.SetContent(mainContainer)

	// Create new app instance
	gardenPlanner := GardenPlanner{
		currentPlan:     nil,
		App:             mainApp,
		Window:          mainWindow,
		MainContainer:   mainContainer,
		Sidebar:         sidebar,
		Toolbar:         toolbar,
		StatusBar:       statusBar,
		GardenContainer: container.New(&ui.GardenLayout{}),
	}

	return &gardenPlanner
}
