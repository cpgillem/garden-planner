package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/controllers"
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
)

type GardenWidget struct {
	widget.BaseWidget

	// Feature Widgets
	features map[models.FeatureID]*FeatureWidget

	// Controller reference
	Controller *controllers.PlanController

	// Events
	OnFeatureDragged       func(id models.FeatureID, e *fyne.DragEvent)
	OnFeatureDragEnd       func(id models.FeatureID)
	OnFeatureHandleDragged func(id models.FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent)
	OnFeatureHandleDragEnd func(id models.FeatureID, edge geometry.BoxEdge)
	OnFeatureTapped        func(id models.FeatureID)
}

// Create a new feature widget.
func (g *GardenWidget) AddFeature(id models.FeatureID) {
	fw := NewFeatureWidget(id, g.Controller)
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

// Opens a plan for viewing.
func (g *GardenWidget) OpenPlan(controller *controllers.PlanController) {
	g.Controller = controller

	// Add features
	for i := range controller.Plan.Features {
		g.AddFeature(models.FeatureID(i))
	}

	g.Refresh()
}

// Create a new garden widget. Requires a plan.
func NewGardenWidget(controller *controllers.PlanController) *GardenWidget {
	gardenWidget := &GardenWidget{
		Controller:             controller,
		features:               map[models.FeatureID]*FeatureWidget{},
		OnFeatureDragged:       func(id models.FeatureID, e *fyne.DragEvent) {},
		OnFeatureDragEnd:       func(id models.FeatureID) {},
		OnFeatureHandleDragged: func(id models.FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent) {},
		OnFeatureHandleDragEnd: func(id models.FeatureID, edge geometry.BoxEdge) {},
		OnFeatureTapped:        func(id models.FeatureID) {},
	}

	gardenWidget.OpenPlan(gardenWidget.Controller)

	gardenWidget.ExtendBaseWidget(gardenWidget)
	return gardenWidget
}

func (w *GardenWidget) CreateRenderer() fyne.WidgetRenderer {
	return newGardenRenderer(w)
}

type gardenRenderer struct {
	parent *GardenWidget
}

// Destroy implements fyne.WidgetRenderer.
func (g gardenRenderer) Destroy() {

}

// Layout implements fyne.WidgetRenderer.
func (g gardenRenderer) Layout(fyne.Size) {
	for i := range g.parent.features {
		box := g.parent.Controller.Plan.Features[i].Box
		g.parent.features[i].Resize(fyne.NewSize(
			box.Size.X*g.parent.Controller.DisplayConfig.Scale,
			box.Size.Y*g.parent.Controller.DisplayConfig.Scale,
		))
		g.parent.features[i].Move(fyne.NewPos(
			box.Location.X*g.parent.Controller.DisplayConfig.Scale,
			box.Location.Y*g.parent.Controller.DisplayConfig.Scale,
		))
	}
}

// MinSize implements fyne.WidgetRenderer.
func (g gardenRenderer) MinSize() fyne.Size {
	size := g.parent.Controller.Plan.Box.Size
	return fyne.NewSize(
		size.X*g.parent.Controller.DisplayConfig.Scale,
		size.Y*g.parent.Controller.DisplayConfig.Scale,
	)
}

// Objects implements fyne.WidgetRenderer.
func (g gardenRenderer) Objects() []fyne.CanvasObject {
	os := []fyne.CanvasObject{}
	for _, f := range g.parent.features {
		os = append(os, f)
	}
	return os
}

// Refresh implements fyne.WidgetRenderer.
func (g gardenRenderer) Refresh() {
	for i := range g.parent.features {
		g.parent.features[i].Refresh()
	}

	g.Layout(g.MinSize())
}

func newGardenRenderer(parent *GardenWidget) gardenRenderer {
	gr := gardenRenderer{
		parent: parent,
	}

	return gr
}
