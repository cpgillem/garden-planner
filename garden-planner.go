package main

import (
	"github.com/cpgillem/garden-planner/widgets"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	mainApp := app.New()

	mainWindow := mainApp.NewWindow("Garden Planner")

	sidebar := layout.NewVBoxLayout()
	toolbar := widget.NewToolbar()
	visual := widgets.NewGardenLayout()

	mainWindow.Show()
	mainApp.Run()
}
