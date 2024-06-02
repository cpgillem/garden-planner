package main

import (
	"fmt"
	"log"

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

const IMPERIAL string = "imperial"
const METRIC string = "metric"

// Represents the state of the application.
type GardenPlanner struct {
	App fyne.App

	// Windows
	Window         fyne.Window
	SettingsWindow SettingsWindow

	// Containers
	MainContainer *fyne.Container
	Sidebar       *fyne.Container
	GardenWidget  *ui.GardenWidget

	// Controllers
	PlanController  controllers.PlanController
	PlantController controllers.PlantController

	// Permanent Widgets
	Toolbar       *widget.Toolbar
	StatusBar     *widget.Label
	PropertyTable *fyne.Container
	FeatureTools  *fyne.Container
	BoxEditor     *ui.BoxEditor

	// Button References
	DeleteFeature    *widget.Button
	TemplateSelector *widget.Select

	// Data
	GardenData    *GardenData
	Formatter     *ui.DimensionFormatter
	DisplayConfig *models.DisplayConfig
}

// Creates a new instance of the app.
func NewGardenPlanner() *GardenPlanner {
	// Setup UI elements
	mainApp := app.NewWithID("net.cpgworld.garden-planner.preferences")

	// Windows
	mainWindow := mainApp.NewWindow("Garden Planner")

	// Main Window
	gardenData := NewGardenData()

	// Load plant data.
	plants, err := ReadObjectFromFile[[]models.Plant]("data/plants.json")
	if err != nil {
		log.Fatal("Could not load plant data.")
		plants = &[]models.Plant{}
	}
	plantController := controllers.NewPlantController(plants)

	displayConfig := models.NewDisplayConfig()
	formatter := ui.NewFormatter()
	gridSpacing, err := formatter.ToDimensionBaseUnit(
		mainApp.Preferences().StringWithFallback("grid_spacing", "12 in"),
		displayConfig.BaseUnit,
	)
	if err != nil {
		gridSpacing = units.NewValue(12, units.Inch)
	}
	sidebar := container.NewVBox()
	blankPlan := models.NewPlan()
	planController := controllers.NewPlanController(blankPlan)
	gardenWidget := ui.NewGardenWidget(
		&planController,
		float32(mainApp.Preferences().FloatWithFallback("display_scale", 2)),
		float32(gridSpacing.Float()),
	)
	toolbar := widget.NewToolbar()
	statusBar := widget.NewLabel("")
	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, gardenWidget)
	propertyTable := container.New(layout.NewFormLayout())
	featureTools := container.NewHBox()
	boxEditor := ui.NewBoxEditor(geometry.NewBoxZero(), ui.AnyUnit, formatter)

	mainWindow.SetContent(mainContainer)

	// Create new app instance
	gardenPlanner := GardenPlanner{
		App:             mainApp,
		Window:          mainWindow,
		MainContainer:   mainContainer,
		Sidebar:         sidebar,
		BoxEditor:       boxEditor,
		Toolbar:         toolbar,
		StatusBar:       statusBar,
		GardenWidget:    gardenWidget,
		FeatureTools:    featureTools,
		PropertyTable:   propertyTable,
		GardenData:      gardenData,
		Formatter:       formatter,
		PlanController:  planController,
		PlantController: plantController,
		DisplayConfig:   &displayConfig,
	}

	// Setup Toolbar
	gardenPlanner.SetupToolbar()
	gardenPlanner.SetupFeatureTools()

	// Other windows
	gardenPlanner.SettingsWindow = NewSettingsWindow(&gardenPlanner)

	mainApp.Preferences().AddChangeListener(gardenPlanner.RereadSettings)

	return &gardenPlanner
}

// After settings are changed, make the appropriate updates.
func (p *GardenPlanner) RereadSettings() {
	spacing := p.App.Preferences().StringWithFallback("grid_spacing", "12 in")
	spacingUnit, err := p.Formatter.ToDimensionBaseUnit(spacing, p.DisplayConfig.BaseUnit)
	if err == nil {
		p.GardenWidget.SetGridSpacing(float32(spacingUnit.Float()))
	}
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
	instance.BoxEditor.SetBox(instance.PlanController.Plan.Features[id].Box)
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
	instance.PlanController = controllers.NewPlanController(plan)
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
	boxEditor := ui.NewBoxEditor(feature.Box, units.Inch, instance.Formatter)
	boxEditor.OnSubmitted = func(newBox geometry.Box) {
		feature.Box = newBox
		instance.GardenWidget.Refresh()
	}
	instance.BoxEditor = boxEditor

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
		entry, err := instance.CreatePropertyWidget(instance.GardenData.Properties[propertyName], feature)
		if err != nil {
			// Don't add anything if the property can't be read.
			fmt.Printf("Warning: %s\n", err.Error())
			continue
		}
		instance.PropertyTable.Add(label)
		instance.PropertyTable.Add(entry)
	}
}

// Creates a widget for modifying a property on a feature.
func (instance *GardenPlanner) CreatePropertyWidget(property models.Property, feature *models.Feature) (fyne.Widget, error) {
	// TODO: Custom widgets for property types.
	// TODO: formatting parameters.
	value := feature.Properties[property.Name]

	switch property.PropertyType {
	case "dimension":
		// Dimensions are stored as strings in files.
		val, err := instance.Formatter.ToDimension(value.(string))
		if err != nil {
			fmt.Printf("Warning on property %s.%s: %s\n", feature.Name, property.Name, err.Error())
		}
		entry := ui.NewDimensionEntry(val, instance.Formatter)
		entry.OnDimensionError = func(err error) {
			dialog.ShowError(err, instance.Window)
		}
		entry.OnValueChanged = func(val units.Value) {
			feature.Properties[property.Name] = instance.Formatter.FormatDimension(val)
			instance.MainContainer.Refresh()
		}
		return entry, nil
	case "decimal":
		// TODO: Numerical entry widget.
		entry := widget.NewEntry()
		entry.SetText(instance.Formatter.FormatDecimal(float32(value.(float64))))
		entry.OnSubmitted = func(s string) {
			setValue, err := instance.Formatter.ToDecimal(s)
			if err != nil {
				dialog.ShowError(err, instance.Window)
				return
			}
			feature.Properties[property.Name] = setValue
			instance.MainContainer.Refresh()
		}
		return entry, nil
	case "integer":
		entry := widget.NewEntry()
		entry.SetText(instance.Formatter.FormatInteger(value.(int)))
		entry.OnSubmitted = func(s string) {
			setValue, err := instance.Formatter.ToInteger(s)
			if err != nil {
				dialog.ShowError(err, instance.Window)
				return
			}
			feature.Properties[property.Name] = setValue
			instance.MainContainer.Refresh()
		}
		return entry, nil
	case "string":
		// Should be a string.
		entry := widget.NewEntry()
		entry.SetText(value.(string))
		entry.OnSubmitted = func(s string) {
			feature.Properties[property.Name] = s
			instance.MainContainer.Refresh()
		}
		return entry, nil
	default:
		// Return a disabled entry and throw an error.
		entry := widget.NewEntry()
		entry.Disable()
		return entry, fmt.Errorf("type not recognized on property %s: %s", property.Name, property.PropertyType)
	}
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

	// Settings
	instance.Toolbar.Append(widget.NewToolbarAction(theme.SettingsIcon(), func() {
		// Display the settings window.
		instance.SettingsWindow.Show()
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
