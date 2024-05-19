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
	OnFeatureTapped        func(id FeatureID)
	OnFeatureHandleDragged func(id FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent)
	OnHandleDragEnd        func()
	OnFeatureDragged       func(id FeatureID, e *fyne.DragEvent)
	OnFeatureDragEnd       func()
	OnFeatureRefresh       func(id FeatureID)

	// Data hooks
	GetPlanSize    func() geometry.Vector
	GetFeatureBox  func(id FeatureID) geometry.AxisAlignedBoundingBox
	GetFeatureName func(id FeatureID) string
}

// Create a new feature widget.
func (g *GardenWidget) addFeature(id FeatureID) {
	fw := NewFeatureWidget(id)
	fw.OnTapped = func() {
		g.OnFeatureTapped(id)
	}
	// Handle drag events are bubbled up to the garden widget.
	fw.OnHandleDragged = func(id FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent) {
		g.OnFeatureHandleDragged(id, edge, e)
	}
	fw.OnHandleDragEnd = func() {
		g.OnHandleDragEnd()
	}
	fw.OnDragged = func(id FeatureID, e *fyne.DragEvent) {
		g.OnFeatureDragged(id, e)
	}
	fw.OnDragEnd = func() {
		g.OnFeatureDragEnd()
	}
	fw.OnRefresh = func(id FeatureID) {
		g.OnFeatureRefresh(id)
	}
	fw.GetName = func(id FeatureID) string {
		return g.GetFeatureName(id)
	}
	g.features = append(g.features, fw)
}

// Opens a plan for viewing.
func (g *GardenWidget) OpenPlan(plan *models.Plan) {
	// Add features
	for i := range plan.Features {
		g.addFeature(FeatureID(i))
	}

	g.Refresh()
}

// Create a new garden widget. Requires a plan.
func NewGardenWidget(plan *models.Plan) *GardenWidget {
	gardenWidget := &GardenWidget{
		Scale:                  2,
		GetPlanSize:            func() geometry.Vector { return geometry.NewVector(0, 0, 0) },
		OnFeatureTapped:        func(id FeatureID) {},
		OnFeatureHandleDragged: func(id FeatureID, edge geometry.BoxEdge, e *fyne.DragEvent) {},
		OnFeatureDragged:       func(id FeatureID, e *fyne.DragEvent) {},
		OnFeatureDragEnd:       func() {},
		GetFeatureBox: func(id FeatureID) geometry.AxisAlignedBoundingBox {
			return geometry.AxisAlignedBoundingBox{
				Location: geometry.NewVector(0, 0, 0),
				Size:     geometry.NewVector(0, 0, 0),
			}
		},
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
	for _, f := range g.parent.features {
		box := g.parent.GetFeatureBox(f.FeatureID)
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
