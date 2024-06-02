package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bcicen/go-units"
	"github.com/cpgillem/garden-planner/ui"
)

type SettingsWindow struct {
	instance *GardenPlanner
	window   fyne.Window

	systemEntry *widget.SelectEntry
	gridEntry   *ui.DimensionEntry

	okButton     *widget.Button
	cancelButton *widget.Button

	// Events
	OnOk func()
}

func NewSettingsWindow(instance *GardenPlanner) SettingsWindow {
	w := SettingsWindow{
		instance:    instance,
		window:      instance.App.NewWindow("Settings"),
		systemEntry: widget.NewSelectEntry([]string{"Imperial", "Metric"}),
		gridEntry: ui.NewDimensionEntry(
			units.NewValue(0, ui.AnyUnit),
			instance.Formatter,
		),
		okButton:     widget.NewButton("OK", func() {}),
		cancelButton: widget.NewButton("Cancel", func() {}),
		OnOk:         func() {},
	}

	// Events
	w.okButton.OnTapped = w.OK
	w.cancelButton.OnTapped = w.Close

	systemLabel := widget.NewLabel("Measurement System")
	gridLabel := widget.NewLabel("Grid Spacing")

	// Containers
	measurementForm := container.New(
		layout.NewFormLayout(),
		systemLabel,
		w.systemEntry,
		gridLabel,
		w.gridEntry,
	)
	measurementTab := container.NewTabItem("Measurement", measurementForm)
	settingsTabs := container.NewAppTabs(measurementTab)
	buttonContainer := container.NewHBox(w.cancelButton, w.okButton)
	settingsWinContainer := container.NewVBox(settingsTabs, buttonContainer)
	w.window.SetContent(settingsWinContainer)

	return w
}

func (w *SettingsWindow) Show() {
	// Measurement system
	switch w.instance.App.Preferences().StringWithFallback("measurement_system", IMPERIAL) {
	case IMPERIAL:
		w.systemEntry.SetText("Imperial")
	case METRIC:
		w.systemEntry.SetText("Metric")
	default:
		fmt.Println("Could not read measurement system setting.")
		w.systemEntry.SetText("Imperial")
	}

	// Grid spacing
	v, err := w.instance.Formatter.ToDimension(w.instance.App.Preferences().StringWithFallback("grid_spacing", "12 in"))
	if err != nil {
		fmt.Println(err.Error())
		w.gridEntry.SetValueAndBaseUnit(units.NewValue(12, units.Inch))
	} else {
		w.gridEntry.SetValueAndBaseUnit(v)
	}

	// Show window
	w.window.Show()
}

func (w *SettingsWindow) OK() {
	// Measurement System
	switch w.systemEntry.Text {
	case "Imperial":
		w.instance.App.Preferences().SetString("measurement_system", IMPERIAL)
	case "Metric":
		w.instance.App.Preferences().SetString("measurement_system", METRIC)
	default:
		dialog.ShowInformation("Validation Error", "Incorrect measurement system.", w.window)
		return
	}

	// Grid spacing
	w.instance.App.Preferences().SetString("grid_spacing", w.gridEntry.GetValueAsText())

	w.Close()
	w.OnOk()
}

func (w *SettingsWindow) Close() {
	w.window.Close()
}
