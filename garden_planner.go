package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bcicen/go-units"
	"github.com/cpgillem/garden-planner/controllers"
	"github.com/cpgillem/garden-planner/geometry"
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
	GardenWidget  *ui.GardenWidget

	// Controllers
	PlanController controllers.PlanController

	// Permanent Widgets
	Toolbar       *widget.Toolbar
	StatusBar     *widget.Label
	PropertyTable *fyne.Container
	FeatureTools  *fyne.Container

	// Button References
	DeleteFeature    *widget.Button
	TemplateSelector *widget.Select

	// Data
	GardenData    *GardenData
	Formatter     *ui.Formatter
	DisplayConfig *models.DisplayConfig
}

// Creates a new instance of the app.
func NewGardenPlanner(gardenData *GardenData) *GardenPlanner {
	// Setup UI elements
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Garden Planner")
	displayConfig := models.NewDisplayConfig()
	sidebar := container.NewVBox()
	blankPlan := models.NewPlan()
	planController := controllers.NewPlanController(blankPlan, &displayConfig)
	gardenWidget := ui.NewGardenWidget(&planController)
	toolbar := widget.NewToolbar()
	statusBar := widget.NewLabel("")
	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, gardenWidget)
	propertyTable := container.New(layout.NewFormLayout())
	// Use inches as a base unit for now.
	formatter := ui.NewFormatter(units.Inch)
	featureTools := container.NewHBox()

	mainWindow.SetContent(mainContainer)

	// Create new app instance
	gardenPlanner := GardenPlanner{
		App:            mainApp,
		Window:         mainWindow,
		MainContainer:  mainContainer,
		Sidebar:        sidebar,
		Toolbar:        toolbar,
		StatusBar:      statusBar,
		GardenWidget:   gardenWidget,
		FeatureTools:   featureTools,
		PropertyTable:  propertyTable,
		GardenData:     gardenData,
		Formatter:      formatter,
		PlanController: planController,
		DisplayConfig:  &displayConfig,
	}

	// Setup Toolbar
	gardenPlanner.SetupToolbar()
	gardenPlanner.SetupFeatureTools()

	return &gardenPlanner
}

func (p *GardenPlanner) Start() {
	p.Window.Show()
	p.App.Run()
}

func (instance *GardenPlanner) FeatureSelected(id models.FeatureID) {
	instance.SelectFeature(id)
}

func (instance *GardenPlanner) FeatureAdded(id models.FeatureID) {
	instance.GardenWidget.AddFeature(id)
	instance.SelectFeature(id)
}

func (instance *GardenPlanner) FeatureRemoved(id models.FeatureID) {
	if instance.PlanController.GetSelectedFeature() == id {
		instance.PropertyTable.RemoveAll()
		instance.DeleteFeature.Disable()
	}
	instance.GardenWidget.RemoveFeature(id)
}

func (instance *GardenPlanner) FeatureDragEnd(id models.FeatureID) {
	instance.Sidebar.Refresh()
}

func (instance *GardenPlanner) FeatureHandleDragEnd(id models.FeatureID, edge geometry.BoxEdge) {
	instance.Sidebar.Refresh()
}

func (instance *GardenPlanner) OpenPlan(plan *models.Plan) {
	instance.ClosePlan()

	// instance.SetupFeatureList()

	instance.Sidebar.Add(instance.FeatureTools)
	// instance.Sidebar.Add(instance.FeatureList)
	instance.Sidebar.Add(instance.PropertyTable)

	// TODO: Make displayconfig loadable from a file.

	// Setup Plan controller.
	instance.PlanController = controllers.NewPlanController(plan, instance.DisplayConfig)
	instance.PlanController.OnFeatureSelected = instance.FeatureSelected
	instance.PlanController.OnFeatureAdded = instance.FeatureAdded
	instance.PlanController.OnFeatureRemoved = instance.FeatureRemoved

	// Setup garden viewer widget.
	instance.GardenWidget.OpenPlan(&instance.PlanController)
	instance.GardenWidget.OnFeatureDragEnd = instance.FeatureDragEnd
	instance.GardenWidget.OnFeatureHandleDragEnd = instance.FeatureHandleDragEnd

	// Enable necessary feature buttons.
	if instance.PlanController.HasSelection() {
		instance.DeleteFeature.Enable()
	}
	instance.TemplateSelector.Enable()
}

// Updates the GUI when a feature is selected.
func (instance *GardenPlanner) SelectFeature(id models.FeatureID) {
	instance.PropertyTable.RemoveAll()
	instance.AddFeatureProperties(id)

	instance.GardenWidget.SelectFeature(id)

	// Enable feature editing buttons
	instance.DeleteFeature.Enable()

	instance.PropertyTable.Refresh()
	instance.GardenWidget.Refresh()
}

func (instance *GardenPlanner) ClearFeatureProperties() {
	instance.PropertyTable.RemoveAll()
}

// Adds the base properties of any landscaping feature to the properties panel.
func (instance *GardenPlanner) AddFeatureProperties(id models.FeatureID) {
	feature := instance.PlanController.Plan.Features[id]
	nameLabel := widget.NewLabel("Name")
	boxLabel := widget.NewLabel("Box")
	boxEditor := ui.NewBoxEditor(&feature.Box, instance.Formatter)
	boxEditor.OnSubmitted = func(boxDelta geometry.Box) {
		instance.GardenWidget.Refresh()
	}

	// Base built-in properties.
	nameEntry := widget.NewEntry()
	nameEntry.MultiLine = false
	nameEntry.SetText(feature.Name)
	nameEntry.OnSubmitted = func(s string) {
		feature.Name = s
		instance.GardenWidget.Refresh()
	}

	instance.PropertyTable.Add(nameLabel)
	instance.PropertyTable.Add(nameEntry)
	instance.PropertyTable.Add(boxLabel)
	instance.PropertyTable.Add(boxEditor)

	// Custom properties on feature.
	for propertyName := range feature.Properties {
		label := widget.NewLabel(instance.GardenData.Properties[propertyName].DisplayName)
		entry := instance.CreatePropertyWidget(instance.GardenData.Properties[propertyName], feature)
		instance.PropertyTable.Add(label)
		instance.PropertyTable.Add(entry)
	}
}

// Creates a widget for modifying a property on a feature.
func (instance *GardenPlanner) CreatePropertyWidget(property models.Property, feature *models.Feature) *widget.Entry {
	// TODO: Custom widgets for property types.
	// TODO: formatting parameters.
	value := feature.Properties[property.Name]
	str := ""

	switch property.PropertyType {
	case "dimension":
		str = instance.Formatter.FormatDimension(float32(value.(float64)))
	case "decimal":
		str = instance.Formatter.FormatDecimal(float32(value.(float64)))
	case "text":
		str = value.(string)
	case "integer":
		str = instance.Formatter.FormatInteger(value.(int))
	}
	entry := widget.NewEntry()
	entry.SetText(str)

	// Setup events.
	entry.OnSubmitted = func(s string) {
		switch property.PropertyType {
		case "dimension":
			setValue, err := instance.Formatter.ToDimension(s)
			if err != nil {
				dialog.ShowError(err, instance.Window)
				break
			}
			feature.Properties[property.Name] = setValue
		case "decimal":
			setValue, err := instance.Formatter.ToDecimal(s)
			if err != nil {
				dialog.ShowError(err, instance.Window)
				break
			}
			feature.Properties[property.Name] = setValue
		case "text":
			feature.Properties[property.Name] = s
		case "integer":
			setValue, err := instance.Formatter.ToInteger(s)
			if err != nil {
				dialog.ShowError(err, instance.Window)
				break
			}
			feature.Properties[property.Name] = setValue
		}

		instance.MainContainer.Refresh()
	}

	return entry
}

// Cleans up the UI elements depending on a current plan.
func (instance *GardenPlanner) ClosePlan() {
	// instance.Content.RemoveAll()
	instance.PropertyTable.RemoveAll()
	instance.DeleteFeature.Disable()
	instance.TemplateSelector.Disable()
}

func (instance *GardenPlanner) SetupToolbar() {
	// Create file
	createButton := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		// Open an empty plan.
		instance.OpenPlan(&models.Plan{})
	})
	instance.Toolbar.Append(createButton)

	// Open file
	instance.Toolbar.Append(widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
		// Display the file open dialog.
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				// Read the plan and open it.
				plan, _ := ReadObject[models.Plan](reader)
				instance.OpenPlan(plan)
			}
		}, instance.Window)
	}))

	// Save file
	instance.Toolbar.Append(widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		// Display the file save dialog.
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				// Save the plan.
				WriteObject(writer, instance.PlanController.Plan)
			}
		}, instance.Window)
	}))
}

func (instance *GardenPlanner) SetupFeatureTools() {
	// Setup template selector for new features.
	templateNames := []string{}
	for _, t := range instance.GardenData.FeatureTemplates {
		// TODO: More robust template selector
		templateNames = append(templateNames, t.Name)
	}
	instance.TemplateSelector = widget.NewSelect(templateNames, func(s string) {
		t := instance.GardenData.FeatureTemplates[s]
		f := models.NewFeature(instance.GardenData.Properties, &t)
		instance.PlanController.AddFeature(f)
	})
	instance.TemplateSelector.Disable()
	instance.TemplateSelector.PlaceHolder = "New Feature..."
	instance.FeatureTools.Add(instance.TemplateSelector)

	// Setup Feature delete button.
	instance.DeleteFeature = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), instance.PlanController.RemoveSelected)
	instance.DeleteFeature.Disable()
	instance.FeatureTools.Add(instance.DeleteFeature)
}

// Creates an entry widget with unit display, conversion, and arithmetic(?) features.
func NewDimensionEntry(formatter *ui.Formatter, unit *units.Unit, initial float32) *widget.Entry {
	entry := widget.NewEntry()

	// Set defaults
	entry.PlaceHolder = unit.Symbol

	return entry
}
