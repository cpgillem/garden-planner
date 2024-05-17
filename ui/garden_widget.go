package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
)

type GardenWidget struct {
	widget.BaseWidget

	// References to feature widgets added to container
	container *fyne.Container

	// Internal data
	Plan *models.Plan

	// Events
	OnFeatureTapped func(feature *models.Feature)
}

// Create a new feature widget.
func (g *GardenWidget) addFeature(feature *models.Feature) {
	fw := NewFeatureWidget(feature)
	fw.OnTapped = func() {
		g.OnFeatureTapped(fw.Feature)
	}
	g.container.Add(fw)
}

// Opens a plan for viewing.
func (g *GardenWidget) OpenPlan(plan *models.Plan) {
	g.Plan = plan

	// Add initial feature widgets.
	for _, f := range plan.Features {
		g.addFeature(&f)
	}

	// Set plan bounds.
	g.container.Layout.(*gardenLayout).box = &plan.Box

	g.Refresh()
}

// Create a new garden widget. Requires a plan.
func NewGardenWidget(plan *models.Plan) *GardenWidget {
	gardenWidget := &GardenWidget{
		container: container.New(newGardenLayout(&plan.Box)),
	}

	gardenWidget.OpenPlan(plan)

	gardenWidget.ExtendBaseWidget(gardenWidget)
	return gardenWidget
}

func (w *GardenWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.container)
}

type gardenLayout struct {
	scale float32
	box   *geometry.AxisAlignedBoundingBox
}

func (g *gardenLayout) vectorSize(v geometry.Vector) fyne.Size {
	return v.Scale(g.scale).ToSize()
}

func (g *gardenLayout) boxSize(box geometry.AxisAlignedBoundingBox) fyne.Size {
	return g.vectorSize(box.Size)
}

func (g *gardenLayout) vectorPosition(v geometry.Vector) fyne.Position {
	return v.Scale(g.scale).ToPosition()
}

func (g *gardenLayout) boxPosition(box geometry.AxisAlignedBoundingBox) fyne.Position {
	return g.vectorPosition(box.Location)
}

// Layout implements fyne.Layout.
func (g *gardenLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		featureWidget := o.(*FeatureWidget)

		o.Resize(g.boxSize(featureWidget.Feature.Box))
		o.Move(g.boxPosition(featureWidget.Feature.Box))
	}
}

// MinSize implements fyne.Layout.
func (g *gardenLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return g.boxSize(*g.box)
}

func newGardenLayout(box *geometry.AxisAlignedBoundingBox) *gardenLayout {
	return &gardenLayout{
		scale: 2,
		box:   box,
	}
}
