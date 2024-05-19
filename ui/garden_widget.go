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
	features []*FeatureWidget

	// Internal data
	Scale float32

	// Controller reference
	Controller *controllers.PlanController

	// Events
	OnFeatureDragged       func(id models.FeatureID, e *fyne.DragEvent)
	OnFeatureDragEnd       func(id models.FeatureID)
	OnFeatureHandleDragged func(id models.FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent)
	OnFeatureHandleDragEnd func(id models.FeatureID, edge geometry.BoxEdge)
}

// Create a new feature widget.
func (g *GardenWidget) addFeature(id models.FeatureID) {
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
	g.features = append(g.features, fw)
}

// Opens a plan for viewing.
func (g *GardenWidget) OpenPlan(controller *controllers.PlanController) {
	g.Controller = controller

	// Add features
	for i := range controller.Plan.Features {
		g.addFeature(models.FeatureID(i))
	}

	g.Refresh()
}

// Create a new garden widget. Requires a plan.
func NewGardenWidget(controller *controllers.PlanController) *GardenWidget {
	gardenWidget := &GardenWidget{
		Scale:                  2,
		Controller:             controller,
		OnFeatureDragged:       func(id models.FeatureID, e *fyne.DragEvent) {},
		OnFeatureDragEnd:       func(id models.FeatureID) {},
		OnFeatureHandleDragged: func(id models.FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent) {},
		OnFeatureHandleDragEnd: func(id models.FeatureID, edge geometry.BoxEdge) {},
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
	for _, f := range g.parent.features {
		box := g.parent.Controller.Plan.Features[f.FeatureID].Box
		f.Resize(fyne.NewSize(
			box.Size.X*g.parent.Scale,
			box.Size.Y*g.parent.Scale,
		))
		f.Move(fyne.NewPos(
			box.Location.X*g.parent.Scale,
			box.Location.Y*g.parent.Scale,
		))
	}
}

// MinSize implements fyne.WidgetRenderer.
func (g gardenRenderer) MinSize() fyne.Size {
	size := g.parent.Controller.Plan.Box.Size
	return fyne.NewSize(
		size.X*g.parent.Scale,
		size.Y*g.parent.Scale,
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
