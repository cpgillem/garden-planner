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
	Plan  *models.Plan
	scale float32

	// Events
	OnFeatureTapped func(feature *models.Feature)
	OnDragEnd       func()
}

// Create a new feature widget.
func (g *GardenWidget) addFeature(feature *models.Feature) {
	fw := NewFeatureWidget(feature)
	fw.OnTapped = func() {
		g.OnFeatureTapped(fw.Feature)
	}
	// Handle drag events are bubbled up to the garden widget.
	fw.OnHandleDragged = func(edge geometry.BoxEdge, e *fyne.DragEvent) {
		dx := e.Dragged.DX / g.scale
		dy := e.Dragged.DY / g.scale

		// Handle edge cases (lol)
		switch edge {
		case geometry.TOP:
			feature.Box.Location.Y += dy
			feature.Box.Size.Y -= dy
		case geometry.BOTTOM:
			feature.Box.Size.Y += dy
		case geometry.LEFT:
			feature.Box.Location.X += dx
			feature.Box.Size.X -= dx
		case geometry.RIGHT:
			feature.Box.Size.X += dx
		}
		g.Refresh()
	}
	fw.OnHandleDragEnd = func() {
		g.OnDragEnd()
	}
	fw.OnDragged = func(e *fyne.DragEvent) {
		feature.Box.Location = *feature.Box.Location.Add(&geometry.Vector{
			X: e.Dragged.DX / g.scale,
			Y: e.Dragged.DY / g.scale,
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
	g.Plan = plan

	// Add features
	for _, feature := range plan.Features {
		g.addFeature(&feature)
	}

	g.Refresh()
}

// Create a new garden widget. Requires a plan.
func NewGardenWidget(plan *models.Plan) *GardenWidget {
	gardenWidget := &GardenWidget{
		scale: 2,
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
			featureWidget.Feature.Box.Size.X*g.parent.scale,
			featureWidget.Feature.Box.Size.Y*g.parent.scale,
		))
		featureWidget.Move(fyne.NewPos(
			featureWidget.Feature.Box.Location.X*g.parent.scale,
			featureWidget.Feature.Box.Location.Y*g.parent.scale,
		))
	}
}

// MinSize implements fyne.WidgetRenderer.
func (g gardenRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		g.parent.Plan.Box.Size.X*g.parent.scale,
		g.parent.Plan.Box.Size.Y*g.parent.scale,
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
