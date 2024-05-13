package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	Toolbar fyne.Widget

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

func (p *GardenPlanner) SetPlan(plan *Plan) {
	p.currentPlan = plan

	// Create a feature widget.
	featureWidget := ui.NewFeatureWidget(&p.currentPlan.RootFeature)

	p.GardenContainer.Add(featureWidget)
}

func NewGardenPlanner() *GardenPlanner {
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Garden Planner")
	sidebar := container.NewVBox()
	gardenContainer := container.New(&ui.GardenLayout{})
	toolbar := widget.NewToolbar()
	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, gardenContainer)

	gardenPlanner := GardenPlanner{
		currentPlan:     nil,
		App:             mainApp,
		Window:          mainWindow,
		MainContainer:   mainContainer,
		Sidebar:         sidebar,
		Toolbar:         toolbar,
		GardenContainer: container.New(&ui.GardenLayout{}),
	}

	gardenPlanner.Window.SetContent(mainContainer)

	return &gardenPlanner
}
