package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/models"
)

type FeatureWidget struct {
	widget.BaseWidget

	Label *widget.Label

	Feature *models.Feature
}

// Create a new widget representing a landscaping feature.
func NewFeatureWidget(feature *models.Feature) *FeatureWidget {
	featureWidget := &FeatureWidget{
		Label:   widget.NewLabelWithData(binding.BindString(&feature.Name)),
		Feature: feature,
	}

	featureWidget.Label.Truncation = fyne.TextTruncateOff
	featureWidget.ExtendBaseWidget(featureWidget)

	return featureWidget
}

func (featureWidget *FeatureWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewCenter(featureWidget.Label)
	return widget.NewSimpleRenderer(c)
}
