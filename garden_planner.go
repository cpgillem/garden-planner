package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	Toolbar       *widget.Toolbar
	StatusBar     *widget.Label
	PropertyTable *fyne.Container
	FeatureList   *widget.List

	// Data
	CurrentPlan *models.Plan
	GardenData  *GardenData
}

func (p *GardenPlanner) Start() {
	p.Window.Show()
	p.App.Run()
}

func (instance *GardenPlanner) OpenPlan(plan *models.Plan) {
	instance.ClosePlan()
	instance.CurrentPlan = plan

	// Create feature widgets.
	gardenContainer := container.New(ui.NewGardenLayout(plan))
	for _, feature := range plan.Features {
		featureWidget := ui.NewFeatureWidget(&feature)
		gardenContainer.Add(featureWidget)
	}
	instance.Content.Add(gardenContainer)
	gardenContainer.Refresh()

	// Setup Feature List

	// Feature length comes from the plan.
	featuresLength := func() int {
		if instance.CurrentPlan == nil {
			return 0
		}
		return int(len(instance.CurrentPlan.Features))
	}

	// When a feature is added, create a label.
	createFeature := func() fyne.CanvasObject {
		// Fix later: why doesn't fyne refresh the width of the sidebar properly?
		label := widget.NewLabel("aaaaaaaaaaaaaaaaaaaaaaaaaaaa")

		return label
	}

	// When a feature is updated, set its text.
	updateFeature := func(id widget.ListItemID, o fyne.CanvasObject) {
		obj := o.(*widget.Label)
		obj.SetText(instance.CurrentPlan.Features[id].Name)
	}

	// Recreate feature list. Fyne doesn't have an obvious way to remove an item.
	instance.FeatureList = widget.NewList(featuresLength, createFeature, updateFeature)

	// When a feature is selected, display its properties.
	instance.FeatureList.OnSelected = func(id widget.ListItemID) {
		instance.SelectFeature(&instance.CurrentPlan.Features[id], instance.GardenData.Properties)
	}

	instance.FeatureList.Refresh()

	instance.Sidebar.Add(instance.FeatureList)
	instance.Sidebar.Add(instance.PropertyTable)
}

// Updates the GUI when a feature is selected.
func (instance *GardenPlanner) SelectFeature(feature *models.Feature, properties map[string]models.Property) {
	instance.PropertyTable.RemoveAll()

	for propertyName := range feature.Properties {
		label := widget.NewLabel(properties[propertyName].DisplayName)
		entry := instance.CreatePropertyWidget(properties[propertyName], feature)
		instance.PropertyTable.Add(label)
		instance.PropertyTable.Add(entry)
	}

	instance.PropertyTable.Refresh()
}

// Creates a widget for modifying a property on a feature.
func (instance *GardenPlanner) CreatePropertyWidget(property models.Property, feature *models.Feature) *widget.Entry {
	// TODO: Custom widgets for property types.
	value := feature.Properties[property.Name]
	str := ""

	switch property.PropertyType {
	case "dimension", "decimal":
		str = fmt.Sprintf("%.3f", value)
	case "text":
		str = value.(string)
	case "integer":
		str = strconv.Itoa(value.(int))
	}
	entry := widget.NewEntry()
	entry.SetText(str)

	// Setup events.
	entry.OnChanged = func(s string) {
		switch property.PropertyType {
		case "dimension", "decimal":
			setValue, err := strconv.ParseFloat(s, 32)
			if err != nil {
				dialog.ShowInformation("Error", "Invalid number.", instance.Window)
				break
			}
			feature.Properties[property.Name] = setValue
		case "text":
			feature.Properties[property.Name] = s
		case "integer":
			setValue, err := strconv.Atoi(s)
			if err != nil {
				dialog.ShowInformation("Error", "Invalid number.", instance.Window)
				break
			}
			feature.Properties[property.Name] = setValue
		}
	}

	return entry
}

// Cleans up the UI elements depending on a current plan.
func (instance *GardenPlanner) ClosePlan() {
	instance.Content.RemoveAll()
	instance.Sidebar.RemoveAll()
	instance.CurrentPlan = nil
}

// Creates a new instance of the app.
func NewGardenPlanner(gardenData *GardenData) *GardenPlanner {
	// Setup UI elements
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Garden Planner")
	sidebar := container.NewVBox()
	content := container.NewVBox()
	toolbar := widget.NewToolbar()
	statusBar := widget.NewLabel("")
	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, content)
	featureList := widget.NewList(func() int { return 0 }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(lii widget.ListItemID, co fyne.CanvasObject) {})
	propertyTable := container.NewGridWithColumns(2)

	mainWindow.SetContent(mainContainer)

	// Create new app instance
	gardenPlanner := GardenPlanner{
		CurrentPlan:   nil,
		App:           mainApp,
		Window:        mainWindow,
		MainContainer: mainContainer,
		Sidebar:       sidebar,
		Toolbar:       toolbar,
		StatusBar:     statusBar,
		Content:       content,
		FeatureList:   featureList,
		PropertyTable: propertyTable,
		GardenData:    gardenData,
	}

	// Setup Toolbar
	createButton := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		// Open an empty plan.
		gardenPlanner.OpenPlan(&models.Plan{})
	})
	toolbar.Append(createButton)
	toolbar.Append(widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
		// Display the file open dialog.
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				// Read the plan and open it.
				plan, _ := ReadObject[models.Plan](reader)
				gardenPlanner.OpenPlan(plan)
			}
		}, gardenPlanner.Window)
	}))
	toolbar.Append(widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		// Display the file save dialog.
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				// Save the plan.
				WriteObject(writer, gardenPlanner.CurrentPlan)
			}
		}, gardenPlanner.Window)
	}))
	toolbar.Append(widget.NewToolbarSeparator())

	return &gardenPlanner
}
