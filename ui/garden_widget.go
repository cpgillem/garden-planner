package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
)

type GardenWidget struct {
	widget.BaseWidget

	// Feature Widgets
	features []*FeatureWidget

	// Internal data
	Scale float32

	// Events
	OnFeatureTapped        func(feature *models.Feature)
	OnFeatureHandleDragged func(feature *models.Feature, edge geometry.BoxEdge, e *fyne.DragEvent)
	OnDragEnd              func()
	OnHandleDragEnd        func()

	// Data hooks
	GetPlanSize func() geometry.Vector
}

// Create a new feature widget.
func (g *GardenWidget) addFeature(feature *models.Feature) {
	fw := NewFeatureWidget(feature)
	fw.OnTapped = func() {
		g.OnFeatureTapped(fw.Feature)
	}
	// Handle drag events are bubbled up to the garden widget.
	fw.OnHandleDragged = func(edge geometry.BoxEdge, e *fyne.DragEvent) {
		g.OnFeatureHandleDragged(fw.Feature, edge, e)
	}
	fw.OnHandleDragEnd = func() {
		g.OnHandleDragEnd()
	}
	fw.OnDragged = func(e *fyne.DragEvent) {
		feature.Box.Location = *feature.Box.Location.Add(&geometry.Vector{
			X: e.Dragged.DX / g.Scale,
			Y: e.Dragged.DY / g.Scale,
			Z: 0,
		})
		g.Refresh()
	}
	fw.OnDragEnd = func() {
		g.OnDragEnd()
	}
	g.features = append(g.features, fw)
}

// Opens a plan for viewing.
func (g *GardenWidget) OpenPlan(plan *models.Plan) {
	// Add features
	for _, feature := range plan.Features {
		g.addFeature(&feature)
	}

	g.Refresh()
}

// Create a new garden widget. Requires a plan.
func NewGardenWidget(plan *models.Plan) *GardenWidget {
	gardenWidget := &GardenWidget{
		Scale:       2,
		GetPlanSize: func() geometry.Vector { return geometry.Vector{X: 0, Y: 0, Z: 0} },
	}

	gardenWidget.OpenPlan(plan)

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
		featureWidget := g.parent.features[i]
		featureWidget.Resize(fyne.NewSize(
			featureWidget.Feature.Box.Size.X*g.parent.Scale,
			featureWidget.Feature.Box.Size.Y*g.parent.Scale,
		))
		featureWidget.Move(fyne.NewPos(
			featureWidget.Feature.Box.Location.X*g.parent.Scale,
			featureWidget.Feature.Box.Location.Y*g.parent.Scale,
		))
	}
}

// MinSize implements fyne.WidgetRenderer.
func (g gardenRenderer) MinSize() fyne.Size {
	size := g.parent.GetPlanSize()
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
	// g.parent.OnRefresh()
}

func newGardenRenderer(parent *GardenWidget) gardenRenderer {
	gr := gardenRenderer{
		parent: parent,
	}

	return gr
}
