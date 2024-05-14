package ui

import (
	"fyne.io/fyne/v2"
	"github.com/cpgillem/garden-planner/models"
)

type GardenLayout struct {
	Plan  *models.Plan
	Scale float32
}

// Layout implements fyne.Layout.
func (g *GardenLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		featureWidget := o.(*FeatureWidget)

		o.Resize(featureWidget.Feature.Box.Size.Scale(g.Scale).ToSize())
		o.Move(featureWidget.Feature.Box.Location.Scale(g.Scale).ToPosition())
	}
}

// MinSize implements fyne.Layout.
func (g *GardenLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return g.Plan.Box.Size.Scale(g.Scale).ToSize()
}

func NewGardenLayout(plan *models.Plan) *GardenLayout {
	return &GardenLayout{
		Plan: plan,
		// Default to 2 pixels per inch for now
		Scale: 2,
	}
}
