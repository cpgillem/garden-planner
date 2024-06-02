package ui

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/controllers"
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
	"golang.org/x/image/colornames"
)

type GardenWidget struct {
	widget.BaseWidget

	// Feature Widgets
	features map[models.FeatureID]*FeatureWidget

	// Background
	background *canvas.Rectangle

	// Gridlines
	hGridlines []*canvas.Line
	vGridlines []*canvas.Line

	// Drawing Settings
	scale       float32
	gridSpacing float32

	// Controller reference
	Controller *controllers.PlanController

	// Events
	OnFeatureDragged       func(id models.FeatureID, e *fyne.DragEvent)
	OnFeatureDragEnd       func(id models.FeatureID)
	OnFeatureHandleDragged func(id models.FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent)
	OnFeatureHandleDragEnd func(id models.FeatureID, edge geometry.BoxEdge)
	OnFeatureTapped        func(id models.FeatureID)
}

// Create a new garden widget. Requires a plan. Agnostic to base units.
//
// controller allows the widget to modify features.
//
// scale is a multiplier on the base unit for display.
//
// gridSpacing defines how many base units between each gridline.
func NewGardenWidget(controller *controllers.PlanController, scale float32, gridSpacing float32) *GardenWidget {
	gardenWidget := &GardenWidget{
		Controller:             controller,
		scale:                  scale,
		gridSpacing:            gridSpacing,
		features:               map[models.FeatureID]*FeatureWidget{},
		OnFeatureDragged:       func(id models.FeatureID, e *fyne.DragEvent) {},
		OnFeatureDragEnd:       func(id models.FeatureID) {},
		OnFeatureHandleDragged: func(id models.FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent) {},
		OnFeatureHandleDragEnd: func(id models.FeatureID, edge geometry.BoxEdge) {},
		OnFeatureTapped:        func(id models.FeatureID) {},
		background:             canvas.NewRectangle(colornames.White),
		hGridlines:             []*canvas.Line{},
		vGridlines:             []*canvas.Line{},
	}

	gardenWidget.OpenPlan(gardenWidget.Controller)
	gardenWidget.Refresh()

	gardenWidget.ExtendBaseWidget(gardenWidget)
	return gardenWidget
}

// Create a new feature widget.
func (g *GardenWidget) AddFeature(id models.FeatureID) {
	fw := NewFeatureWidget(id, g.Controller, g.scale)
	fw.OnDragEnd = func() {
		g.Refresh()
		g.OnFeatureDragEnd(fw.FeatureID)
	}
	fw.OnDragged = func(e *fyne.DragEvent) {
		g.Refresh()
		g.OnFeatureDragged(fw.FeatureID, e)
	}
	fw.OnHandleDragged = func(edge geometry.BoxEdge, e *fyne.DragEvent) {
		g.Refresh()
		g.OnFeatureHandleDragged(fw.FeatureID, edge, e)
	}
	fw.OnHandleDragEnd = func(edge geometry.BoxEdge) {
		g.OnFeatureHandleDragEnd(fw.FeatureID, edge)
	}
	fw.OnTapped = func() {
		g.SelectFeature(fw.FeatureID)
		g.OnFeatureTapped(fw.FeatureID)
	}
	g.features[id] = fw
}

func (g *GardenWidget) RemoveFeature(id models.FeatureID) {
	// Deselect all features.
	g.SelectNone()

	// Remove feature.
	g.features[id].Hide()
}

func (g *GardenWidget) SelectFeature(id models.FeatureID) {
	g.SelectNone()

	// Select this feature.
	g.features[id].Select()
	g.Refresh()
}

func (g *GardenWidget) SelectNone() {
	// Deselect other features.
	for i := range g.features {
		g.features[i].Deselect()
	}
}

// Recreates the gridline cache upon changes to the plan.
func (g *GardenWidget) CalculateGridlines() {
	// Calculate gridline counts.
	vGridlines := int(math.Floor(float64(g.Controller.Plan.Box.GetWidth()) / float64(g.gridSpacing)))
	hGridlines := int(math.Floor(float64(g.Controller.Plan.Box.GetHeight()) / float64(g.gridSpacing)))

	// Empty gridline cache.
	g.hGridlines = []*canvas.Line{}
	g.vGridlines = []*canvas.Line{}

	// Add un-laid-out gridlines to cache.
	for i := 0; i < vGridlines; i++ {
		g.vGridlines = append(g.vGridlines, canvas.NewLine(colornames.Gray))
	}

	for i := 0; i < hGridlines; i++ {
		g.hGridlines = append(g.hGridlines, canvas.NewLine(colornames.Gray))
	}
}

func (g *GardenWidget) SetScale(s float32) {
	g.scale = s
	g.Refresh()
}

func (g *GardenWidget) SetGridSpacing(s float32) {
	g.gridSpacing = s
	g.Refresh()
}

// Opens a plan for viewing.
func (g *GardenWidget) OpenPlan(controller *controllers.PlanController) {
	g.Controller = controller

	// Calculate gridlines.
	g.CalculateGridlines()

	// Add features
	for i := range controller.Plan.Features {
		g.AddFeature(models.FeatureID(i))
	}

	g.Refresh()
}

func (w *GardenWidget) CreateRenderer() fyne.WidgetRenderer {
	return newGardenRenderer(w)
}

// Events

// On scroll, adjust the scale.
func (w *GardenWidget) Scrolled(e *fyne.ScrollEvent) {
	adjDY := e.Scrolled.DY / 250
	w.SetScale(w.scale + adjDY)
}

type gardenRenderer struct {
	parent *GardenWidget

	hGridlines uint
	vGridlines uint
}

func newGardenRenderer(parent *GardenWidget) gardenRenderer {
	gr := gardenRenderer{
		parent:     parent,
		hGridlines: 0,
		vGridlines: 0,
	}

	return gr
}

// Destroy implements fyne.WidgetRenderer.
func (g gardenRenderer) Destroy() {

}

// Layout implements fyne.WidgetRenderer.
func (g gardenRenderer) Layout(s fyne.Size) {
	// Layout background.
	g.parent.background.StrokeWidth = 1
	g.parent.background.StrokeColor = colornames.Black
	g.parent.background.Move(fyne.NewPos(0, 0))
	g.parent.background.Resize(fyne.NewSize(
		g.parent.Controller.Plan.Box.GetWidth()*g.parent.scale,
		g.parent.Controller.Plan.Box.GetHeight()*g.parent.scale,
	))

	// Layout horizontal gridlines.
	for i := range g.parent.hGridlines {
		y := float32(float64(i) * float64(g.parent.gridSpacing) * float64(g.parent.scale))
		var left float32 = 0
		var right float32 = g.parent.Controller.Plan.Box.GetWidth() * g.parent.scale
		g.parent.hGridlines[i].Position1 = fyne.NewPos(
			left,
			y,
		)
		g.parent.hGridlines[i].Position2 = fyne.NewPos(
			right,
			y,
		)
	}

	// Layout vertical gridlines.
	for i := range g.parent.vGridlines {
		x := float32(float64(i) * float64(g.parent.gridSpacing) * float64(g.parent.scale))
		var top float32 = 0
		var bottom float32 = g.parent.Controller.Plan.Box.GetHeight() * g.parent.scale
		g.parent.vGridlines[i].Position1 = fyne.NewPos(
			x,
			top,
		)
		g.parent.vGridlines[i].Position2 = fyne.NewPos(
			x,
			bottom,
		)
	}

	// Layout features.
	for i := range g.parent.features {
		box := g.parent.Controller.Plan.Features[i].Box
		g.parent.features[i].Resize(fyne.NewSize(
			box.Size.X*g.parent.scale,
			box.Size.Y*g.parent.scale,
		))
		g.parent.features[i].Move(fyne.NewPos(
			box.Location.X*g.parent.scale,
			box.Location.Y*g.parent.scale,
		))
	}

}

// MinSize implements fyne.WidgetRenderer.
func (g gardenRenderer) MinSize() fyne.Size {
	size := g.parent.Controller.Plan.Box.Size
	return fyne.NewSize(
		size.X*g.parent.scale,
		size.Y*g.parent.scale,
	)
}

// Objects implements fyne.WidgetRenderer.
func (g gardenRenderer) Objects() []fyne.CanvasObject {
	os := []fyne.CanvasObject{}

	// Add background objects.
	os = append(os, g.parent.background)

	// Add gridlines.
	for _, g := range g.parent.hGridlines {
		os = append(os, g)
	}
	for _, g := range g.parent.vGridlines {
		os = append(os, g)
	}

	// Add features.
	for _, f := range g.parent.features {
		os = append(os, f)
	}
	return os
}

// Refresh implements fyne.WidgetRenderer.
func (g gardenRenderer) Refresh() {
	g.parent.CalculateGridlines()

	for i := range g.parent.features {
		g.parent.features[i].Refresh()
	}

	g.Layout(g.MinSize())
}
