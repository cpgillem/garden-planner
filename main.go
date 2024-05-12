package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/ui"
)

func main() {
	mainApp := app.New()

	mainWindow := mainApp.NewWindow("Garden Planner")

	sidebar := container.NewVBox()
	garden := container.New(&ui.GardenLayout{})

	toolbar := widget.NewToolbar()

	mainContainer := container.NewBorder(toolbar, nil, sidebar, nil, garden)

	mainWindow.SetContent(mainContainer)

	mainWindow.Show()
	mainApp.Run()
}
