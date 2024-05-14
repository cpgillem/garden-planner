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
	MainContainer *fyne.Container
	Sidebar       *fyne.Container
	Content       *fyne.Container

	// Permanent Widgets
	Toolbar   *widget.Toolbar
	StatusBar *widget.Label

	// Data
	currentPlan *models.Plan
}

func (p *GardenPlanner) Start() {
	p.Window.Show()
	p.App.Run()
}

func (instance *GardenPlanner) OpenPlan(plan *models.Plan) {
	instance.currentPlan = plan

	// Create feature widgets.
	gardenContainer := container.New(ui.NewGardenLayout(plan))
	for _, feature := range plan.Features {
		featureWidget := ui.NewFeatureWidget(&feature)
		gardenContainer.Add(featureWidget)
	}
	instance.Content.Add(gardenContainer)
	gardenContainer.Refresh()

	// Setup Sidebar
	featureList := widget.NewList(
		// Length
		func() int {
			if instance.currentPlan == nil {
				return 0
			}
			return int(len(instance.currentPlan.Features))
		},
		// Create
		func() fyne.CanvasObject {
			// Fix later: why doesn't fyne refresh the width of the sidebar properly?
			label := widget.NewLabel("aaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			label.Truncation = fyne.TextTruncateOff
			return label
		},
		// Update
		func(id widget.ListItemID, o fyne.CanvasObject) {
			obj := o.(*widget.Label)
			obj.SetText(instance.currentPlan.Features[id].Name)
		},
	)

	featureList.Refresh()

	instance.Sidebar.Add(featureList)
}

// Cleans up the UI elements depending on a current plan.
func (instance *GardenPlanner) ClosePlan() {
	instance.Content.RemoveAll()
	instance.Sidebar.RemoveAll()
	instance.currentPlan = nil
}

// Creates a new instance of the app.
func NewGardenPlanner() *GardenPlanner {
	// Setup UI elements
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Garden Planner")
	sidebar := container.NewVBox()
	content := container.NewVBox()
	toolbar := widget.NewToolbar()
	statusBar := widget.NewLabel("")
	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, content)

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
		currentPlan:   nil,
		App:           mainApp,
		Window:        mainWindow,
		MainContainer: mainContainer,
		Sidebar:       sidebar,
		Toolbar:       toolbar,
		StatusBar:     statusBar,
		Content:       content,
	}

	return &gardenPlanner
}
